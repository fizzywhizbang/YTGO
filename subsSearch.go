package main

import (
	"fmt"
	"strconv"

	"github.com/fizzywhizbang/YTGO/database"
	"github.com/fizzywhizbang/YTGO/functions"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

func showSubsSearch(searchstring string, searchType string, status string) {
	GlobalStatus = status

	verticalLayout := widgets.NewQVBoxLayout()

	mainWidget := widgets.NewQWidget(nil, 0)

	//this is where it differs
	channels := database.ChannelSearch(config.Db_name, GlobalStatus, searchstring, searchType)
	treeWidget := widgets.NewQTreeWidget(nil)
	treeWidget.SetColumnCount(6)
	treeWidget.SetObjectName("treewidget")
	treeWidget.Header().SetStretchLastSection(false)
	treeWidget.Header().SetSectionsClickable(true)
	treeWidget.SetSortingEnabled(true)
	treeWidget.SortByColumn(sectionClicked, core.Qt__SortOrder(0))
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
		treewidgetItem := widgets.NewQTreeWidgetItem2([]string{channel.Displayname, functions.DateConvertTrim(channel.Lastcheck, 10), functions.DateConvertTrim(database.GetLastDownload(config.Db_name, channel.Yt_channelid), 10), functions.DateConvertTrim(channel.Date_added, 10), database.GetStatus(config.Db_name, strconv.Itoa(channel.Archive)), strconv.Itoa(channel.Last_feed_count)}, channel.ID)
		treewidgetItem.SetData(0, int(core.Qt__UserRole), core.NewQVariant12(channel.Yt_channelid))
		treeWidget.AddTopLevelItem(treewidgetItem)
		counter++
	}

	// treeWidget.ConnectKeyReleaseEvent(keyPressEvent)
	treeWidget.ConnectKeyReleaseEvent(func(event *gui.QKeyEvent) {
		//get selected sub and then pass to the master key event in libs
		index := treeWidget.IndexFromItem(treeWidget.CurrentItem(), 0)
		indexSelected = index.Row()
		data := index.Data(int(core.Qt__UserRole)).ToString()
		GlobalChannelID = data
		chaninfo := database.GetChanInfo(config.Db_name, data)
		Window.StatusBar().ShowMessage("Subscription Selected: "+chaninfo.Displayname+" "+data, 0)

		// keyPressEvent(event, w)
	})

	treeWidget.ConnectContextMenuEvent(func(event *gui.QContextMenuEvent) {
		contextMenu(GlobalChannelID, event)
	})
	treeWidget.ConnectClicked(func(index *core.QModelIndex) {
		indexSelected = index.Row()
		data := index.Data(int(core.Qt__UserRole)).ToString()
		//set global channel id for subsequent actions
		GlobalChannelID = data
		chaninfo := database.GetChanInfo(config.Db_name, data)
		Window.StatusBar().ShowMessage("Subscription Selected: "+chaninfo.Displayname+" "+data, 0)
		// widgets.QMessageBox_Information(nil, "OK", data, widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
	})
	treeWidget.ConnectDoubleClicked(func(index *core.QModelIndex) {
		indexSelected = index.Row()
		data := index.Data(int(core.Qt__UserRole)).ToString()
		//set global channel id for subsequent actions
		GlobalChannelID = data
		//double click means open the settings for this channel
		if GlobalChannelID != "" {
			ChannelSettings(GlobalChannelID)
		}
		chaninfo := database.GetChanInfo(config.Db_name, data)
		Window.StatusBar().ShowMessage("Subscription Selected: "+chaninfo.Displayname+" "+data, 0)
		// widgets.QMessageBox_Information(nil, "OK", "Open Subscription Settings for "+data, widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
	})
	treeWidget.Header().ConnectSectionClicked(func(logicalIndex int) {
		sectionClicked = logicalIndex

	})
	treeWidget.ResizeColumnToContents(0)
	treeWidget.SetCurrentItem(treeWidget.TopLevelItem(indexSelected))

	treeWidget.ScrollToItem(treeWidget.TopLevelItem(indexSelected), widgets.QAbstractItemView__PositionAtCenter)
	//end loop
	verticalLayout.AddWidget(treeWidget, 0, 0)
	subCount = counter
	toolbar := toolbarInit(widgets.NewQToolBar2(nil))

	toolbar.AddSeparator()

	verticalLayout.SetMenuBar(toolbar)
	mainWidget.SetLayout(verticalLayout)

	// // Set main widget as the central widget of the window
	Window.SetCentralWidget(mainWidget)

	// // Show the window
	Window.Show()

}
