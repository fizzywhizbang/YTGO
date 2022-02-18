package database

import (
	"database/sql"
	"log"
	"os"

	"github.com/fizzywhizbang/YTGO/functions"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

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
	Watched      int
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
		defer db.Close()
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
		videoSQL += "`watched` INTEGER,"
		videoSQL += "CONSTRAINT fk_channel FOREIGN KEY (publisher) REFERENCES channel(yt_channelid) ON DELETE CASCADE ON UPDATE CASCADE)"
		_, err = db.Exec(videoSQL)

		if err != nil {
			log.Println(err)
		}

		return true
	}
	return false
}

func DbConnect(dbname string) *sql.DB {
	connectString := dbname + "?_cache_size=-10000&_journal_mode=WAL"
	db, err := sql.Open("sqlite3", connectString)
	if err != nil {
		log.Println("Unable to connect to the database (Line 29)")
	}
	return db
}
