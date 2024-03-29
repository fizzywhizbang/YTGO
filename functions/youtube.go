package functions

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dyatlov/go-opengraph/opengraph"
	"github.com/fizzywhizbang/YTGO/database"
	"github.com/therecipe/qt/widgets"
)

var Found = false
var S interface{} = ""

// youtube feed struct
type Feed struct {
	XMLName xml.Name `xml:"feed"`
	Title   string   `xml:"title"`
	Entries []Entry  `xml:"entry"`
}

// youtube feed entry struct
type Entry struct {
	XMLName   xml.Name   `xml:"entry"`
	ID        string     `xml:"id"`
	VideoId   string     `xml:"videoId"`
	Chanid    string     `xml:"channelId"`
	Title     string     `xml:"title"`
	Published string     `xml:"published"`
	MGroup    MediaGroup `xml:"group"`
}

// youtube media group struct
type MediaGroup struct {
	XMLName     xml.Name `xml:"group"`
	Description string   `xml:"description"`
}

// url constants
const (
	YtVideoInfoURL = "https://www.youtube.com/get_video_info?video_id="
	YtFeedURL      = "https://www.youtube.com/feeds/videos.xml?channel_id="
	YtWatchPrefix  = "https://www.youtube.com/watch?v="
	YtChanPrefix   = "https://www.youtube.com/channel/"
	YtSearchPrefix = "https://www.youtube.com/results?search_query="
)

func getChanID(url string) string {

	data, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(data.Body)
	if err != nil {
		log.Fatal(err)
	}

	s := doc.Find("script")
	for i := range s.Nodes {

		if strings.Contains(s.Eq(i).Text(), "ytInitialData") && !strings.Contains(s.Eq(i).Text(), "(function") {

			text := strings.Replace(s.Eq(i).Text(), "var ytInitialData = ", "", -1)
			text = strings.Replace(text, ";", "", -1)

			m := map[string]interface{}{}

			// Parsing/Unmarshalling JSON encoding/json
			err := json.Unmarshal([]byte(text), &m)

			if err != nil {
				panic(err)
			}
			parseMap(m)
			return S.(string)

		}

	}
	// Creating the maps for JSON
	return "Unknown"
}

func parseMap(aMap map[string]interface{}) interface{} {
	uc := ""

	Found = false
	for key, val := range aMap {
		switch concreteVal := val.(type) {
		case map[string]interface{}:
			if !Found {
				parseMap(val.(map[string]interface{}))
			}

		default:
			uc = fmt.Sprintf("%v", concreteVal)

			if len(key) > 1 && strings.HasPrefix(uc, "UC") && !Found {
				Found = true

				S = concreteVal
				return concreteVal
			}

		}
	}
	return uc
}

// get channel information from youtube
func GetChanInfoFromYT(chanid string) database.Channel {
	url := YtChanPrefix + chanid

	if strings.HasPrefix(chanid, "UC") {
		url = YtChanPrefix + chanid

	} else {
		//get UC value
		if strings.HasPrefix(chanid, "@") {
			ucChanID := getChanID("https://www.youtube.com/" + chanid)
			url = YtChanPrefix + ucChanID

			chanid = ucChanID
		} else {
			ucChanID := getChanID("https://www.youtube.com/user/" + chanid)
			url = YtChanPrefix + ucChanID
			chanid = ucChanID
		}

	}
	fmt.Println(url)
	var channel database.Channel

	resp, err := http.Get(url)
	// // handle the error if there is one
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	//get opengraph content
	og := opengraph.NewOpenGraph()
	nerr := og.ProcessHTML(strings.NewReader(string(html)))
	if nerr != nil {
		fmt.Println(nerr)
	}
	channel.Displayname = og.Title
	channel.Notes = og.Description
	// chaninfo := []string{og.Title, og.Description}
	return channel
}

