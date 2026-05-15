package main

import (
	"strconv"

	"github.com/fizzywhizbang/YTGO/database"
	"github.com/fizzywhizbang/YTGO/functions"
	qt "github.com/mappu/miqt/qt6"
)

func toolbarInit(toolbar *qt.QToolBar) *qt.QToolBar {

	statusSelector := qt.NewQComboBox(nil)
	statusitems := []string{"Main"}
	statuses := database.GetAllStatus(config.Db_name)
	//create vertical layout
	for statuses.Next() {
		var status database.Category
		err := statuses.Scan(&status.ID, &status.Name)
		functions.CheckErr(err, "Unable to get status from db (toolbar.go line 23)")

		statusitems = append(statusitems, status.Name)
	}
	statusSelector.AddItems(statusitems)
	toolbar.AddWidget(statusSelector.QWidget)
	if GlobalStatus != "" {
		//set default selector item
		statusSelector.SetCurrentText(database.GetStatus(config.Db_name, GlobalStatus))
	}
	statusSelector.OnCurrentTextChanged(func(text string) {

		if text == "Main" {
			//reset index selected
			indexSelected = 0
			//reset section/header clicked
			sectionClicked = 0
			GlobalStatus = ""
			createHomeWindow()
		} else {
			GlobalStatus := database.GetStatusName(config.Db_name, text)

			if GlobalSearchType != "" && globalSearchTags != "" {
				//continue with current filter
				showSubsSearch(globalSearchTags, GlobalSearchType, GlobalStatus)
			} else {
				GlobalSearchType = "" //resetting because
				showSubs(GlobalStatus)
			}

		}

	})

	toolbar.SetToolButtonStyle(qt.ToolButtonTextOnly)
	toolbar.SetMovable(true)
	//search
	selector := qt.NewQComboBox(nil)
	items := []string{"Select Search Type", "Tags", "Notes", "Channel Name", "Channel Directory", "Channel ID", "Channel with Video Title"}
	selector.AddItems(items)
	if GlobalSearchType != "" {
		selector.SetCurrentText(GlobalSearchType)
	}

	toolbar.AddWidget(selector.QWidget)
	searchTags := qt.NewQLineEdit(nil)
	searchTags.SetMaximumWidth(400)
	if len(globalSearchTags) > 1 {
		searchTags.SetText(globalSearchTags)
	}
	selector.OnCurrentTextChanged(func(text string) {
		switch text {
		case "Select Search Type":
			globalSearchTags = ""
			GlobalSearchType = ""
			//reset subcount
			subCount = 0
			searchTags.SetPlaceholderText("")
			//clear entry
			searchTags.SetText("")
		case "Channel ID":
			searchTags.SetPlaceholderText("UCxxxx")
		default:
			searchTags.SetPlaceholderText("spaces will be anded or use & for and | for or")

		}
		//
	})
	toolbar.AddWidget(searchTags.QWidget)
	searchTags.OnKeyReleaseEvent(func(super func(event *qt.QKeyEvent), event *qt.QKeyEvent) {
		if int32(event.Key()) == int32(qt.Key_Return) || int32(event.Key()) == int32(qt.Key_Enter) {
			if selector.CurrentText() != "Select Search Type" {
				if searchTags.Text() != "" {
					if selector.CurrentText() == "Channel ID" {
						channel := database.GetChanInfo(config.Db_name, searchTags.Text())
						//search multiple criteria
						if channel.Displayname == "" {
							action := qt.QMessageBox_Warning(nil, "Search not found", "There is no channel with ID"+searchTags.Text()+" found\nWould you like to add it?")
							if action == qt.QMessageBox__Ok {
								addChannel(searchTags.Text())
							}
						} else {
							action := qt.QMessageBox_Question(nil, "Channel Exists", "This channel exists do you want to view the settings?")
							if action == qt.QMessageBox__Open {
								ChannelSettings(searchTags.Text())
							}
						}
					} else {
						//set global search tags
						globalSearchTags = searchTags.Text()
						//unset status so we can use the shortcut keys in this view
						GlobalSearchType = selector.CurrentText()

						showSubsSearch(searchTags.Text(), selector.CurrentText(), GlobalStatus)
					}

				}
			}

		}
	})

	tagButton := qt.NewQPushButton3("tags")
	tagButton.OnClicked(func() {
		//open tag window for searching
		showTagSearch()
	})
	toolbar.AddWidget(tagButton.QWidget)

	countLabel := qt.NewQLabel(toolbar.QWidget)
	if subCount > 0 || len(GlobalSearchType) > 1 {
		countLabel.SetText(strconv.Itoa(subCount))
	} else {
		ct := database.CheckCount(config.Db_name, GlobalStatus)
		if ct > 0 {
			countLabel.SetText(strconv.Itoa(ct))
		}

	}

	toolbar.AddWidget(countLabel.QWidget)
	return toolbar
}
