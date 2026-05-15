package main

import (
	"strconv"

	"github.com/fizzywhizbang/YTGO/database"
	"github.com/fizzywhizbang/YTGO/functions"
	qt "github.com/mappu/miqt/qt6"
)

func showSubs(status string) {

	GlobalStatus = status //set global status to what's been selected for sorting later

	verticalLayout := qt.NewQVBoxLayout(nil)

	mainWidget := qt.NewQWidget(nil)

	toolbar := toolbarInit(qt.NewQToolBar2("toolbar"))

	toolbar.AddSeparator()
	//set menubar
	verticalLayout.SetMenuBar(toolbar.QWidget)

	//this is where it differs
	channels := database.GetChannels(config.Db_name, status, orderby)

	treeWidget := qt.NewQTreeWidget(nil)
	treeWidget.SetColumnCount(6)
	treeWidget.SetObjectName(*qt.NewQAnyStringView3("treewidget"))
	treeWidget.Header().SetStretchLastSection(false)
	treeWidget.Header().SetSectionsClickable(true)
	treeWidget.SetSortingEnabled(true)
	treeWidget.SortByColumn(sectionClicked, qt.SortOrder(0))
	treeWidget.SetAlternatingRowColors(true)
	// treeWidget.SetSelectionMode(qt.QAbstractItemView__ExtendedSelection)

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

			treewidgetItem := qt.NewQTreeWidgetItem2([]string{channel.Displayname, functions.DateConvertTrim(channel.Lastcheck, 16), functions.DateConvertTrim(database.GetLastDownload(config.Db_name, channel.Yt_channelid), 16), functions.DateConvertTrim(channel.Date_added, 10), database.GetStatus(config.Db_name, strconv.Itoa(channel.Archive)), strconv.Itoa(channel.Last_feed_count)})
			treewidgetItem.SetData(0, int(qt.UserRole), qt.NewQVariant11(channel.Yt_channelid))
			treeWidget.AddTopLevelItem(treewidgetItem)

		}

		treeWidget.Header().OnSectionClicked(func(logicalIndex int) {
			sectionClicked = logicalIndex

		})

		treeWidget.OnKeyReleaseEvent(func(super func(event *qt.QKeyEvent), event *qt.QKeyEvent) {
			//get selected sub and then pass to the master key event in libs
			index := treeWidget.IndexFromItem(treeWidget.CurrentItem())
			indexSelected = index.Row()
			data := index.DataWithRole(int(qt.UserRole)).ToString()
			GlobalChannelID = data
			chaninfo := database.GetChanInfo(config.Db_name, data)
			Window.StatusBar().ShowMessage("Subscription Selected: " + chaninfo.Displayname + " " + data)

		})

		treeWidget.OnContextMenuEvent(func(super func(event *qt.QContextMenuEvent), event *qt.QContextMenuEvent) {
			contextMenu(GlobalChannelID, event)
		})

		treeWidget.OnClicked(func(index *qt.QModelIndex) {
			data := index.DataWithRole(int(qt.UserRole)).ToString()
			indexSelected = index.Row()
			//set global channel id for subsequent actions
			GlobalChannelID = data
			chaninfo := database.GetChanInfo(config.Db_name, data)
			Window.StatusBar().ShowMessage("Subscription Selected: " + chaninfo.Displayname + " " + data)

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

		})
	}
	treeWidget.ResizeColumnToContents(0)
	treeWidget.ResizeColumnToContents(1)
	treeWidget.ResizeColumnToContents(2)
	treeWidget.SetCurrentItem(treeWidget.TopLevelItem(indexSelected))

	treeWidget.ScrollToItem(treeWidget.TopLevelItem(indexSelected))
	//end loop
	verticalLayout.AddWidget(treeWidget.QWidget)

	mainWidget.SetLayout(verticalLayout.QLayout)

	// // Set main widget as the central widget of the window
	Window.SetCentralWidget(mainWidget)

	// // Show the window
	Window.Show()

}

func contextMenu(chanid string, event *qt.QContextMenuEvent) {

	menu := qt.NewQMenu(Window.QWidget)

	menu.AddActionWithText("Refresh View").OnTriggered(func() {
		showSubs(GlobalStatus)
	})

	menu.AddActionWithText("Download New").OnTriggered(func() {
		functions.UpdateChan(config.Db_name, config.FolderWatch, chanid, true, true)

	})

	menu.AddActionWithText("Open URL").OnTriggered(func() {
		functions.Openbrowser(chanid, config.Defbrowser)

	})

	menu.AddActionWithText("Update DB").OnTriggered(func() {
		ct := functions.UpdateChan(config.Db_name, config.FolderWatch, chanid, false, false)
		qt.QMessageBox_Information(nil, "Updated Database", strconv.Itoa(ct)+" videos added to database")
	})
	menu.AddActionWithText("Delete Channel").OnTriggered(func() {
		if GlobalChannelID != "" {
			channel := database.GetChanInfo(config.Db_name, GlobalChannelID)
			action := qt.QMessageBox_Question(nil, "Warning", "Are you sure you want to delete "+channel.Displayname+"?")
			if action == qt.QMessageBox__Yes {
				database.DeleteChannel(config.Db_name, GlobalChannelID)
				if GlobalStatus == "" && globalSearchTags != "" {
					showSubsSearch(globalSearchTags, GlobalSearchType, GlobalStatus)
				} else {
					showSubs(GlobalStatus)
				}
			}
		} else {
			qt.QMessageBox_Information(nil, "Oops", "No channel selected.\nSelect the channel name first.")
		}

	})
	menu.AddSeparator()
	statuses := database.GetAllStatus(config.Db_name)

	for statuses.Next() {
		var status database.Category
		err := statuses.Scan(&status.ID, &status.Name)
		functions.CheckErr(err, "Unable to get database status (showsubs.go)")

		//create actions
		menu.AddActionWithText("mv-to->" + status.Name).OnTriggered(func() {
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
				action := qt.QMessageBox_Question(nil, "Notice", "Due to the number of subs in this status refresh will not be automatic\n Do you want to refresh?")
				if action == qt.QMessageBox__Yes {
					if GlobalStatus == "" && globalSearchTags != "" {
						showSubsSearch(globalSearchTags, GlobalSearchType, GlobalStatus)
					} else {
						showSubs(GlobalStatus)
					}
				}
			}

		})
	}

	menu.ExecWithPos(event.GlobalPos().ToPointF().ToPoint())

}
