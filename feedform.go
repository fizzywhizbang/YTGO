package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aquilax/truncate"
	"github.com/fizzywhizbang/YTGO/database"
	"github.com/fizzywhizbang/YTGO/functions"
	qt "github.com/mappu/miqt/qt6"
)

func feedWindow(chanid string) {
	window := qt.NewQMainWindow(nil)
	chaninfo := database.GetChanInfo(config.Db_name, chanid)
	title := chaninfo.Displayname + " Channel Feed"
	window.SetWindowTitle(title)
	window.SetMinimumSize2(800, 400)
	window.OnKeyPressEvent(func(super func(event *qt.QKeyEvent), event *qt.QKeyEvent) {
		if int32(event.Key()) == int32(qt.Key_Escape) {
			//close window
			window.Close()
		}
	})

	mainWidget := qt.NewQWidget(nil)
	mainWidget.SetContentsMargins(0, 2, 0, 0)

	youtubefeed := YtFeedURL + chanid

	resp, err := http.Get(youtubefeed)
	// // handle the error if there is one
	functions.CheckErr(err, "Unable to get youtube feed ("+youtubefeed+")")

	byteValue, err := ioutil.ReadAll(resp.Body)

	functions.CheckErr(err, "Error reading feed body")
	// we initialize our Feed array
	var feed functions.Feed

	// we unmarshal our byteArray which contains our
	// xmlFiles content into 'users' which we defined above
	xml.Unmarshal(byteValue, &feed)

	/*
		create form widget (I believe the problem is sizing but if I can static set sizes it should be visually appealing)
		insert elements into form widget
		each form widget will have it's own connect.clicked with the value it needs to be

	*/

	formlayout := qt.NewQFormLayout(nil)
	group := qt.NewQHBoxLayout(nil)
	header1 := qt.NewQLabel3("VideoID")
	header1.SetFixedWidth(100)
	group.AddWidget(header1.QWidget)
	header2 := qt.NewQLabel3("Title")
	header2.SetFixedWidth(400)
	group.AddWidget(header2.QWidget)
	header3 := qt.NewQLabel3("Date")
	header3.SetFixedWidth(100)
	group.AddWidget(header3.QWidget)
	header4 := qt.NewQLabel3("Status")
	header4.SetFixedWidth(100)
	group.AddWidget(header4.QWidget)
	header5 := qt.NewQLabel3("Actions")
	group.AddWidget(header5.QWidget)

	formlayout.AddRowWithLayout(group.QLayout)
	database.UpdateChecked(config.Db_name, chanid)
	database.UpdateFeedCT(config.Db_name, chanid, len(feed.Entries))
	if len(feed.Entries) >= 1 {
		//db fields yt_videoid, title, description, publisher, publish_date(unix), watched(if added to download then 1 else 0)
		for i := 0; i < len(feed.Entries); i++ {
			group := qt.NewQHBoxLayout(nil)

			date, _ := time.Parse(time.RFC3339, feed.Entries[i].Published)
			videxists := false
			if database.GetVideoExist(config.Db_name, feed.Entries[i].VideoId) == 1 {
				videxists = true
			}
			//fields VideoID, Title, Date, Downloaded, Action (view,mark,etc)

			videoIDLabel := qt.NewQLabel3(feed.Entries[i].VideoId)
			videoIDLabel.SetFixedWidth(100)
			group.AddWidget(videoIDLabel.QWidget)

			titleLabel := qt.NewQLabel3(truncate.Truncate(feed.Entries[i].Title, 70, "...", truncate.PositionEnd))
			titleLabel.SetFixedWidth(400)
			group.AddWidget(titleLabel.QWidget)

			dateLabel := qt.NewQLabel3(date.Format("2006-01-02"))
			dateLabel.SetFixedWidth(100)
			group.AddWidget(dateLabel.QWidget)

			actionCombo := qt.NewQComboBox(nil)
			list := []string{"Actions", "Download " + strconv.Itoa(i) + "", "Skip " + strconv.Itoa(i) + "", "View " + strconv.Itoa(i) + "", "Find Similar " + strconv.Itoa(i) + ""}
			actionCombo.AddItems(list)

			status := "False"
			if videxists {
				videoData := database.GetVideoInfo(config.Db_name, feed.Entries[i].VideoId)
				fmt.Println(videoData.ID, videoData.Downloaded)
				status = "Queued"
				if videoData.Downloaded == 1 {
					status = "Downloaded"
				}
				if videoData.Downloaded == 2 {
					status = "Skipped"
				}

			}

			downloadedLabel := qt.NewQLabel3(status)
			downloadedLabel.SetFixedWidth(100)
			group.AddWidget(downloadedLabel.QWidget)

			group.AddWidget(actionCombo.QWidget)
			actionCombo.OnCurrentTextChanged(func(text string) {
				action := strings.Split(text, " ")

				if action[0] == "Download" {
					row, _ := strconv.Atoi(action[1])
					fmt.Println("Download", feed.Entries[row].VideoId)
					date, _ := time.Parse(time.RFC3339, feed.Entries[row].Published)
					functions.MkCrawljob(config.Db_name, config.FolderWatch, GlobalChannelID, feed.Entries[row].Title, feed.Entries[row].VideoId, date.Format("2006-01-02"), 1)
					qt.QMessageBox_Information(nil, "OK", "Added to Queue "+feed.Entries[row].Title)
				}
				if action[0] == "Skip" {
					row, _ := strconv.Atoi(action[1])
					fmt.Println("Skip", feed.Entries[row].VideoId)
					date, _ := time.Parse(time.RFC3339, feed.Entries[row].Published)
					unixdate := functions.ConvertYMDtoUnix(date.Format("2006-01-02"))
					database.InsertVideo(config.Db_name, feed.Entries[row].VideoId, feed.Entries[row].Title, feed.Entries[row].Title, GlobalChannelID, unixdate, "2")
					qt.QMessageBox_Information(nil, "OK", feed.Entries[row].Title+" Recorded as skipped")
				}
				if action[0] == "View" {
					row, _ := strconv.Atoi(action[1])
					url := YtWatchPrefix + feed.Entries[row].VideoId
					functions.Openbrowser(url, config.Defbrowser)
				}
				if action[0] == "Find" {
					row, _ := strconv.Atoi(action[2])
					title := strings.Replace(feed.Entries[row].Title, " ", "+", -1)
					url := YtSearchPrefix + title + "&sp=CAI%253D" //order by upload date
					functions.Openbrowser(url, config.Defbrowser)
				}
			})

			formlayout.AddRowWithLayout(group.QLayout)
		}
	}

	mainWidget.SetLayout(formlayout.QLayout)
	window.SetCentralWidget(mainWidget)
	window.Show()
}
