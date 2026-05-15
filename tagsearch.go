package main

import (
	"strings"

	"github.com/fizzywhizbang/YTGO/database"
	"github.com/fizzywhizbang/YTGO/functions"
	qt "github.com/mappu/miqt/qt6"
)

func showTagSearch() {
	window := qt.NewQMainWindow(nil)
	window.SetWindowTitle("Tag Search")
	// window.SetMinimumWidth(200)
	window.SetMaximumSize2(250, 200)
	mainWidget := qt.NewQWidget(nil)
	mainWidget.SetContentsMargins(0, 2, 0, 0)
	searchTags := []string{}
	window.OnKeyPressEvent(func(super func(e *qt.QKeyEvent), e *qt.QKeyEvent) {
		if int32(e.Key()) == int32(qt.Key_Escape) {
			//close window
			window.Close()
		}
	})

	scrollArea := qt.NewQScrollArea(window.QWidget)
	scrollArea.SetHorizontalScrollBarPolicy(qt.ScrollBarAlwaysOff)
	scrollArea.SetVerticalScrollBarPolicy(qt.ScrollBarAlwaysOn)
	scrollArea.SetWidgetResizable(true)
	scrollArea.SetWidget(mainWidget)
	//create form layout
	form := qt.NewQFormLayout(nil)

	tags := database.GetAllTags(config.Db_name, "tag")

	rowCounter := 0
	count := database.TagCount(config.Db_name)
	if count >= 1 {
		for tags.Next() {
			var tag database.Tags
			err := tags.Scan(&tag.ID, &tag.Name)
			functions.CheckErr(err, "unable to retrieve tags")
			//create form items

			checkbox := qt.NewQCheckBox3(tag.Name)
			if GlobalSearchType == "Tags" && contains(tag.Name) {
				checkbox.SetChecked(true)
				//add to array because it was there before
				searchTags = append(searchTags, "#"+tag.Name)
			}

			form.InsertRow5(rowCounter, checkbox.QWidget)
			checkbox.OnClicked(func() {
				checked := checkbox.IsChecked()
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

	mainWidget.SetLayout(form.QLayout)

	scroll_layout := qt.NewQVBoxLayout2()
	scroll_layout.AddWidget(scrollArea.QWidget)
	scroll_layout.SetContentsMargins(0, 0, 0, 0)
	containerWidget := qt.NewQWidget(nil)
	containerWidget.SetLayout(scroll_layout.QLayout)

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
