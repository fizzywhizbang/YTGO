package database

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/fizzywhizbang/YTGO/functions"
	_ "github.com/mattn/go-sqlite3"
)

type Status struct {
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

func DBCheck(dbname string) bool {
	if !functions.Exists(dbname) {
		//create database
		os.Create(dbname)
		db := DbConnect(dbname)

		channelSQL := "CREATE TABLE `channel` ("
		channelSQL += "`id` INTEGER PRIMARY KEY AUTOINCREMENT,"
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
		defer db.Close()
		return true
	}
	return false
}

func DbConnect(dbname string) *sql.DB {
	connectString := dbname + "?_cache_size=-10000&_journal_mode=WAL"
	db, err := sql.Open("sqlite3", connectString)
	functions.CheckErr(err, connectString)

	return db
}

//status related queries
func GetAllStatus(dbname string) *sql.Rows {

	DB := DbConnect(dbname)
	sql := "SELECT * from status order by status_name"
	results, err := DB.Query(sql)

	functions.CheckErr(err, "Unable to get all statuses (GetAllStatus func)")

	defer DB.Close()
	return results
}
func GetStatus(dbname, status string) string {
	DB := DbConnect(dbname)
	results := DB.QueryRow("SELECT * FROM status WHERE ID=?", status)
	var statusModel Status
	err := results.Scan(&statusModel.ID, &statusModel.Name)
	functions.CheckErr(err, "Something went wrong getting status (GetStatus func)")

	defer DB.Close()
	return statusModel.Name
}
func GetStatusName(dbname, status string) string {
	DB := DbConnect(dbname)
	results := DB.QueryRow("SELECT * FROM status WHERE status_name=?", status)
	var statusModel Status
	err := results.Scan(&statusModel.ID, &statusModel.Name)
	functions.CheckErr(err, "Unable to get status by name (GetStatusName func)")

	defer DB.Close()
	return strconv.Itoa(statusModel.ID)
}
func CheckCount(dbname, status string) (count int) {

	DB := DbConnect(dbname)
	sql := "SELECT count(*) from channel where archive=" + status

	if status == "" {
		//count all
		sql = "SELECT count(*) from channel"
	}

	result := DB.QueryRow(sql)
	err := result.Scan(&count)
	functions.CheckErr(err, "Unable to get count from channels")

	defer DB.Close()
	return count

}

///// end status related

//video related queries
func GetLatestVideos(dbname string) *sql.Rows {
	DB := DbConnect(dbname)
	unixTime := time.Now().Unix()
	begin := unixTime - 86400
	bs := strconv.Itoa(int(begin))
	et := strconv.Itoa(int(unixTime))
	GetLatestVideos := "select * from video where publish_date between " + bs + " and " + et + " order by publish_date desc"

	results, err := DB.Query(GetLatestVideos)
	functions.CheckErr(err, "error getting latest videos")

	defer DB.Close()
	return results

}

///// end video related queries

//channel related queries
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

//end channel related queries
