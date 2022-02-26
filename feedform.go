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
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

func feedWindow(chanid string) {
	window := widgets.NewQMainWindow(nil, 0)
	chaninfo := database.GetChanInfo(config.Db_name, chanid)
	title := chaninfo.Displayname + " Channel Feed"
	window.SetWindowTitle(title)
	window.SetMinimumSize2(800, 400)
	window.ConnectKeyPressEvent(func(e *gui.QKeyEvent) {
		if int32(e.Key()) == int32(core.Qt__Key_Escape) {
			//close window
			window.Close()
		}
	})

	mainWidget := widgets.NewQWidget(nil, 0)
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

	formlayout := widgets.NewQFormLayout(nil)
	group := widgets.NewQHBoxLayout()
	header1 := widgets.NewQLabel2("VideoID", nil, 0)
	header1.SetFixedWidth(100)
	group.AddWidget(header1, 0, 0)
	header2 := widgets.NewQLabel2("Title", nil, 0)
	header2.SetFixedWidth(400)
	group.AddWidget(header2, 0, 0)
	header3 := widgets.NewQLabel2("Date", nil, 0)
	header3.SetFixedWidth(100)
	group.AddWidget(header3, 0, 0)
	header4 := widgets.NewQLabel2("Status", nil, 0)
	header4.SetFixedWidth(100)
	group.AddWidget(header4, 0, 0)
	header5 := widgets.NewQLabel2("Actions", nil, 0)
	group.AddWidget(header5, 0, 0)

	formlayout.AddRow6(group)
	database.UpdateChecked(config.Db_name, chanid)
	if len(feed.Entries) >= 1 {
		//db fields yt_videoid, title, description, publisher, publish_date(unix), watched(if added to download then 1 else 0)
		for i := 0; i < len(feed.Entries); i++ {
			group := widgets.NewQHBoxLayout()

			date, _ := time.Parse(time.RFC3339, feed.Entries[i].Published)
			videxists := false
			if database.GetVideoExist(config.Db_name, feed.Entries[i].VideoId) == 1 {
				videxists = true
			}
			//fields VideoID, Title, Date, Downloaded, Action (view,mark,etc)

			videoIDLabel := widgets.NewQLabel2(feed.Entries[i].VideoId, nil, 0)
			videoIDLabel.SetFixedWidth(100)
			group.AddWidget(videoIDLabel, 0, 0)

			titleLabel := widgets.NewQLabel2(truncate.Truncate(feed.Entries[i].Title, 70, "...", truncate.PositionEnd), nil, 0)
			titleLabel.SetFixedWidth(400)
			group.AddWidget(titleLabel, 0, 0)

			dateLabel := widgets.NewQLabel2(date.Format("2006-01-02"), nil, 0)
			dateLabel.SetFixedWidth(100)
			group.AddWidget(dateLabel, 0, 0)

			actionCombo := widgets.NewQComboBox(nil)
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

			downloadedLabel := widgets.NewQLabel2(status, nil, 0)
			downloadedLabel.SetFixedWidth(100)
			group.AddWidget(downloadedLabel, 0, 0)

			group.AddWidget(actionCombo, 0, 0)
			actionCombo.ConnectCurrentTextChanged(func(text string) {
				action := strings.Split(text, " ")

				if action[0] == "Download" {
					row, _ := strconv.Atoi(action[1])
					fmt.Println("Download", feed.Entries[row].VideoId)
					date, _ := time.Parse(time.RFC3339, feed.Entries[row].Published)
					functions.MkCrawljob(config.Db_name, config.FolderWatch, GlobalChannelID, feed.Entries[row].Title, feed.Entries[row].VideoId, date.Format("2006-01-02"), 1)
					widgets.QMessageBox_Information(nil, "OK", "Added to Queue "+feed.Entries[row].Title, widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
				}
				if action[0] == "Skip" {
					row, _ := strconv.Atoi(action[1])
					fmt.Println("Skip", feed.Entries[row].VideoId)
					date, _ := time.Parse(time.RFC3339, feed.Entries[row].Published)
					unixdate := functions.ConvertYMDtoUnix(date.Format("2006-01-02"))
					database.InsertVideo(config.Db_name, feed.Entries[row].VideoId, feed.Entries[row].Title, feed.Entries[row].Title, GlobalChannelID, unixdate, "2")
					widgets.QMessageBox_Information(nil, "OK", feed.Entries[row].Title+" Recorded as skipped", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
				}
				if action[0] == "View" {
					row, _ := strconv.Atoi(action[1])
					url := YtWatchPrefix + feed.Entries[row].VideoId
					functions.Openbrowser(url, config.Defbrowser)
				}
				if action[0] == "Find" {
					row, _ := strconv.Atoi(action[1])
					url := YtSearchPrefix + feed.Entries[row].VideoId + "&sp=CAI%253D" //order by upload date
					functions.Openbrowser(url, config.Defbrowser)
				}
			})

			formlayout.AddRow6(group)
		}
	}

	mainWidget.SetLayout(formlayout)
	window.SetCentralWidget(mainWidget)
	window.Show()
}
