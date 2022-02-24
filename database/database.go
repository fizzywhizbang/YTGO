package database

import (
	"database/sql"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Category struct {
	ID   int
	Name string
}

type Video struct {
	ID           int
	YT_videoid   string
	Title        string
	Description  string
	Publisher    string
	Publish_date int
	Downloaded   int
}

type Channel struct {
	ID              int
	Displayname     string
	Dldir           string
	Yt_channelid    string
	Lastpub         int
	Lastcheck       int
	Archive         int
	Notes           string
	Date_added      int
	Last_feed_count int
}

type Tags struct {
	ID   int
	Name string
}

type Search struct {
	ID   int
	Name string
	Link string
}

func DBCheck(dbname string) bool {
	if !Exists(dbname) {
		//create database
		os.Create(dbname)
		db := DbConnect(dbname)

		channelSQL := "CREATE TABLE `channel` ("
		channelSQL += "`id` INTEGER PRIMARY KEY AUTOINCREMENT,"
		channelSQL += "`displayname` VARCHAR(255),"
		channelSQL += "`dldir` VARCHAR(255),"
		channelSQL += "`yt_channelid` VARCHAR(255),"
		channelSQL += "`lastpub` INTEGER,"
		channelSQL += "`lastcheck` INTEGER,"
		channelSQL += "`archive` INTEGER,"
		channelSQL += "`notes` TEXT,"
		channelSQL += "`date_added` INTEGER,"
		channelSQL += "`last_feed_count` INTEGER)"
		_, err := db.Exec(channelSQL)
		if err != nil {
			log.Println(err)
		}

		statusSQL := "CREATE TABLE `status` ("
		statusSQL += "`id` INTEGER PRIMARY KEY AUTOINCREMENT,"
		statusSQL += "`status_name` VARCHAR(255))"
		_, err = db.Exec(statusSQL)

		if err != nil {
			log.Println(err)
		}
		statusInsert := "INSERT INTO status (status_name) values ('Active')" //inserting default value becasue that's the one we're gonna scan
		_, err = db.Exec(statusInsert)

		if err != nil {
			log.Println(err)
		}

		videoSQL := "CREATE TABLE `video` ("
		videoSQL += "`id` INTEGER PRIMARY KEY AUTOINCREMENT,"
		videoSQL += "`yt_videoid` VARCHAR(255),"
		videoSQL += "`title` TEXT,"
		videoSQL += "`description` TEXT,"
		videoSQL += "`publisher` VARCHAR(255),"
		videoSQL += "`publish_date` INTEGER,"
		videoSQL += "`downloaded` INTEGER,"
		videoSQL += "CONSTRAINT fk_channel FOREIGN KEY (publisher) REFERENCES channel(yt_channelid) ON DELETE CASCADE ON UPDATE CASCADE)"
		_, err = db.Exec(videoSQL)

		if err != nil {
			log.Println(err)
		}

		tagSQL := "CREATE TABLE `tags` ("
		tagSQL += "`id` INTEGER PRIMARY KEY AUTOINCREMENT,"
		tagSQL += "`tag` VARCHAR(255))"
		_, err = db.Exec(tagSQL)

		if err != nil {
			log.Println(err)
		}

		searchSQL := "CREATE TABLE `searches` ("
		searchSQL += "`id` INTEGER PRIMARY KEY AUTOINCREMENT,"
		searchSQL += "`name` VARCHAR(255),"
		searchSQL += "`link` TEXT)"
		_, err = db.Exec(searchSQL)

		if err != nil {
			log.Println(err)
		}

		defer db.Close()
		return true
	}
	return false
}

func DbConnect(dbname string) *sql.DB {
	connectString := dbname + "?_cache_size=-10000&_journal_mode=WAL"
	db, err := sql.Open("sqlite3", connectString)
	CheckErr(err, connectString)

	return db
}

//status related queries
func GetAllStatus(dbname string) *sql.Rows {

	DB := DbConnect(dbname)
	sql := "SELECT * from status order by status_name"
	results, err := DB.Query(sql)

	CheckErr(err, "Unable to get all categories (GetAllStatus func)")

	defer DB.Close()
	return results
}
func GetStatus(dbname, status string) string {
	DB := DbConnect(dbname)
	results := DB.QueryRow("SELECT * FROM status WHERE ID=?", status)
	var statusModel Category
	err := results.Scan(&statusModel.ID, &statusModel.Name)
	CheckErr(err, "Something went wrong getting category (GetStatus func)")

	defer DB.Close()
	return statusModel.Name
}
func GetStatusName(dbname, status string) string {
	DB := DbConnect(dbname)
	results := DB.QueryRow("SELECT * FROM status WHERE status_name=?", status)
	var statusModel Category
	err := results.Scan(&statusModel.ID, &statusModel.Name)
	CheckErr(err, "Unable to get category by name (GetStatusName func)")

	defer DB.Close()
	return strconv.Itoa(statusModel.ID)
}
func StatusInsert(dbname, status string) {
	DB := DbConnect(dbname)
	_, err := DB.Exec("INSERT into status (status_name) values (?)", status)
	CheckErr(err, "Unable to get insert category")
	DB.Close()
}
func StatusCount(dbname string) (count int) {
	DB := DbConnect(dbname)
	result := DB.QueryRow("Select count(*) from status")

	err := result.Scan(&count)
	if err != nil {
		return 0
	}
	defer DB.Close()
	return count
}
func StatusUpdate(dbname, id, status string) {
	DB := DbConnect(dbname)
	_, err := DB.Exec("update status set status_name=? where id=?", status, id)
	if err != nil {
		panic(err.Error())
	}
}
func GetStatusIDI(dbname, status string) int {
	DB := DbConnect(dbname)
	results := DB.QueryRow("SELECT * FROM status where status_name=?", status)
	var statusModel Category
	err := results.Scan(&statusModel.ID, &statusModel.Name)
	if err != nil {
		panic(err.Error())
	}
	defer DB.Close()
	return statusModel.ID
}

///// end status related

//tag related queries

func GetAllTags(dbname, ob string) *sql.Rows {
	DB := DbConnect(dbname)
	sql := "SELECT * from tags order by " + ob
	results, err := DB.Query(sql)

	CheckErr(err, "unable to get all tags")
	defer DB.Close()
	return results
}

func TagCount(dbname string) (count int) {
	DB := DbConnect(dbname)
	result := DB.QueryRow("Select count(*) from tags")
	err := result.Scan(&count)
	CheckErr(err, "unable to get tag count")
	defer DB.Close()
	return count
}

func TagUpdate(dbname, id, tag string) {
	DB := DbConnect(dbname)
	_, err := DB.Exec("update tags set tag=? where id=?", tag, id)
	CheckErr(err, "unable to update tag")
}
func TagInsert(dbname, tag string) {
	DB := DbConnect(dbname)
	_, err := DB.Exec("INSERT into tags (tag) values (?)", tag)
	CheckErr(err, "Unable to insert tag")
}

//end tag related queries

//video related queries
func GetLatestVideos(dbname string) *sql.Rows {
	DB := DbConnect(dbname)
	unixTime := time.Now().Unix()
	begin := unixTime - 86400
	bs := strconv.Itoa(int(begin))
	et := strconv.Itoa(int(unixTime))
	GetLatestVideos := "select * from video where publish_date between " + bs + " and " + et + " and downloaded=1 order by publish_date desc"

	results, err := DB.Query(GetLatestVideos)
	CheckErr(err, "error getting latest videos")

	defer DB.Close()
	return results

}
func GetVideoExist(dbname, videoid string) (count int) {
	DB := DbConnect(dbname)
	result := DB.QueryRow("SELECT count(*) from video where yt_videoid=?", videoid)

	err := result.Scan(&count)
	if err != nil {
		return 0
	}
	defer DB.Close()
	return count
}

func GetVideoInfo(dbname, videoid string) Video {
	DB := DbConnect(dbname)
	result := DB.QueryRow("SELECT * from video where yt_videoid=?", videoid)
	var video Video
	err := result.Scan(&video.ID, &video.YT_videoid, &video.Title, &video.Description, &video.Publisher, &video.Publish_date, &video.Downloaded)
	if err != nil {

		return Video{0, "", "", "", "", 0, 0}
	}

	defer DB.Close()
	return video
}

func InsertVideo(dbname, videoid, title, description, publisher, publish_date, downloaded string) {
	DB := DbConnect(dbname)
	if GetVideoExist(dbname, videoid) == 0 {
		titleReplaceQuotes := strings.ReplaceAll(title, `"`, `\"`)
		descriptionReplaceQuotes := strings.ReplaceAll(description, `"`, `\"`)
		query := "insert into video (yt_videoid, title, description, publisher, publish_date, downloaded) values "
		query += "(\"" + videoid + "\", \"" + titleReplaceQuotes + "\",\"" + descriptionReplaceQuotes + "\",\"" + publisher + "\",\"" + publish_date + "\",\"" + downloaded + "\")"

		_, err := DB.Exec(query)
		if err != nil {
			log.Println(err)
		}
		defer DB.Close()

		//update channel last pub
		UpdateChanLastPub(dbname, publisher, publish_date)
	} else {
		//update the video information

		_, err := DB.Exec("UPDATE video set downloaded=? where yt_videoid=?", downloaded, videoid)
		if err != nil {
			log.Println(err)
		}
		defer DB.Close()
	}

}
func GetChannelVids(dbname, publisher string) *sql.Rows {
	DB := DbConnect(dbname)
	results, err := DB.Query("SELECT * FROM video WHERE publisher=? order by publish_date", publisher)

	CheckErr(err, "Unable to get channel videos (database.go)")
	defer DB.Close()
	return results
}

///// end video related queries

//channel related queries
func CheckCount(dbname, status string) (count int) {

	DB := DbConnect(dbname)
	sql := "SELECT count(*) from channel where archive=" + status

	if status == "" {
		//count all
		sql = "SELECT count(*) from channel"
	}

	result := DB.QueryRow(sql)
	err := result.Scan(&count)
	CheckErr(err, "Unable to get count from channels")

	defer DB.Close()
	return count

}
func GetLastCheck(dbname string) int {
	DB := DbConnect(dbname)
	results := DB.QueryRow("select max(lastcheck) from channel")
	var channel Channel
	err := results.Scan(&channel.Lastcheck)
	if err != nil {
		return 0
	}
	defer DB.Close()
	return channel.Lastcheck
}

func GetChanInfo(dbname, ytid string) Channel {
	DB := DbConnect(dbname)
	results := DB.QueryRow("SELECT * FROM channel WHERE yt_channelid=?", ytid)

	var channel Channel
	err := results.Scan(&channel.ID, &channel.Displayname, &channel.Dldir, &channel.Yt_channelid, &channel.Lastpub, &channel.Lastcheck, &channel.Archive, &channel.Notes, &channel.Date_added, &channel.Last_feed_count)
	if err != nil {
		return channel
	}
	defer DB.Close()
	return channel
}

func GetChannels(dbname, arch, ob string) *sql.Rows {
	DB := DbConnect(dbname)
	results, err := DB.Query("SELECT * FROM channel where archive=? ORDER BY ? asc", arch, ob)
	CheckErr(err, "Unable to get channels (database.go)")
	defer DB.Close()
	return results
}
func GetLastDownload(dbname, chanid string) int {
	DB := DbConnect(dbname)

	results := DB.QueryRow("SELECT MAX(publish_date) as publish_date FROM video WHERE publisher=?", chanid)
	var video Video
	err := results.Scan(&video.Publish_date)
	if err != nil {
		return 0
	}
	defer DB.Close()
	return video.Publish_date
}
func DeleteChannel(dbname, chanid string) {
	DB := DbConnect(dbname)
	_, err := DB.Exec("delete from channel where yt_channelid=?", chanid)

	if err != nil {
		log.Println(err)
	}
	defer DB.Close()

}
func MoveTo(dbname, chanid, status string) {
	DB := DbConnect(dbname)

	_, err := DB.Exec("update channel set archive=? where yt_channelid=?", status, chanid)
	if err != nil {
		log.Println(err)
	}
	defer DB.Close()
}
func UpdateChanLastPub(dbname, chanid, unix string) {
	DB := DbConnect(dbname)
	_, err := DB.Exec("update channel set lastpub=? where yt_channelid=?", unix, chanid)
	if err != nil {
		log.Println(err)
	}
	defer DB.Close()
}
func UpdateChecked(dbname, chanid string) {
	DB := DbConnect(dbname)
	timestamp := time.Now().Unix()
	t := strconv.FormatInt(timestamp, 10)

	_, err := DB.Exec("UPDATE channel set lastcheck=? where yt_channelid=?", t, chanid)
	if err != nil {
		log.Println(err)
	}
	defer DB.Close()
}
func UpdateFeedCT(dbname, chanid string, feedCount int) {
	DB := DbConnect(dbname)
	_, err := DB.Exec("UPDATE channel set last_feed_count=? where yt_channelid=?", feedCount, chanid)
	if err != nil {
		panic(err.Error())
	}
	DB.Close()
}

func GetChanExist(dbname, chanid string) (count int) {
	DB := DbConnect(dbname)
	result := DB.QueryRow("SELECT count(*) from channel where yt_channelid=?", chanid)

	err := result.Scan(&count)
	if err != nil {
		return 0
	}
	defer DB.Close()
	return count
}

func ChannelSearch(dbname, GlobalStatus, str, searchType string) *sql.Rows {
	//first split the string at qoutes so we can distinguish between phrases and words
	re := regexp.MustCompile(`[^\s"]+|"([^"]*)"`)
	args := re.FindAllString(str, -1)
	query := "SELECT * from channel where "
	if searchType == "Tags" {
		if strings.Contains(str, "|") || strings.Contains(str, "&") {
			for i := 0; i < len(args); i++ {

				str := strings.Replace(args[i], "\"", "", -1)
				if str != "|" && str != "&" {
					query += "notes like \"%" + str + "%\""
				} else {
					if str == "&" {
						query += " and "
					} else {
						query += " or "
					}
				}

			}
		} else {
			for i := 0; i < len(args); i++ {
				if i < len(args) && i > 0 {
					query += " and "
				}
				str := strings.Replace(args[i], "\"", "", -1)
				query += "notes like \"%" + str + "%\""

			}
		}
	}

	if searchType == "Notes" {
		if strings.Contains(str, "|") || strings.Contains(str, "&") {
			for i := 0; i < len(args); i++ {

				str := strings.Replace(args[i], "\"", "", -1)
				if str != "|" && str != "&" {
					query += "notes like \"%" + str + "%\""
				} else {
					if str == "&" {
						query += " and "
					} else {
						query += " or "
					}
				}

			}
		} else {
			// for i := 0; i < len(args); i++ {
			// 	if i < len(args) && i > 0 {
			// 		query += " and "
			// 	}
			// 	str := strings.Replace(args[i], "\"", "", -1)
			// 	query += "notes like \"%" + str + "%\""

			// }
			query += "notes like \"%" + str + "%\""
		}
	}

	if searchType == "Channel Name" {
		if strings.Contains(str, "|") || strings.Contains(str, "&") {
			for i := 0; i < len(args); i++ {

				str := strings.Replace(args[i], "\"", "", -1)
				if str != "|" && str != "&" {
					query += "displayname like \"%" + str + "%\""
				} else {
					if str == "&" {
						query += " and "
					} else {
						query += " or "
					}
				}

			}
		} else {
			for i := 0; i < len(args); i++ {
				if i < len(args) && i > 0 {
					query += " and "
				}
				str := strings.Replace(args[i], "\"", "", -1)
				query += "displayname like \"%" + str + "%\""

			}
		}
	}
	//Channel Directory
	if searchType == "Channel Directory" {
		if strings.Contains(str, "|") || strings.Contains(str, "&") {
			for i := 0; i < len(args); i++ {

				str := strings.Replace(args[i], "\"", "", -1)
				if str != "|" && str != "&" {
					query += "dldir like \"%" + str + "%\""
				} else {
					if str == "&" {
						query += " and "
					} else {
						query += " or "
					}
				}

			}
		} else {
			for i := 0; i < len(args); i++ {
				if i < len(args) && i > 0 {
					query += " and "
				}
				str := strings.Replace(args[i], "\"", "", -1)
				query += "dldir like \"%" + str + "%\""

			}
		}
	}
	//select * from channel where yt_channelid in(select distinct(publisher) from video where title like "%something%")
	if searchType == "Channel with Video Title" {
		query += " yt_channelid in(select distinct(publisher) from video where "
		if strings.Contains(str, "|") || strings.Contains(str, "&") {
			for i := 0; i < len(args); i++ {

				str := strings.Replace(args[i], "\"", "", -1)
				if str != "|" && str != "&" {
					query += "title like \"%" + str + "%\""
				} else {
					if str == "&" {
						query += " and "
					} else {
						query += " or "
					}
				}

			}
		} else {
			for i := 0; i < len(args); i++ {
				if i < len(args) && i > 0 {
					query += " and "
				}
				str := strings.Replace(args[i], "\"", "", -1)
				query += "title like \"%" + str + "%\""

			}
		}
		query += ")"
	}
	if GlobalStatus != "" {
		//append current status to search
		query += " and archive=" + GlobalStatus
	}
	DB := DbConnect(dbname)
	results, _ := DB.Query(query)
	defer DB.Close()
	return results

}
func ModChanSettings(dbname, channelURL, newChanURL, displayname, channelDirectory, textArea string, statusSelector int) bool {
	s := statusSelector

	q := "UPDATE channel set "
	if channelURL != newChanURL {
		q += "yt_channelid=\"" + newChanURL + "\","
	}
	q += "displayname=\"" + displayname + "\","
	q += "dldir=\"" + channelDirectory + "\","
	q += "notes=\"" + textArea + "\","
	q += "archive=" + strconv.Itoa(s) + " "
	q += "where yt_channelid=\"" + channelURL + "\""
	DB := DbConnect(dbname)
	_, err := DB.Exec(q)
	if err != nil {
		return false
	}
	defer DB.Close()
	return true
}

//
func InsertChannel(dbname string, channel Channel) bool {
	query := "INSERT into channel (displayname, dldir, yt_channelid, lastcheck, archive, notes, date_added) values "
	query += "(\"" + channel.Displayname + "\",\"" + channel.Dldir + "\",\"" + channel.Yt_channelid + "\",\"" + strconv.Itoa(channel.Lastcheck) + "\","
	query += "\"" + strconv.Itoa(channel.Archive) + "\",\"" + channel.Notes + "\",\"" + strconv.Itoa(channel.Date_added) + "\")"
	DB := DbConnect(dbname)
	_, err := DB.Exec(query)

	if err != nil {
		return false
	}
	defer DB.Close()
	return true
}

//end channel related queries

//search related queries

func GetAllSearches(dbname string) *sql.Rows {
	DB := DbConnect(dbname)

	sql := "SELECT * from searches order by id"
	results, err := DB.Query(sql)

	if err != nil {
		panic(err.Error())
	}
	defer DB.Close()
	return results
}

func SearchCount(dbname string) (count int) {
	DB := DbConnect(dbname)
	result := DB.QueryRow("Select count(*) from searches")
	err := result.Scan(&count)
	if err != nil {
		return 0
	}
	defer DB.Close()
	return count
}

func SearchUpdate(dbname, id, name, link string) bool {
	DB := DbConnect(dbname)
	query := "update searches set name=\"" + name + "\", link=\"" + link + "\" where id=" + id
	_, err := DB.Exec(query)
	if err != nil {
		return false
	}
	defer DB.Close()
	return true
}
func SearchInsert(dbname, name, link string) bool {
	DB := DbConnect(dbname)
	query := "INSERT into searches (name, link) values (\"" + name + "\", \"" + link + "\")"
	_, err := DB.Exec(query)
	if err != nil {
		return false
	}
	defer DB.Close()
	return true
}

func SearchDelete(dbname, id string) bool {
	DB := DbConnect(dbname)
	query := "delete from searches where id=" + id
	_, err := DB.Exec(query)
	if err != nil {
		return false
	}
	defer DB.Close()
	return true
}

//end search related queries
