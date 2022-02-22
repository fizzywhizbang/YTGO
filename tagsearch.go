package main

import (
	"strings"

	"github.com/fizzywhizbang/YTGO/database"
	"github.com/fizzywhizbang/YTGO/functions"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

func showTagSearch() {
	window := widgets.NewQMainWindow(nil, 0)
	window.SetWindowTitle("Tag Search")
	// window.SetMinimumWidth(200)
	window.SetMaximumSize2(250, 200)
	mainWidget := widgets.NewQWidget(nil, 0)
	mainWidget.SetContentsMargins(0, 2, 0, 0)
	searchTags := []string{}
	window.ConnectKeyPressEvent(func(e *gui.QKeyEvent) {
		if int32(e.Key()) == int32(core.Qt__Key_Escape) {
			//close window
			window.Close()
		}
	})

	scrollArea := widgets.NewQScrollArea(window)
	scrollArea.SetHorizontalScrollBarPolicy(core.Qt__ScrollBarAlwaysOff)
	scrollArea.SetVerticalScrollBarPolicy(core.Qt__ScrollBarAlwaysOn)
	scrollArea.SetWidgetResizable(true)
	scrollArea.SetWidget(mainWidget)
	//create form layout
	form := widgets.NewQFormLayout(nil)

	tags := database.GetAllTags(config.Db_name, "tag")

	rowCounter := 0
	count := database.TagCount(config.Db_name)
	if count >= 1 {
		for tags.Next() {
			var tag database.Tags
			err := tags.Scan(&tag.ID, &tag.Name)
			functions.CheckErr(err, "unable to retrieve tags")
			//create form items

			checkbox := widgets.NewQCheckBox2(tag.Name, nil)
			if GlobalSearchType == "Tags" && contains(tag.Name) {
				checkbox.SetChecked(true)
				//add to array because it was there before
				searchTags = append(searchTags, "#"+tag.Name)
			}

			form.InsertRow5(rowCounter, checkbox)
			checkbox.ConnectClicked(func(checked bool) {
				if checked {
					//true
					searchTags = append(searchTags, "#"+checkbox.Text())

				} else {
					//false
					searchTags = remove(searchTags, "#"+checkbox.Text())
				}
				//turn this into text for global search

				GlobalSearchType = "Tags"
				globalSearchTags = strings.Join(searchTags, " ")
				if len(globalSearchTags) > 1 {
					showSubsSearch(globalSearchTags, GlobalSearchType, GlobalStatus)
				}

			})
			rowCounter++
		}

	}

	mainWidget.SetLayout(form)

	scroll_layout := widgets.NewQVBoxLayout2(nil)
	scroll_layout.AddWidget(scrollArea, 0, 0)
	scroll_layout.SetContentsMargins(0, 0, 0, 0)
	containerWidget := widgets.NewQWidget(nil, 0)
	containerWidget.SetLayout(scroll_layout)

	window.SetCentralWidget(containerWidget)

	// Show the window
	window.Show()

}

func remove(slice []string, s string) []string {
	for i, v := range slice {
		if v == s {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func contains(n string) bool {
	s := strings.ReplaceAll(globalSearchTags, "#", "")
	slice := strings.Split(s, " ")
	for _, a := range slice {
		if a == n {
			return true
		}
	}
	return false
}