// update channel video counts and stuff
func UpdateChan(dbname, fwatch, chanid string, dl bool, msg bool) int {
	youtubefeed := YtFeedURL + chanid

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
	x := 0
	if len(feed.Entries) >= 1 {
		//db fields yt_videoid, title, description, publisher, publish_date(unix), watched(if added to download then 1 else 0)
		for i := 0; i < len(feed.Entries); i++ {
			//check and if exist do nothing
			unixdate := DateConvertToUnix(feed.Entries[i].Published)

			exists := database.GetVideoExist(dbname, feed.Entries[i].VideoId)
			// fmt.Println(feed.Entries[i].VideoId + " " + strconv.Itoa(exists))
			if exists == 0 {

				//if dl == true then create download and set downloaded == 1 else add to database with value of 2 Skipped
				title := strings.Replace(feed.Entries[i].Title, "\\", "", -1)
				if dl {

					database.InsertVideo(dbname, feed.Entries[i].VideoId, title, title, chanid, unixdate, "1")
					MkCrawljob(dbname, fwatch, chanid, feed.Entries[i].Title, feed.Entries[i].VideoId, feed.Entries[i].Published, 0)
				} else {
					database.InsertVideo(dbname, feed.Entries[i].VideoId, title, title, chanid, unixdate, "2")
				}
				x++
			}
		}
	}
	if msg {
		if dl {
			widgets.QMessageBox_Information(nil, "Videos Added", strconv.Itoa(x)+" Videos added to the queue", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		} else {
			widgets.QMessageBox_Information(nil, "Videos Added", strconv.Itoa(x)+" Videos added to the database", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		}
	}

	//update last check timestamp
	database.UpdateChecked(dbname, chanid)
	database.UpdateFeedCT(dbname, chanid, len(feed.Entries))
	return x
}

// get information about a specific video
func GetVideoInfo(videoid string) database.Video {
	var video database.Video
	fmt.Println(video)
	url := YtWatchPrefix + videoid
	//we'll get it from the OG data
	//get opengraph content
	resp, err := http.Get(url)
	// // handle the error if there is one
	if err != nil {
		fmt.Println(err)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	datePub := ""
	publisher := ""
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		prop, _ := s.Attr("itemprop")
		if prop == "datePublished" {
			val, bVal := s.Attr("content")
			if bVal {
				datePub = val
			}

		}
	})
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		prop2, _ := s.Attr("itemprop")
		if prop2 == "channelId" {
			val2, bVal2 := s.Attr("content")
			if bVal2 {
				publisher = val2
			}
		}
	})
	title := ""
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		prop2, _ := s.Attr("itemprop")
		if prop2 == "name" {
			val2, bVal2 := s.Attr("content")
			if bVal2 {
				title = val2
			}
		}
	})
	description := ""
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		prop2, _ := s.Attr("itemprop")
		if prop2 == "description" {
			val2, bVal2 := s.Attr("content")
			if bVal2 {
				description = val2
			}
		}
	})
	uxDatePub := ConvertYMDtoUnix(datePub)
	video.YT_videoid = videoid
	video.Title = title
	video.Description = description
	video.Publisher = publisher
	dt, _ := strconv.Atoi(uxDatePub)
	video.Publish_date = dt
	video.Downloaded = 0

	defer resp.Body.Close()
	return video
}

// create a crawljob for jdownloader
func MkCrawljob(dbname, fwatch, chanid, title, videoid, date string, updatedb int) {
	chaninfo := database.GetChanInfo(dbname, chanid)
	packagename := "<jd:packagename>"
	filename := fwatch + chaninfo.Displayname + "_" + videoid + ".crawljob"
	file, err := os.Create(filename)
	CheckErr(err, "Unable to create file for crawljob")
	defer file.Close()
	fmt.Fprintln(file, "#chantitle "+chaninfo.Displayname)
	fmt.Fprintln(file, "#download "+title)
	url := YtWatchPrefix + videoid
	fmt.Fprintln(file, "text=\""+url+"\"")
	fmt.Fprintln(file, "autoConfirm=TRUE")
	fmt.Fprintln(file, "autoStart=TRUE")
	fmt.Fprintln(file, "downloadFolder="+chaninfo.Dldir+"/"+packagename)
	fmt.Fprintln(file, "downloadPassword=null")
	fmt.Fprintln(file, "enabled=true")
	fmt.Fprintln(file, "forcedStart=Default")
	fmt.Fprintln(file, "priority=DEFAULT")

	//check channel id before video because we might want to download a video for which there is no sub
	if database.GetChanExist(dbname, chanid) != 0 {
		if updatedb == 1 {
			unixdate := ConvertYMDtoUnix(date)
			database.InsertVideo(dbname, videoid, title, title, chanid, unixdate, "1")

		}

	}

}
