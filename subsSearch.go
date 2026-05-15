package main

import (
	"fmt"
	"strconv"

	"github.com/fizzywhizbang/YTGO/database"
	"github.com/fizzywhizbang/YTGO/functions"
	qt "github.com/mappu/miqt/qt6"
)

func showSubsSearch(searchstring string, searchType string, status string) {
	GlobalStatus = status

	verticalLayout := qt.NewQVBoxLayout(nil)

	mainWidget := qt.NewQWidget(nil)

	//this is where it differs
	channels := database.ChannelSearch(config.Db_name, GlobalStatus, searchstring, searchType)
	treeWidget := qt.NewQTreeWidget(nil)
	treeWidget.SetColumnCount(6)
	treeWidget.SetObjectName(*qt.NewQAnyStringView3("treewidget"))
	treeWidget.Header().SetStretchLastSection(false)
	treeWidget.Header().SetSectionsClickable(true)
	treeWidget.SetSortingEnabled(true)
	treeWidget.SortByColumn(sectionClicked, qt.SortOrder(0))
	treeWidget.SetAlternatingRowColors(true)
	tableColors := "alternate-background-color: #88DD88; background-color:#FFFFFF; color:#000000; font-size: 12px;"
	treeWidget.SetStyleSheet(tableColors)
	treeWidget.Header()

	treeWidget.SetHeaderLabels([]string{"Channel Name", "Checked", "Downloaded", "Added", "Status", "Feed_CT"})

	//loop through channels
	counter := 0
	for channels.Next() {
		var channel database.Channel
		err := channels.Scan(&channel.ID, &channel.Displayname, &channel.Dldir, &channel.Yt_channelid, &channel.Lastpub, &channel.Lastcheck, &channel.Archive, &channel.Notes, &channel.Date_added, &channel.Last_feed_count)
		if err != nil {
			fmt.Println("something went wrong with the channel scan")
		}

		//filter by will be added
		treewidgetItem := qt.NewQTreeWidgetItem2([]string{channel.Displayname, functions.DateConvertTrim(channel.Lastcheck, 10), functions.DateConvertTrim(database.GetLastDownload(config.Db_name, channel.Yt_channelid), 10), functions.DateConvertTrim(channel.Date_added, 10), database.GetStatus(config.Db_name, strconv.Itoa(channel.Archive)), strconv.Itoa(channel.Last_feed_count)})
		treewidgetItem.SetData(0, int(qt.UserRole), qt.NewQVariant11(channel.Yt_channelid))
		treeWidget.AddTopLevelItem(treewidgetItem)
		counter++
	}

	// treeWidget.OnKeyReleaseEvent(keyPressEvent)
	treeWidget.OnKeyReleaseEvent(func(super func(event *qt.QKeyEvent), event *qt.QKeyEvent) {
		//get selected sub and then pass to the master key event in libs
		index := treeWidget.IndexFromItem(treeWidget.CurrentItem())
		indexSelected = index.Row()
		data := index.DataWithRole(int(qt.UserRole)).ToString()
		GlobalChannelID = data
		chaninfo := database.GetChanInfo(config.Db_name, data)
		Window.StatusBar().ShowMessage("Subscription Selected: " + chaninfo.Displayname + " " + data)

		// keyPressEvent(event, w)
	})

	treeWidget.OnContextMenuEvent(func(super func(event *qt.QContextMenuEvent), event *qt.QContextMenuEvent) {
		contextMenu(GlobalChannelID, event)
	})
	treeWidget.OnClicked(func(index *qt.QModelIndex) {
		indexSelected = index.Row()
		data := index.DataWithRole(int(qt.UserRole)).ToString()
		//set global channel id for subsequent actions
		GlobalChannelID = data
		chaninfo := database.GetChanInfo(config.Db_name, data)
		Window.StatusBar().ShowMessage("Subscription Selected: " + chaninfo.Displayname + " " + data)
		// qt.QMessageBox_Information(nil, "OK", data, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
	})
	treeWidget.OnDoubleClicked(func(index *qt.QModelIndex) {
		indexSelected = index.Row()
		data := index.DataWithRole(int(qt.UserRole)).ToString()
		//set global channel id for subsequent actions
		GlobalChannelID = data
		//double click means open the settings for this channel
		if GlobalChannelID != "" {
			ChannelSettings(GlobalChannelID)
		}
		chaninfo := database.GetChanInfo(config.Db_name, data)
		Window.StatusBar().ShowMessage("Subscription Selected: " + chaninfo.Displayname + " " + data)
		// qt.QMessageBox_Information(nil, "OK", "Open Subscription Settings for "+data, qt.QMessageBox__Ok, qt.QMessageBox__Ok)
	})
	treeWidget.Header().OnSectionClicked(func(logicalIndex int) {
		sectionClicked = logicalIndex

	})
	treeWidget.ResizeColumnToContents(0)
	treeWidget.SetCurrentItem(treeWidget.TopLevelItem(indexSelected))

	treeWidget.ScrollToItem(treeWidget.TopLevelItem(indexSelected))
	//end loop
	verticalLayout.AddWidget(treeWidget.QWidget)
	subCount = counter
	toolbar := toolbarInit(qt.NewQToolBar3())

	toolbar.AddSeparator()

	verticalLayout.SetMenuBar(toolbar.QWidget)
	mainWidget.SetLayout(verticalLayout.QLayout)

	// // Set main widget as the central widget of the window
	Window.SetCentralWidget(mainWidget)

	// // Show the window
	Window.Show()

}
