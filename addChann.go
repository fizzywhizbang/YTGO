package main

import (
	"strconv"
	"time"

	"github.com/fizzywhizbang/YTGO/database"
	"github.com/fizzywhizbang/YTGO/functions"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

func addChannel(chanid string) {
	config := ConfigParser()
	// Create main window
	window := widgets.NewQMainWindow(nil, 0)
	window.SetWindowTitle("Add Channel")
	window.SetMinimumSize2(800, 400)

	window.ConnectKeyPressEvent(func(e *gui.QKeyEvent) {
		if int32(e.Key()) == int32(core.Qt__Key_Escape) {
			//close window
			window.Close()
		}
	})
	// Create main widget and set the layout
	mainWidget := widgets.NewQWidget(nil, 0)
	mainWidget.SetContentsMargins(0, 2, 0, 0)

	//create form layout
	layout := widgets.NewQFormLayout(nil)

	layout.SetFieldGrowthPolicy(widgets.QFormLayout__ExpandingFieldsGrow)

	// Create a line edit and add it to the layout
	input := widgets.NewQLineEdit(nil)
	input.SetPlaceholderText("UCxxxx")
	input.SetText(chanid)
	input.SetToolTip("Press enter to lookup channel information")
	layout.AddRow3("Channel URL: ", input)

	input2 := widgets.NewQLineEdit(nil)
	layout.InsertRow3(1, "Channel Name: ", input2)

	label := widgets.NewQLabel(nil, 0)
	label.SetText(config.BaseDL)
	layout.InsertRow3(2, "Base Directory: ", label)

	input3 := widgets.NewQLineEdit(nil)
	layout.InsertRow3(3, "Directory: ", input3)

	input4 := widgets.NewQTextEdit(nil)
	layout.InsertRow3(4, "Notes/Tags: ", input4)

	input.ConnectKeyReleaseEvent(func(event *gui.QKeyEvent) {
		if int32(event.Key()) == int32(core.Qt__Key_Return) || int32(event.Key()) == int32(core.Qt__Key_Enter) {
			chaninfo := functions.GetChanInfoFromYT(input.Text())
			input2.SetText(chaninfo.Displayname)
			dir := chaninfo.Displayname
			input3.SetText(dir)
			input4.SetText(chaninfo.Notes)
			//see if channel exists
			if database.GetChanExist(config.Db_name, input.Text()) == 1 {
				action := widgets.QMessageBox_Question(nil, "Channel Exists", "This channel exists do you want to view the settings?", widgets.QMessageBox__Open|widgets.QMessageBox__Cancel, 0)

				if action == widgets.QMessageBox__Open {
					ChannelSettings(input.Text())
					window.Close()
				}
			}
		}
	})

	//tags
	tags := database.GetAllTags(config.Db_name, "tag")
	tagSelector := widgets.NewQComboBox(nil)
	tagSelector.SetToolTip("Under System->Edit Tags you can add more tags")
	tagItems := []string{}
	tagItems = append(tagItems, "Select to add to Notes/Tags box")
	for tags.Next() {
		var tag database.Tags
		err := tags.Scan(&tag.ID, &tag.Name)
		functions.CheckErr(err, "Unable to get tags from database (addChann.go)")
		tagItems = append(tagItems, tag.Name)
	}

	tagSelector.AddItems(tagItems)
	layout.InsertRow3(5, "Tag Selector: ", tagSelector)

	tagSelector.ConnectCurrentTextChanged(func(text string) {

		if text != "Select to add to Notes/Tags box" {
			//get current text in notes and append to it
			currentNotes := input4.ToPlainText()
			currentNotes += "\n#" + text
			input4.SetText(currentNotes)
			//since I have selected the tag selector let's keep focus on it for faster tag selection
			tagSelector.SetFocus2()
		}
	})

	//selector needs to be generated from the database
	statuses := database.GetAllStatus(config.Db_name)
	selector := widgets.NewQComboBox(nil)
	items := []string{}
	for statuses.Next() {
		var status database.Category
		err := statuses.Scan(&status.ID, &status.Name)
		functions.CheckErr(err, "Unable to get statuses from database (addchann.go)")
		items = append(items, status.Name)
	}

	selector.AddItems(items)
	layout.InsertRow3(6, "Status: ", selector)

	optionGroup := widgets.NewQHBoxLayout()
	//mark all downloaded
	markAll := widgets.NewQCheckBox(nil)
	markAll.SetText("Mark All Downloaded")
	optionGroup.AddWidget(markAll, 0, 0)
	//Download all Videos
	downloadAll := widgets.NewQCheckBox(nil)
	downloadAll.SetText("Download All")
	optionGroup.AddWidget(downloadAll, 0, 0)
	//View Settings
	viewSettings := widgets.NewQCheckBox(nil)
	viewSettings.SetText("View Settings")
	optionGroup.AddWidget(viewSettings, 0, 0)
	//add button
	addButton := widgets.NewQPushButton(nil)
	addButton.SetText("Add")
	optionGroup.AddWidget(addButton, 0, 0)
	//cancel button
	cancelButton := widgets.NewQPushButton(nil)
	cancelButton.SetText("Cancel")
	optionGroup.AddWidget(cancelButton, 0, 0)

	cancelButton.ConnectClicked(func(checked bool) { window.Close() })

	//progress bar
	progressBar := widgets.NewQProgressBar(nil)
	progressBar.SetMinimum(0)
	progressBar.SetMaximum(100)
	// progressBar.SetValue(progressBar.Maximum() / 2)
	layout.InsertRow5(8, progressBar)

	addButton.ConnectClicked(func(checked bool) {
		//add channel
		var channel database.Channel
		//displayname, dldir, yt_channelid, lastcheck, archive, notes, date_added
		channel.Yt_channelid = input.Text()
		channel.Displayname = input2.Text()
		channel.Dldir = input3.Text()
		channel.Notes = functions.MysqlRealEscapeString(input4.ToPlainText())
		channel.Archive = database.GetStatusIDI(config.Db_name, selector.CurrentText())
		channel.Lastcheck = int(time.Now().Unix())
		channel.Date_added = int(time.Now().Unix())
		progressBar.SetValue(progressBar.Maximum() / 3)
		result := database.InsertChannel(config.Db_name, channel)
		ct := 0
		if result {

			//do action on marking and downloading
			if markAll.IsChecked() && !downloadAll.IsChecked() {
				ct = functions.UpdateChan(config.Db_name, config.FolderWatch, channel.Yt_channelid, false, false)
				progressBar.SetValue(progressBar.Maximum() / 2)
			}
			if downloadAll.IsChecked() {
				//ignore markall and download all from feed
				ct = functions.UpdateChan(config.Db_name, config.FolderWatch, channel.Yt_channelid, true, false)
				progressBar.SetValue(progressBar.Maximum() / 2)
			}
			//if view settings open settings window after closing this one
			if viewSettings.IsChecked() {
				window.Close()
				GlobalChannelID = channel.Yt_channelid
				time.Sleep(2 * time.Second)
				progressBar.SetValue(progressBar.Maximum())
				ChannelSettings(GlobalChannelID)
			}
			progressBar.SetValue(progressBar.Maximum())
			time.Sleep(time.Second)
			widgets.QMessageBox_Information(nil, "OK", "Added "+input2.Text()+" with "+strconv.Itoa(ct)+" videos", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		} else {
			widgets.QMessageBox_Warning(nil, "Warning", "Something went wrong and I can't handle it", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		}

	})

	layout.InsertRow6(7, optionGroup)

	instructionsLabel := widgets.NewQLabel2("After inserting the Channel ID press enter and I'll fetch the channel details", nil, 0)
	instructionsLabel.SetFont(gui.NewQFont2("Times", 12, 1, true))
	layout.AddRow5(instructionsLabel)

	// Set main widget as the central widget of the window
	mainWidget.SetLayout(layout)
	// mainWidget.Layout().QLayoutItem.SetAlignment(core.Qt__AlignLeft)
	window.SetCentralWidget(mainWidget)

	// Show the window
	window.Show()

}
