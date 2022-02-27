package monitor

import (
	"database/sql"
	"strings"
)

type ResultRow map[string]interface{}

const (
	GetStatus = "SELECT * FROM status WHERE ID="

	GetStatusID = "SELECT * FROM status where status_name="

	GetActive = "SELECT * FROM channel where archive = "

	GetChanInfo = "SELECT * FROM channel WHERE yt_channelid="

	GetChanVids = "SELECT * FROM video WHERE publisher="

	GetVideoStatus = "SELECT downloaded FROM video WHERE yt_videoid="

	GetLatestVideos = "select * from video where publish_date between unix_timestamp() - 86400 and unix_timestamp() order by publish_date desc"

	GetLastCheck = "select max(lastcheck) from channel"

	InsertNewChannel = "INSERT into channel (displayname, dldir, yt_channelid, lastcheck, archive, notes, date_added) values "
)

func ConnectDB() *sql.DB {
	config := loadConfig()

	connectString := config.Db_name + "?_cache_size=-10000&_journal_mode=WAL&_fk=true"
	db, err := sql.Open("sqlite3", connectString)
	CheckErr(err, connectString)

	return db

}
func updateVideoStatus(videoid string) {
	DB = ConnectDB()
	_, err := DB.Exec("update video set downloaded=0 where yt_videoid==?", videoid)
	if err != nil {
		panic(err.Error())
	}
	DB.Close()
}

func insertUpdate(q string) bool {
	DB = ConnectDB()
	query, err := DB.Query(q)

	if err != nil {
		panic(err.Error())
	}
	query.Close()
	defer DB.Close()
	return true
}

func getVideoForQueue() *sql.Rows {
	DB = ConnectDB()
	results, err := DB.Query("select * from video where downloaded=0")

	if err != nil {
		return results
	}
	defer DB.Close()
	return results
}

func getLastCheck() int {
	DB = ConnectDB()
	results := DB.QueryRow(GetLastCheck)
	var channel Channel
	err := results.Scan(&channel.lastcheck)
	if err != nil {
		panic(err.Error())
	}
	defer DB.Close()
	return channel.lastcheck
}

func getChanInfo(ytid string) *sql.Row {
	DB = ConnectDB()
	sql := GetChanInfo + "'" + ytid + "'"
	results := DB.QueryRow(sql)
	defer DB.Close()
	return results
}

func getChanInfo2(ytid string) Channel {
	chaninfo := getChanInfo(ytid)

	var channel Channel

	err := chaninfo.Scan(&channel.id, &channel.displayname, &channel.dldir, &channel.yt_channelid, &channel.lastpub, &channel.lastcheck, &channel.archive, &channel.notes, &channel.date_added, &channel.last_feed_count)
	if err != nil {
		return channel
	}
	return channel
}
func getChanName(ytid string) string {
	chaninfo := getChanInfo(ytid)
	var channel Channel

	err := chaninfo.Scan(&channel.id, &channel.displayname, &channel.dldir, &channel.yt_channelid, &channel.lastpub, &channel.lastcheck, &channel.archive, &channel.notes, &channel.date_added, &channel.last_feed_count)
	if err != nil {
		return "None"
	}
	return channel.displayname
}

func getChannels(arch string, ob string, sb string) *sql.Rows {
	DB = ConnectDB()
	query := GetActive + arch + " order by " + ob + " " + sb
	results, _ := DB.Query(query)
	defer DB.Close()
	return results
}

func getVideoExist(videoid string) (count int) {
	DB = ConnectDB()
	sql := "SELECT count(*) from video where yt_videoid=\"" + videoid + "\""

	result := DB.QueryRow(sql)

	err := result.Scan(&count)
	if err != nil {
		return 0
	}
	defer DB.Close()
	return count
}

func updateChanLastPub(chanid string, unix string) {
	query := "update channel set lastpub=\"" + unix + "\" where yt_channelid=\"" + chanid + "\""
	insertUpdate(query)
}

func insertVideo(videoid string, title string, description string, publisher string, publish_date string, downloaded string) {
	titleReplaceQuotes := strings.ReplaceAll(title, `"`, `\"`)
	descriptionReplaceQuotes := strings.ReplaceAll(description, `"`, `\"`)
	query := "insert into video (yt_videoid, title, description, publisher, publish_date, downloaded) values "
	query += "(\"" + videoid + "\", \"" + titleReplaceQuotes + "\",\"" + descriptionReplaceQuotes + "\",\"" + publisher + "\",\"" + publish_date + "\",\"" + downloaded + "\")"

	insertUpdate(query)
	//update channel last pub
	updateChanLastPub(publisher, publish_date)

}
