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

func toolbarInit(toolbar *widgets.QToolBar) *widgets.QToolBar {

	statusSelector := widgets.NewQComboBox(nil)
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
	toolbar.AddWidget(statusSelector)
	if GlobalStatus != "" {
		//set default selector item
		statusSelector.SetCurrentText(database.GetStatus(config.Db_name, GlobalStatus))
	}
	statusSelector.ConnectCurrentTextChanged(func(text string) {

		if text == "Main" {
			//reset index selected
			indexSelected = 0
			//reset section/header clicked
			sectionClicked = 0
			GlobalStatus = ""
			createHomeWindow()
		} else {
			GlobalStatus := database.GetStatusName(config.Db_name, text)
			fmt.Println(GlobalStatus)
			if GlobalSearchType != "" && globalSearchTags != "" {
				//continue with current filter
				showSubsSearch(globalSearchTags, GlobalSearchType, GlobalStatus)
			} else {
				GlobalSearchType = "" //resetting because
				showSubs(GlobalStatus)
			}

		}

	})

	toolbar.SetToolButtonStyle(core.Qt__ToolButtonTextOnly)
	toolbar.SetMovable(true)
	//search
	selector := widgets.NewQComboBox(nil)
	items := []string{"Select Search Type", "Tags", "Notes", "Channel Name", "Channel Directory", "Channel ID", "Channel with Video Title"}
	selector.AddItems(items)
	if GlobalSearchType != "" {
		selector.SetCurrentText(GlobalSearchType)
	}

	toolbar.AddWidget(selector)
	searchTags := widgets.NewQLineEdit(nil)
	searchTags.SetMaximumWidth(400)
	if len(globalSearchTags) > 1 {
		searchTags.SetText(globalSearchTags)
	}
	selector.ConnectCurrentTextChanged(func(text string) {
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
	toolbar.AddWidget(searchTags)
	searchTags.ConnectKeyReleaseEvent(func(event *gui.QKeyEvent) {
		if int32(event.Key()) == int32(core.Qt__Key_Return) || int32(event.Key()) == int32(core.Qt__Key_Enter) {
			if selector.CurrentText() != "Select Search Type" {
				if searchTags.Text() != "" {
					if selector.CurrentText() == "Channel ID" {
						channel := database.GetChanInfo(config.Db_name, searchTags.Text())
						//search multiple criteria
						if channel.Displayname == "" {
							action := widgets.QMessageBox_Warning(nil, "Search not found", "There is no channel with ID"+searchTags.Text()+" found\nWould you like to add it?", widgets.QMessageBox__Ok, widgets.QMessageBox__Cancel)
							if action == widgets.QMessageBox__Ok {
								addChannel(searchTags.Text())
							}
						} else {
							action := widgets.QMessageBox_Question(nil, "Channel Exists", "This channel exists do you want to view the settings?", widgets.QMessageBox__Open|widgets.QMessageBox__Cancel, 0)
							if action == widgets.QMessageBox__Open {
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

	tagButton := widgets.NewQPushButton2("tags", nil)
	tagButton.ConnectClicked(func(checked bool) {
		//open tag window for searching
		showTagSearch()
	})
	toolbar.AddWidget(tagButton)

	countLabel := widgets.NewQLabel(toolbar, 0)
	if subCount > 0 || len(GlobalSearchType) > 1 {
		countLabel.SetText(strconv.Itoa(subCount))
	} else {
		ct := database.CheckCount(config.Db_name, GlobalStatus)
		if ct > 0 {
			countLabel.SetText(strconv.Itoa(ct))
		}

	}

	toolbar.AddWidget(countLabel)
	return toolbar
}
