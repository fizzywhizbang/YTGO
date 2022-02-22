package main

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
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
	layout := widgets.NewQVBoxLayout()

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
	treeWidget := widgets.NewQTreeWidget(nil)
	treeWidget.SetColumnCount(6)
	treeWidget.SetObjectName("treewidget")
	treeWidget.Header().SetStretchLastSection(false)
	treeWidget.Header().SetSectionsClickable(true)
	treeWidget.SetAlternatingRowColors(true)
	tableColors := "alternate-background-color: #88DD88; background-color:#FFFFFF; color:#000000; font-size: 12px;"
	treeWidget.SetStyleSheet(tableColors)
	treeWidget.Header()
	treeWidget.SetHeaderLabels([]string{"VidoeID", "Title", "Date", "Status", "View", "Mark"})

	if len(feed.Entries) >= 1 {
		//db fields yt_videoid, title, description, publisher, publish_date(unix), watched(if added to download then 1 else 0)
		for i := 0; i < len(feed.Entries); i++ {

			date, _ := time.Parse(time.RFC3339, feed.Entries[i].Published)
			videxists := "False"
			if database.GetVideoExist(config.Db_name, feed.Entries[i].VideoId) == 1 {
				videxists = "True"
			}
			treewidgetItem := widgets.NewQTreeWidgetItem2([]string{feed.Entries[i].VideoId, truncate.Truncate(feed.Entries[i].Title, 70, "...", truncate.PositionEnd), date.Format("2006-01-02"), videxists, "View", "Mark"}, 0)
			treewidgetItem.SetData(0, int(core.Qt__UserRole), core.NewQVariant12(feed.Entries[i].VideoId))

			treeWidget.AddTopLevelItem(treewidgetItem)
		}
	}
	database.UpdateFeedCT(config.Db_name, chanid, len(feed.Entries))
	database.UpdateChecked(config.Db_name, chanid)
	treeWidget.ResizeColumnToContents(0)
	treeWidget.ResizeColumnToContents(1)
	treeWidget.ResizeColumnToContents(2)
	treeWidget.ResizeColumnToContents(3)
	treeWidget.ResizeColumnToContents(4)
	treeWidget.ResizeColumnToContents(5)
	treeWidget.ConnectDoubleClicked(func(index *core.QModelIndex) {
		data := index.Data(int(core.Qt__UserRole)).ToString()
		item := treeWidget.CurrentColumn()

		//get item
		//one item will download, one will view in browser, one will search for similar subjects
		//https://www.youtube.com/watch?v=
		if item == 0 {
			// fmt.Println(data)

			functions.MkCrawljob(config.Db_name, config.FolderWatch, chanid, treeWidget.CurrentItem().Text(1), data, treeWidget.CurrentItem().Text(2), 1)
			//double click means open the settings for this channel
			widgets.QMessageBox_Information(nil, "OK", "Added to Queue "+treeWidget.CurrentItem().Text(1), widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)

		}
		if item == 4 {
			url := YtWatchPrefix + treeWidget.CurrentItem().Text(0)
			functions.Openbrowser(config.Defbrowser, url)
		}
		if item == 5 {
			if database.GetVideoExist(config.Db_name, treeWidget.CurrentItem().Text(0)) == 0 {
				//if this is true then it's coming from the feed
				unixdate := functions.ConvertYMDtoUnix(treeWidget.CurrentItem().Text(2))
				database.InsertVideo(config.Db_name, treeWidget.CurrentItem().Text(0), treeWidget.CurrentItem().Text(1), treeWidget.CurrentItem().Text(1), chanid, unixdate, "1")
				widgets.QMessageBox_Information(nil, "OK", treeWidget.CurrentItem().Text(1)+" Recorded as downloaded", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
			}
		}
		if item == 1 {
			text := cleanUnicode(treeWidget.CurrentItem().Text(item))
			//seting sort order by date uploaded
			/*
				sp=CAI%253D #search by upload date
				sp=CAASAhAB #search by  relevance
				sp=EgQIAhAB #today and relevance
				sp=CAISBAgCEAE%253D #today & upload date
				sp=CAMSBAgCEAE%253D #today and view count
				sp=CAESBAgCEAE%253D #today and rating
				sp=CAESAhAB #rating only
			*/
			url := YtSearchPrefix + text + "&sp=CAI%253D" //order by upload date
			functions.Openbrowser(config.Defbrowser, url)
		}

	})

	layout.AddWidget(treeWidget, 0, 0)
	mainWidget.SetLayout(layout)
	window.SetCentralWidget(mainWidget)
	window.Show()
}
