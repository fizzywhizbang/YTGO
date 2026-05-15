package main

import (
	"strconv"
	"time"

	"github.com/fizzywhizbang/YTGO/database"
	"github.com/fizzywhizbang/YTGO/functions"
	qt "github.com/mappu/miqt/qt6"
)

func addChannel(chanid string) {
	config := ConfigParser()
	// Create main window
	window := qt.NewQMainWindow(nil)
	window.SetWindowTitle("Add Channel")
	window.SetMinimumSize2(800, 400)

	window.OnKeyPressEvent(func(super func(*qt.QKeyEvent), e *qt.QKeyEvent) {
		if int32(e.Key()) == int32(qt.Key_Escape) {
			//close window
			window.Close()
		}
	})
	// Create main widget and set the layout
	mainWidget := qt.NewQWidget(nil)
	mainWidget.SetContentsMargins(0, 2, 0, 0)

	//create form layout
	layout := qt.NewQFormLayout(nil)

	layout.SetFieldGrowthPolicy(qt.QFormLayout__ExpandingFieldsGrow)

	// Create a line edit and add it to the layout
	input := qt.NewQLineEdit(nil)
	input.SetPlaceholderText("UCxxxx")
	input.SetText(chanid)
	input.SetToolTip("Press enter to lookup channel information")
	layout.AddRow3("Channel URL: ", input.QWidget)

	input2 := qt.NewQLineEdit(nil)
	layout.InsertRow3(1, "Channel Name: ", input2.QWidget)

	label := qt.NewQLabel(nil)
	label.SetText(config.BaseDL)
	layout.InsertRow3(2, "Base Directory: ", label.QWidget)

	input3 := qt.NewQLineEdit(nil)
	layout.InsertRow3(3, "Directory: ", input3.QWidget)

	input4 := qt.NewQTextEdit(nil)
	layout.InsertRow3(4, "Notes/Tags: ", input4.QWidget)

	input.OnKeyReleaseEvent(func(super func(*qt.QKeyEvent), event *qt.QKeyEvent) {
		if int32(event.Key()) == int32(qt.Key_Return) || int32(event.Key()) == int32(qt.Key_Enter) {
			chaninfo := functions.GetChanInfoFromYT(input.Text())
			input2.SetText(chaninfo.Displayname)
			dir := chaninfo.Displayname
			input3.SetText(dir)
			input4.SetText(chaninfo.Notes)
			//see if channel exists
			if database.GetChanExist(config.Db_name, input.Text()) == 1 {
				action := qt.QMessageBox_Question2(nil, "Channel Exists", "This channel exists do you want to view the settings?", qt.QMessageBox__Open|qt.QMessageBox__Cancel, qt.QMessageBox__Cancel)

				if action == int(qt.QMessageBox__Open) {
					ChannelSettings(input.Text())
					window.Close()
				}
			}
		}
	})

	//tags
	tags := database.GetAllTags(config.Db_name, "tag")
	tagSelector := qt.NewQComboBox(nil)
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
	layout.InsertRow3(5, "Tag Selector: ", tagSelector.QWidget)

	tagSelector.OnCurrentTextChanged(func(text string) {

		if text != "Select to add to Notes/Tags box" {
			//get current text in notes and append to it
			currentNotes := input4.ToPlainText()
			currentNotes += "\n#" + text
			input4.SetText(currentNotes)
			//since I have selected the tag selector let's keep focus on it for faster tag selection
			tagSelector.SetFocus()
		}
	})

	//selector needs to be generated from the database
	statuses := database.GetAllStatus(config.Db_name)
	selector := qt.NewQComboBox(nil)
	items := []string{}
	for statuses.Next() {
		var status database.Category
		err := statuses.Scan(&status.ID, &status.Name)
		functions.CheckErr(err, "Unable to get statuses from database (addchann.go)")
		items = append(items, status.Name)
	}

	selector.AddItems(items)
	layout.InsertRow3(6, "Status: ", selector.QWidget)

	optionGroup := qt.NewQHBoxLayout(nil)
	//mark all downloaded
	markAll := qt.NewQCheckBox(nil)
	markAll.SetText("Mark All Downloaded")
	optionGroup.AddWidget(markAll.QWidget)
	//Download all Videos
	downloadAll := qt.NewQCheckBox(nil)
	downloadAll.SetText("Download All")
	optionGroup.AddWidget(downloadAll.QWidget)
	//View Settings
	viewSettings := qt.NewQCheckBox(nil)
	viewSettings.SetText("View Settings")
	optionGroup.AddWidget(viewSettings.QWidget)
	//add button
	addButton := qt.NewQPushButton(nil)
	addButton.SetText("Add")
	optionGroup.AddWidget(addButton.QWidget)
	//cancel button
	cancelButton := qt.NewQPushButton(nil)
	cancelButton.SetText("Cancel")
	optionGroup.AddWidget(cancelButton.QWidget)

	cancelButton.OnClicked(func() { window.Close() })

	//progress bar
	progressBar := qt.NewQProgressBar(nil)
	progressBar.SetMinimum(0)
	progressBar.SetMaximum(100)
	// progressBar.SetValue(progressBar.Maximum() / 2)
	layout.InsertRow5(8, progressBar.QWidget)

	addButton.OnClicked(func() {
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
			qt.QMessageBox_Information(nil, "OK", "Added "+input2.Text()+" with "+strconv.Itoa(ct)+" videos")
		} else {
			qt.QMessageBox_Warning(nil, "Warning", "Something went wrong and I can't handle it")
		}

	})

	layout.InsertRow6(7, optionGroup.QLayout)

	instructionsLabel := qt.NewQLabel3("After inserting the Channel ID press enter and I'll fetch the channel details")
	instructionsLabel.SetFont(qt.NewQFont6("Times", 12))
	layout.AddRowWithWidget(instructionsLabel.QWidget)
	// Set main widget as the central widget of the window
	mainWidget.SetLayout(layout.QLayout)
	// mainWidget.Layout().QLayoutItem.SetAlignment(qt.Qt__AlignLeft)
	window.SetCentralWidget(mainWidget)

	// Show the window
	window.Show()

}
