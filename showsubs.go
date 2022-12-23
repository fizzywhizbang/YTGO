package main

import (
	"strconv"

	"github.com/fizzywhizbang/YTGO/database"
	"github.com/fizzywhizbang/YTGO/functions"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

func showSubs(status string) {

	GlobalStatus = status //set global status to what's been selected for sorting later

	verticalLayout := widgets.NewQVBoxLayout()

	mainWidget := widgets.NewQWidget(nil, 0)

	toolbar := toolbarInit(widgets.NewQToolBar2(nil))

	toolbar.AddSeparator()
	//set menubar
	verticalLayout.SetMenuBar(toolbar)

	//this is where it differs
	channels := database.GetChannels(config.Db_name, status, orderby)

	treeWidget := widgets.NewQTreeWidget(nil)
	treeWidget.SetColumnCount(6)
	treeWidget.SetObjectName("treewidget")
	treeWidget.Header().SetStretchLastSection(false)
	treeWidget.Header().SetSectionsClickable(true)
	treeWidget.SetSortingEnabled(true)
	treeWidget.SortByColumn(sectionClicked, core.Qt__SortOrder(0))
	treeWidget.SetAlternatingRowColors(true)
	// treeWidget.SetSelectionMode(widgets.QAbstractItemView__ExtendedSelection)

	tableColors := "alternate-background-color: #88DD88; background-color:#FFFFFF; color:#000000; font-size: 12px;"
	treeWidget.SetStyleSheet(tableColors)
	treeWidget.Header()

	treeWidget.SetHeaderLabels([]string{"Channel Name", "Checked", "Downloaded", "Added", "Status", "Feed_CT"})

	//loop through channels if there are any
	count := database.CheckCount(config.Db_name, GlobalStatus)
	if count >= 1 {

		for channels.Next() {
			var channel database.Channel
			err := channels.Scan(&channel.ID, &channel.Displayname, &channel.Dldir, &channel.Yt_channelid, &channel.Lastpub, &channel.Lastcheck, &channel.Archive, &channel.Notes, &channel.Date_added, &channel.Last_feed_count)
			functions.CheckErr(err, "Unable to retrieve the channels (showsubs.go)")
			//filter by will be added

			treewidgetItem := widgets.NewQTreeWidgetItem2([]string{channel.Displayname, functions.DateConvertTrim(channel.Lastcheck, 16), functions.DateConvertTrim(database.GetLastDownload(config.Db_name, channel.Yt_channelid), 16), functions.DateConvertTrim(channel.Date_added, 10), database.GetStatus(config.Db_name, strconv.Itoa(channel.Archive)), strconv.Itoa(channel.Last_feed_count)}, channel.ID)
			treewidgetItem.SetData(0, int(core.Qt__UserRole), core.NewQVariant12(channel.Yt_channelid))
			treeWidget.AddTopLevelItem(treewidgetItem)

		}

		treeWidget.Header().ConnectSectionClicked(func(logicalIndex int) {
			sectionClicked = logicalIndex

		})

		treeWidget.ConnectKeyReleaseEvent(func(event *gui.QKeyEvent) {
			//get selected sub and then pass to the master key event in libs
			index := treeWidget.IndexFromItem(treeWidget.CurrentItem(), 0)
			indexSelected = index.Row()
			data := index.Data(int(core.Qt__UserRole)).ToString()
			GlobalChannelID = data
			chaninfo := database.GetChanInfo(config.Db_name, data)
			Window.StatusBar().ShowMessage("Subscription Selected: "+chaninfo.Displayname+" "+data, 0)

		})

		treeWidget.ConnectContextMenuEvent(func(event *gui.QContextMenuEvent) {
			contextMenu(GlobalChannelID, event)
		})

		treeWidget.ConnectClicked(func(index *core.QModelIndex) {
			data := index.Data(int(core.Qt__UserRole)).ToString()
			indexSelected = index.Row()
			//set global channel id for subsequent actions
			GlobalChannelID = data
			chaninfo := database.GetChanInfo(config.Db_name, data)
			Window.StatusBar().ShowMessage("Subscription Selected: "+chaninfo.Displayname+" "+data, 0)

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

		})
	}
	treeWidget.ResizeColumnToContents(0)
	treeWidget.ResizeColumnToContents(1)
	treeWidget.ResizeColumnToContents(2)
	treeWidget.SetCurrentItem(treeWidget.TopLevelItem(indexSelected))

	treeWidget.ScrollToItem(treeWidget.TopLevelItem(indexSelected), widgets.QAbstractItemView__PositionAtCenter)
	//end loop
	verticalLayout.AddWidget(treeWidget, 0, 0)

	mainWidget.SetLayout(verticalLayout)

	// // Set main widget as the central widget of the window
	Window.SetCentralWidget(mainWidget)

	// // Show the window
	Window.Show()

}

func contextMenu(chanid string, event *gui.QContextMenuEvent) {

	menu := widgets.NewQMenu(Window)

	menu.AddAction("Refresh View").ConnectTriggered(func(checked bool) {
		showSubs(GlobalStatus)
	})

	menu.AddAction("Download New").ConnectTriggered(func(checked bool) {
		functions.UpdateChan(config.Db_name, config.FolderWatch, chanid, true, true)

	})

	menu.AddAction("Open URL").ConnectTriggered(func(checked bool) {
		functions.Openbrowser(chanid, config.Defbrowser)

	})

	menu.AddAction("Update DB").ConnectTriggered(func(checked bool) {
		ct := functions.UpdateChan(config.Db_name, config.FolderWatch, chanid, false, false)
		widgets.QMessageBox_Information(nil, "Updated Database", strconv.Itoa(ct)+" videos added to database", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
	})
	menu.AddAction("Delete Channel").ConnectTriggered(func(checked bool) {
		if GlobalChannelID != "" {
			channel := database.GetChanInfo(config.Db_name, GlobalChannelID)
			action := widgets.QMessageBox_Question(nil, "Warning", "Are you sure you want to delete "+channel.Displayname+"?", widgets.QMessageBox__Yes|widgets.QMessageBox__No, 0)
			if action == widgets.QMessageBox__Yes {
				database.DeleteChannel(config.Db_name, GlobalChannelID)
				if GlobalStatus == "" && globalSearchTags != "" {
					showSubsSearch(globalSearchTags, GlobalSearchType, GlobalStatus)
				} else {
					showSubs(GlobalStatus)
				}
			}
		} else {
			widgets.QMessageBox_Information(nil, "Oops", "No channel selected.\nSelect the channel name first.", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		}

	})
	menu.AddSeparator()
	statuses := database.GetAllStatus(config.Db_name)

	for statuses.Next() {
		var status database.Category
		err := statuses.Scan(&status.ID, &status.Name)
		functions.CheckErr(err, "Unable to get database status (showsubs.go)")

		//create actions
		menu.AddAction("mv-to->" + status.Name).ConnectTriggered(func(checked bool) {
			database.MoveTo(config.Db_name, chanid, strconv.Itoa(status.ID))
			statusCount := database.CheckCount(config.Db_name, strconv.Itoa(status.ID))
			//refresh view if count for view < 75 and this is because of a sloooooooo refresh if you have a lot of subs
			if statusCount < 75 {
				if GlobalStatus == "" && globalSearchTags != "" {
					showSubsSearch(globalSearchTags, GlobalSearchType, GlobalStatus)
				} else {
					showSubs(GlobalStatus)
				}
			} else {
				action := widgets.QMessageBox_Question(nil, "Notice", "Due to the number of subs in this status refresh will not be automatic\n Do you want to refresh?", widgets.QMessageBox__Yes|widgets.QMessageBox__No, 0)
				if action == widgets.QMessageBox__Yes {
					if GlobalStatus == "" && globalSearchTags != "" {
						showSubsSearch(globalSearchTags, GlobalSearchType, GlobalStatus)
					} else {
						showSubs(GlobalStatus)
					}
				}
			}

		})
	}

	menu.Exec2(event.GlobalPos().QPoint_PTR(), nil)

}
