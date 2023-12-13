package main

import (
	"strconv"

	"github.com/fizzywhizbang/YTGO/database"
	"github.com/fizzywhizbang/YTGO/functions"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

func ChannelSettings(ytchanid string) {

	// Create main window
	window := widgets.NewQMainWindow(nil, 0)
	window.SetWindowTitle("Edit Channel")
	window.SetMinimumSize2(800, 400)
	FormSelected = "EditChannel"

	window.ConnectKeyPressEvent(func(e *gui.QKeyEvent) {
		if int32(e.Key()) == int32(core.Qt__Key_Escape) {
			//close window
			window.Close()
		}
	})

	//get data
	channel := database.GetChanInfo(config.Db_name, ytchanid)

	// Create main widget and set the layout
	mainWidget := widgets.NewQWidget(nil, 0)
	mainWidget.SetContentsMargins(0, 2, 0, 0)

	//crate page for tab

	//create form layout
	mainFormLayout := widgets.NewQFormLayout(nil)

	//create tabbed container
	tabContainer := widgets.NewQTabWidget(nil)
	tabContainer.SetMinimumWidth(790)

	//details
	//Date added
	tableWidget := widgets.NewQTableWidget(nil)
	tableWidget.SetColumnCount(2)
	tableWidget.SetRowCount(6)
	tableWidget.SetHorizontalHeaderLabels([]string{"Title", "Data"})
	tableWidget.SetAlternatingRowColors(true)
	tableColors := "alternate-background-color: #88DD88; background-color:#FFFFFF; color:#000000; font-size: 12px;"
	tableWidget.SetStyleSheet(tableColors)

	tableWidget.SetItem(0, 0, widgets.NewQTableWidgetItem2("Date Added", 0))
	tableWidget.SetItem(0, 1, widgets.NewQTableWidgetItem2(functions.DateConvert(channel.Date_added), 0))

	tableWidget.SetItem(1, 0, widgets.NewQTableWidgetItem2("Last Download", 0))
	tableWidget.SetItem(1, 1, widgets.NewQTableWidgetItem2(functions.DateConvert(channel.Lastpub), 0))

	tableWidget.SetItem(2, 0, widgets.NewQTableWidgetItem2("Last Download", 0))
	tableWidget.SetItem(2, 1, widgets.NewQTableWidgetItem2(strconv.Itoa(channel.Lastpub), 0))

	tableWidget.SetItem(3, 0, widgets.NewQTableWidgetItem2("Last Check", 0))
	tableWidget.SetItem(3, 1, widgets.NewQTableWidgetItem2(functions.DateConvert(channel.Lastcheck), 0))

	tableWidget.SetItem(4, 0, widgets.NewQTableWidgetItem2("Directory", 0))
	tableWidget.SetItem(4, 1, widgets.NewQTableWidgetItem2(channel.Dldir, 0))

	tableWidget.SetItem(5, 0, widgets.NewQTableWidgetItem2("Last Feed Count", 0))
	tableWidget.SetItem(5, 1, widgets.NewQTableWidgetItem2(strconv.Itoa(channel.Last_feed_count), 0))

	// detailsWidget.SetLayout(detailsLayout)
	tableWidget.ResizeColumnToContents(0)
	tableWidget.ResizeColumnToContents(1)
	tabContainer.AddTab(tableWidget, "Details")

	//settings Tab
	settingsWidget := widgets.NewQWidget(nil, 0)
	layout := widgets.NewQFormLayout(nil)
	layout.SetFieldGrowthPolicy(widgets.QFormLayout__ExpandingFieldsGrow)
	// Create a line edit and add it to the layout
	input := widgets.NewQLineEdit(nil)
	input.SetText(ytchanid)

	layout.AddRow3("Channel URL: ", input)

	input2 := widgets.NewQLineEdit(nil)
	input2.SetText(channel.Displayname)
	layout.InsertRow3(1, "Channel Name: ", input2)

	input3 := widgets.NewQLineEdit(nil)
	input3.SetText(channel.Dldir)
	layout.InsertRow3(2, "Directory: ", input3)

	input4 := widgets.NewQTextEdit(nil)
	input4.SetText(channel.Notes)
	layout.InsertRow3(3, "Notes/Tags: ", input4)

	//tags
	tags := database.GetAllTags(config.Db_name, "tag")
	tagSelector := widgets.NewQComboBox(nil)
	tagItems := []string{}
	tagItems = append(tagItems, "Select to add to Notes/Tags box")
	for tags.Next() {
		var tag database.Tags
		err := tags.Scan(&tag.ID, &tag.Name)
		functions.CheckErr(err, "unable to get tags (subsettings.go)")
		tagItems = append(tagItems, tag.Name)
	}

	tagSelector.AddItems(tagItems)
	layout.InsertRow3(4, "Tag Selector: ", tagSelector)

	tagSelector.ConnectCurrentTextChanged(func(text string) {

		if text != "Select to add to Notes/Tags box" {
			//get current text in notes and append to it
			currentNotes := input4.ToPlainText()
			currentNotes += "\n#" + text
			input4.SetText(currentNotes)
		}
	})

	selector := widgets.NewQComboBox(nil)
	statuses := database.GetAllStatus(config.Db_name)
	statusSlice := []string{}
	for statuses.Next() {
		var status database.Category
		err := statuses.Scan(&status.ID, &status.Name)
		statusSlice = append(statusSlice, status.Name)
		functions.CheckErr(err, "Unable to get statuses (subsettings.go)")
	}
	// statuses := []string{"Active", "Archive", "Manual", "Delete", "FVG"}
	selector.AddItems(statusSlice)

	selector.SetCurrentText(database.GetStatus(config.Db_name, strconv.Itoa(channel.Archive)))

	layout.InsertRow3(5, "Status: ", selector)

	optionGroup := widgets.NewQHBoxLayout()
	//save button
	saveButton := widgets.NewQPushButton(nil)
	saveButton.SetText("Save Changes")
	optionGroup.AddWidget(saveButton, 0, 0)
	//cancel button
	cancelButton := widgets.NewQPushButton(nil)
	cancelButton.SetText("Cancel")
	optionGroup.AddWidget(cancelButton, 0, 0)
	//
	optionGroup2 := widgets.NewQHBoxLayout()
	//Download new button
	dlNew := widgets.NewQPushButton(nil)
	dlNew.SetText("Download New")
	optionGroup2.AddWidget(dlNew, 0, 0)
	//update database button
	updateDB := widgets.NewQPushButton(nil)
	updateDB.SetText("Update Database (no dl)")
	optionGroup2.AddWidget(updateDB, 0, 0)

	//goto URL button
	gotoURLButton := widgets.NewQPushButton(nil)
	gotoURLButton.SetText("Go To URL")
	optionGroup2.AddWidget(gotoURLButton, 0, 0)

	//cancel action
	cancelButton.ConnectClicked(func(checked bool) { window.Close() })
	//goto url action
	gotoURLButton.ConnectClicked(func(checked bool) { functions.Openbrowser(config.Defbrowser, channel.Yt_channelid) })

	layout.InsertRow6(6, optionGroup)
	layout.InsertRow6(7, optionGroup2)
	settingsWidget.SetLayout(layout)
	tabContainer.AddTab(settingsWidget, "Settings")

	videoDL := contentListDL(channel.Yt_channelid)
	tabContainer.AddTab(videoDL, "Downloaded Video")

	tabContainer.ConnectCurrentChanged(func(index int) {
		//index2 is downloaded videos
		if index == 2 {
			videoDL.Clear()                            //clear content
			contentFill(channel.Yt_channelid, videoDL) //reload content
		}

	})
	// tabContainer.AddTab(widgets.NewQLabel2("Downloaded Videos", nil, 0), "Downloaded Videos")
	// removing and opting for having a separate feed window
	// rssFeed := contentList(channel.yt_channelid)
	// tabContainer.AddTab(rssFeed, "RSS Feed")

	updateDB.ConnectClicked(func(checked bool) {
		functions.UpdateChan(config.Db_name, config.FolderWatch, channel.Yt_channelid, false, true)
		database.UpdateChecked(config.Db_name, channel.Yt_channelid)
		// feedCheck(channel.yt_channelid)
	})
	dlNew.ConnectClicked(func(checked bool) {
		functions.UpdateChan(config.Db_name, config.FolderWatch, channel.Yt_channelid, true, true)
		database.UpdateChecked(config.Db_name, channel.Yt_channelid)
	})
	saveButton.ConnectClicked(func(checked bool) {
		//fmt.Println(getStatusIDI(selector.CurrentText()))
		//if channel id is changed prompt are you sure
		result := false
		if channel.Yt_channelid != input.Text() {
			action := widgets.QMessageBox_Question(nil, "Warning", "Are you sure you want to update "+channel.Yt_channelid+" to "+input.Text(), widgets.QMessageBox__Yes|widgets.QMessageBox__No, 0)
			if action == widgets.QMessageBox__Yes {
				result = database.ModChanSettings(config.Db_name, channel.Yt_channelid, input.Text(), input2.Text(), input3.Text(), functions.MysqlRealEscapeString(input4.ToPlainText()), database.GetStatusIDI(config.Db_name, selector.CurrentText()))
				//we also need to refresh the view because the channel id changed
				refreshFunc(Window, App)
			}
		} else {
			result = database.ModChanSettings(config.Db_name, channel.Yt_channelid, input.Text(), input2.Text(), input3.Text(), input4.ToPlainText(), database.GetStatusIDI(config.Db_name, selector.CurrentText()))
		}
		if result {
			widgets.QMessageBox_Information(nil, "OK", "Update Complete", widgets.QMessageBox__Ok, widgets.QMessageBox__Cancel)
		} else {
			widgets.QMessageBox_Information(nil, "OK", "Failed to update", widgets.QMessageBox__Ok, widgets.QMessageBox__Cancel)
		}

	})

	tabContainer.AddTab(downloadVideoForm(channel.Displayname, channel.Dldir, channel.Yt_channelid), "Download Video")

	mainFormLayout.AddWidget(tabContainer)

	// Set main widget as the central widget of the window
	mainWidget.SetLayout(mainFormLayout)

	// mainWidget.Layout().QLayoutItem.SetAlignment(core.Qt__AlignLeft)
	window.SetCentralWidget(mainWidget)

	// Show the window
	window.Show()

}
