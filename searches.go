package main

import (
	"strconv"

	"github.com/fizzywhizbang/YTGO/database"
	"github.com/fizzywhizbang/YTGO/functions"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

func showSearches() {
	window := widgets.NewQMainWindow(nil, 0)
	window.SetWindowTitle("Searches")
	window.SetMinimumSize2(800, 300)
	mainWidget := widgets.NewQWidget(nil, 0)
	mainWidget.SetContentsMargins(0, 2, 0, 0)
	window.ConnectKeyPressEvent(func(e *gui.QKeyEvent) {
		if int32(e.Key()) == int32(core.Qt__Key_Escape) {
			//close window
			window.Close()
		}
	})
	tableWidget := widgets.NewQTableWidget(mainWidget)

	searches := database.GetAllSearches(config.Db_name)
	tableWidget.SetColumnCount(4)
	tableWidget.SetRowCount(database.SearchCount(config.Db_name) + 1)
	tableWidget.SetHorizontalHeaderLabels([]string{"ID", "Name", "Link (delete link and it will remove from db)", "GoTo"})
	tableWidget.SetAlternatingRowColors(true)
	tableColors := "alternate-background-color: #88DD88; background-color:#FFFFFF; color:#000000; font-size: 12px;"
	tableWidget.SetStyleSheet(tableColors)

	rowCounter := 0
	for searches.Next() {
		var search database.Search
		err := searches.Scan(&search.ID, &search.Name, &search.Link)
		functions.CheckErr(err, "Could not retrieve searches")
		id := widgets.NewQTableWidgetItem2(strconv.Itoa(search.ID), 0)
		name := widgets.NewQTableWidgetItem2(search.Name, 0)
		link := widgets.NewQTableWidgetItem2(search.Link, 0)
		tableWidget.SetItem(rowCounter, 0, id)
		tableWidget.SetItem(rowCounter, 1, name)
		tableWidget.SetItem(rowCounter, 2, link)
		gobutton := widgets.NewQTableWidgetItem2("GO", 0)
		gobutton.SetFlags(core.Qt__ItemIsEnabled) //disable this cell as it's there to follow the link
		tableWidget.SetItem(rowCounter, 3, gobutton)

		name.SetData(1, core.NewQVariant12(search.Name))

		id.SetData(0, core.NewQVariant12(strconv.Itoa(search.ID)))
		id.SetFlags(core.Qt__NoItemFlags)
		rowCounter++
	}

	widgets.NewQTableWidgetItem2("", 0)
	widgets.NewQTableWidgetItem2("", 0)

	tableWidget.ResizeColumnToContents(0)
	tableWidget.ResizeColumnToContents(1)
	tableWidget.SetColumnWidth(2, 400)
	tableWidget.ResizeColumnToContents(3)

	tableWidget.ConnectCellDoubleClicked(func(row, column int) {

		if column == 3 {
			url := tableWidget.Item(row, 2).Text()
			functions.Openbrowser(config.Defbrowser, url)
		}
	})

	tableWidget.ConnectCellChanged(func(row, column int) {

		// index := tableWidget.IndexFromItem(tableWidget.CurrentItem())
		id := tableWidget.Item(row, 0).Text()
		data := tableWidget.Item(row, 1).Text()
		link := tableWidget.Item(row, 2).Text()
		if link != "" { //if link empty do nothing
			if id == "" {
				//new record
				database.SearchInsert(config.Db_name, data, link)
				window.Close()
				showSearches()
			} else {
				//update status in database
				database.SearchUpdate(config.Db_name, id, data, link)
				//reload view
				window.Close()
				showSearches()
			}
		} else {
			//if link empty but id is valid then delete the item
			if id != "" {
				database.SearchDelete(config.Db_name, id)
				window.Close()
				showSearches()
			}
		}

	})
	tableWidget.SetSortingEnabled(true)

	window.SetCentralWidget(tableWidget)

	// Show the window
	window.Show()

}
