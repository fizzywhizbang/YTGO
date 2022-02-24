package monitor

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
)

var orderby = "id"
var DB *sql.DB
var ConfigFile = ""
var ConfigDir = ""

func main() {
	MonitorStart(ConfigFile)
}

func MonitorStart(configfile string) {
	if len(configfile) == 0 {
		homedir, err := os.UserHomeDir()
		checkErr(err)
		ConfigFile = homedir + ".config/ytmon/ytmon.json"
	} else {
		ConfigFile = configfile
	}
	executeChannelMonitor()
	executeQueueMonitor()
}

func queueCheck() {
	//if queue is running ignore signal

	videos := getVideoForQueue()
	videosUpdate := []string{}
	for videos.Next() {
		var video Video
		err := videos.Scan(&video.id, &video.yt_videoid, &video.title, &video.description, &video.publisher, &video.publish_date, &video.watched)
		if err != nil {
			fmt.Println("No videos to queue")
		} else {
			chaninfo := getChanInfo2(video.publisher)
			processedString := strings.Join(cleanText(video.title), " ")

			savename := chaninfo.displayname + " " + dateConvertTrim(video.publish_date, 10) + " " + processedString
			fmt.Println("downloading ", savename, " ", video.publish_date)

			//download queued video and when complete update status

			mkCrawljob(chaninfo.yt_channelid, video.title, video.yt_videoid, dateConvertTrim(video.publish_date, 10), 0)
			//mark video downloaded
			videosUpdate = append(videosUpdate, video.yt_videoid)

		}
	}

	//no update videos downloaded
	for i := 0; i < len(videosUpdate); i++ {
		updateVideoStatus(videosUpdate[i])
	}

}

func channelCheck() {
	lastCheck := getLastCheck()
	now := time.Now().Unix()

	diff := now - int64(lastCheck)
	channelQueue := []string{}
	channelDisplaynames := []string{}
	fmt.Println("last check ", diff, " seconds ago")
	if diff >= 1800 { //this will stop us from getting banned because we kept restarting our program :) 1/2 hour but scans are normally 1 hour interval
		//put results in a slice to free up the database because sqlite doesn't like to share
		channels := getChannels("1", orderby, "asc")

		for channels.Next() {
			var channel Channel
			err := channels.Scan(&channel.id, &channel.displayname, &channel.dldir, &channel.yt_channelid, &channel.lastpub, &channel.lastcheck, &channel.archive, &channel.notes, &channel.date_added, &channel.last_feed_count)
			if err != nil {
				fmt.Println("something went wrong with the channel scan")
			}

			channelQueue = append(channelQueue, channel.yt_channelid)
			channelDisplaynames = append(channelDisplaynames, channel.displayname)
		}
		for i := 0; i < len(channelQueue); i++ {
			fmt.Println("Checking ", channelDisplaynames[i], " for updates ", time.Now())
			getChannelVideos(channelQueue[i])
			time.Sleep(time.Second * 3) //wait three seconds between checks so as no to piss off youtube

		}
	} else {
		fmt.Println("Queue locked defering channel check until later")
	}

}

func cleanText(text string) []string {
	words := regexp.MustCompile(`[\p{L}\d_]+`)
	return words.FindAllString(text, -1)
}

func getChannelVideos(chanid string) {
	youtubefeed := ytFeedURL + chanid

	resp, err := http.Get(youtubefeed)
	// // handle the error if there is one
	if err != nil {
		panic(err)
	}
	byteValue, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	// we initialize our Users array
	var feed Feed

	// we unmarshal our byteArray which contains our
	// xmlFiles content into 'users' which we defined above
	xml.Unmarshal(byteValue, &feed)

	resultCount := 0

	if len(feed.Entries) >= 1 {
		//db fields yt_videoid, title, description, publisher, publish_date(unix), watched(if added to download then 1 else 0)

		for i := 0; i < len(feed.Entries); i++ {
			//check and if exist do nothing

			date := friendlyDate(feed.Entries[i].Published)

			unixdate := convertYMDtoUnix(date)
			exists := getVideoExist(feed.Entries[i].VideoId)
			fmt.Println(feed.Entries[i].VideoId + " " + strconv.Itoa(exists))
			if exists == 0 {
				//insert into database with watched status 0 and begin queue check
				insertVideo(feed.Entries[i].VideoId, feed.Entries[i].Title, feed.Entries[i].Title, chanid, unixdate, "0")
				i++
				resultCount++
			}

		}

	}

	//update last check timestamp

}

func executeQueueMonitor() {
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.SingletonMode()
	scheduler.Every(10).Minutes().Do(queueCheck)
	scheduler.StartBlocking()
}

func executeChannelMonitor() {
	//set download queue to true to keep the database from being locked
	scheduler2 := gocron.NewScheduler(time.UTC)
	scheduler2.Every(60).Minutes().Do(channelCheck)
	scheduler2.StartBlocking()
}
