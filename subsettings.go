package main

import (
	"strconv"

	"github.com/fizzywhizbang/YTGO/database"
	"github.com/fizzywhizbang/YTGO/functions"
	qt "github.com/mappu/miqt/qt6"
)

func ChannelSettings(ytchanid string) {

	// Create main window
	window := qt.NewQMainWindow(nil)
	window.SetWindowTitle("Edit Channel")
	window.SetMinimumSize2(800, 400)
	FormSelected = "EditChannel"

	window.OnKeyPressEvent(func(super func(event *qt.QKeyEvent), event *qt.QKeyEvent) {
		if int32(event.Key()) == int32(qt.Key_Escape) {
			//close window
			window.Close()
		}
	})

	//get data
	channel := database.GetChanInfo(config.Db_name, ytchanid)

	// Create main widget and set the layout
	mainWidget := qt.NewQWidget(nil)
	mainWidget.SetContentsMargins(0, 2, 0, 0)

	//crate page for tab

	//create form layout
	mainFormLayout := qt.NewQFormLayout(nil)

	//create tabbed container
	tabContainer := qt.NewQTabWidget(nil)
	tabContainer.SetMinimumWidth(790)

	//details
	//Date added
	tableWidget := qt.NewQTableWidget(nil)
	tableWidget.SetColumnCount(2)
	tableWidget.SetRowCount(6)
	tableWidget.SetHorizontalHeaderLabels([]string{"Title", "Data"})
	tableWidget.SetAlternatingRowColors(true)
	tableColors := "alternate-background-color: #88DD88; background-color:#FFFFFF; color:#000000; font-size: 12px;"
	tableWidget.SetStyleSheet(tableColors)

	tableWidget.SetItem(0, 0, qt.NewQTableWidgetItem2("Date Added"))
	tableWidget.SetItem(0, 1, qt.NewQTableWidgetItem2(functions.DateConvert(channel.Date_added)))

	tableWidget.SetItem(1, 0, qt.NewQTableWidgetItem2("Last Download"))
	tableWidget.SetItem(1, 1, qt.NewQTableWidgetItem2(functions.DateConvert(channel.Lastpub)))

	tableWidget.SetItem(2, 0, qt.NewQTableWidgetItem2("Last Download"))
	tableWidget.SetItem(2, 1, qt.NewQTableWidgetItem2(strconv.Itoa(channel.Lastpub)))

	tableWidget.SetItem(3, 0, qt.NewQTableWidgetItem2("Last Check"))
	tableWidget.SetItem(3, 1, qt.NewQTableWidgetItem2(functions.DateConvert(channel.Lastcheck)))

	tableWidget.SetItem(4, 0, qt.NewQTableWidgetItem2("Directory"))
	tableWidget.SetItem(4, 1, qt.NewQTableWidgetItem2(channel.Dldir))

	tableWidget.SetItem(5, 0, qt.NewQTableWidgetItem2("Last Feed Count"))
	tableWidget.SetItem(5, 1, qt.NewQTableWidgetItem2(strconv.Itoa(channel.Last_feed_count)))

	// detailsWidget.SetLayout(detailsLayout)
	tableWidget.ResizeColumnToContents(0)
	tableWidget.ResizeColumnToContents(1)
	tabContainer.AddTab(tableWidget.QWidget, "Details")

	//settings Tab
	settingsWidget := qt.NewQWidget(nil)
	layout := qt.NewQFormLayout(nil)
	layout.SetFieldGrowthPolicy(qt.QFormLayout__ExpandingFieldsGrow)
	// Create a line edit and add it to the layout
	input := qt.NewQLineEdit(nil)
	input.SetText(ytchanid)

	layout.AddRow3("Channel URL: ", input.QWidget)

	input2 := qt.NewQLineEdit(nil)
	input2.SetText(channel.Displayname)
	layout.InsertRow3(1, "Channel Name: ", input2.QWidget)

	input3 := qt.NewQLineEdit(nil)
	input3.SetText(channel.Dldir)
	layout.InsertRow3(2, "Directory: ", input3.QWidget)

	input4 := qt.NewQTextEdit(nil)
	input4.SetText(channel.Notes)
	layout.InsertRow3(3, "Notes/Tags: ", input4.QWidget)

	//tags
	tags := database.GetAllTags(config.Db_name, "tag")
	tagSelector := qt.NewQComboBox(nil)
	tagItems := []string{}
	tagItems = append(tagItems, "Select to add to Notes/Tags box")
	for tags.Next() {
		var tag database.Tags
		err := tags.Scan(&tag.ID, &tag.Name)
		functions.CheckErr(err, "unable to get tags (subsettings.go)")
		tagItems = append(tagItems, tag.Name)
	}

	tagSelector.AddItems(tagItems)
	layout.InsertRow3(4, "Tag Selector: ", tagSelector.QWidget)

	tagSelector.OnCurrentTextChanged(func(text string) {

		if text != "Select to add to Notes/Tags box" {
			//get current text in notes and append to it
			currentNotes := input4.ToPlainText()
			currentNotes += "\n#" + text
			input4.SetText(currentNotes)
		}
	})

	selector := qt.NewQComboBox(nil)
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

	layout.InsertRow3(5, "Status: ", selector.QWidget)

	optionGroup := qt.NewQHBoxLayout(nil)
	//save button
	saveButton := qt.NewQPushButton(nil)
	saveButton.SetText("Save Changes")
	optionGroup.AddWidget(saveButton.QWidget)
	//cancel button
	cancelButton := qt.NewQPushButton(nil)
	cancelButton.SetText("Cancel")
	optionGroup.AddWidget(cancelButton.QWidget)
	//
	optionGroup2 := qt.NewQHBoxLayout(nil)
	//Download new button
	dlNew := qt.NewQPushButton(nil)
	dlNew.SetText("Download New")
	optionGroup2.AddWidget(dlNew.QWidget)
	//update database button
	updateDB := qt.NewQPushButton(nil)
	updateDB.SetText("Update Database (no dl)")
	optionGroup2.AddWidget(updateDB.QWidget)

	//goto URL button
	gotoURLButton := qt.NewQPushButton(nil)
	gotoURLButton.SetText("Go To URL")
	optionGroup2.AddWidget(gotoURLButton.QWidget)

	//cancel action
	cancelButton.OnClicked(func() { window.Close() })
	//goto url action
	gotoURLButton.OnClicked(func() { functions.Openbrowser(config.Defbrowser, channel.Yt_channelid) })

	layout.InsertRow6(6, optionGroup.QLayout)
	layout.InsertRow6(7, optionGroup2.QLayout)
	settingsWidget.SetLayout(layout.QLayout)
	tabContainer.AddTab(settingsWidget, "Settings")

	videoDL := contentListDL(channel.Yt_channelid)
	tabContainer.AddTab(videoDL.QWidget, "Downloaded Video")

	tabContainer.OnCurrentChanged(func(index int) {
		//index2 is downloaded videos
		if index == 2 {
			videoDL.Clear()                            //clear content
			contentFill(channel.Yt_channelid, videoDL) //reload content
		}

	})
	// tabContainer.AddTab(qt.NewQLabel2("Downloaded Videos", nil, 0), "Downloaded Videos")
	// removing and opting for having a separate feed window
	// rssFeed := contentList(channel.yt_channelid)
	// tabContainer.AddTab(rssFeed, "RSS Feed")

	updateDB.OnClicked(func() {
		functions.UpdateChan(config.Db_name, config.FolderWatch, channel.Yt_channelid, false, true)
		database.UpdateChecked(config.Db_name, channel.Yt_channelid)
		// feedCheck(channel.yt_channelid)
	})
	dlNew.OnClicked(func() {
		functions.UpdateChan(config.Db_name, config.FolderWatch, channel.Yt_channelid, true, true)
		database.UpdateChecked(config.Db_name, channel.Yt_channelid)
	})
	saveButton.OnClicked(func() {
		//fmt.Println(getStatusIDI(selector.CurrentText()))
		//if channel id is changed prompt are you sure
		result := false
		if channel.Yt_channelid != input.Text() {
			action := qt.QMessageBox_Question(nil, "Warning", "Are you sure you want to update "+channel.Yt_channelid+" to "+input.Text())
			if action == qt.QMessageBox__Yes {
				result = database.ModChanSettings(config.Db_name, channel.Yt_channelid, input.Text(), input2.Text(), input3.Text(), functions.MysqlRealEscapeString(input4.ToPlainText()), database.GetStatusIDI(config.Db_name, selector.CurrentText()))
				//we also need to refresh the view because the channel id changed
				refreshFunc(Window, App)
			}
		} else {
			result = database.ModChanSettings(config.Db_name, channel.Yt_channelid, input.Text(), input2.Text(), input3.Text(), input4.ToPlainText(), database.GetStatusIDI(config.Db_name, selector.CurrentText()))
		}
		if result {
			qt.QMessageBox_Information(nil, "OK", "Update Complete")
		} else {
			qt.QMessageBox_Information(nil, "OK", "Failed to update")
		}

	})

	tabContainer.AddTab(downloadVideoForm(channel.Displayname, channel.Dldir, channel.Yt_channelid), "Download Video")

	mainFormLayout.AddWidget(tabContainer.QWidget)

	// Set main widget as the central widget of the window
	mainWidget.SetLayout(mainFormLayout.QLayout)

	// mainWidget.Layout().QLayoutItem.SetAlignment(qt.AlignLeft)
	window.SetCentralWidget(mainWidget)

	// Show the window
	window.Show()

}
