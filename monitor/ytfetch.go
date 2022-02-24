package monitor

import (
	"encoding/xml"
	"fmt"
	"os"
)

type Feed struct {
	XMLName xml.Name `xml:"feed"`
	Title   string   `xml:"title"`
	Entries []Entry  `xml:"entry"`
}

type Entry struct {
	XMLName   xml.Name   `xml:"entry"`
	ID        string     `xml:"id"`
	VideoId   string     `xml:"videoId"`
	Chanid    string     `xml:"channelId"`
	Title     string     `xml:"title"`
	Published string     `xml:"published"`
	MGroup    MediaGroup `xml:"group"`
}

type MediaGroup struct {
	XMLName     xml.Name `xml:"group"`
	Description string   `xml:"description"`
}

func mkCrawljob(chanid string, title string, videoid string, date string, updatedb int) {
	chaninfo := getChanInfo2(chanid)
	config := loadConfig()
	packagename := "<jd:packagename>"
	filename := config.FolderWatch + chaninfo.displayname + "_" + videoid + ".crawljob"
	file, err := os.Create(filename)
	checkErr(err)
	defer file.Close()
	fmt.Fprintln(file, "#chantitle "+chaninfo.displayname)
	fmt.Fprintln(file, "#download "+title)
	url := ytWatchPrefix + videoid
	fmt.Fprintln(file, "text=\""+url+"\"")
	fmt.Fprintln(file, "autoConfirm=TRUE")
	fmt.Fprintln(file, "autoStart=TRUE")
	fmt.Fprintln(file, "downloadFolder="+chaninfo.dldir+"/"+packagename)
	fmt.Fprintln(file, "downloadPassword=null")
	fmt.Fprintln(file, "enabled=true")
	fmt.Fprintln(file, "forcedStart=Default")
	fmt.Fprintln(file, "priority=DEFAULT")

	//check channel id before video because we might want to download a video for which there is no sub
	if getChanName(chanid) != "None" {
		if updatedb == 1 {
			if getVideoExist(videoid) == 0 {
				//if this is true then it's coming from the feed
				unixdate := convertYMDtoUnix(date)
				insertVideo(videoid, title, title, chanid, unixdate, "1")
			}
		}

	}

}
