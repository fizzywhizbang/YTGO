package main

import (
	"strconv"

	"github.com/fizzywhizbang/YTGO/database"
	"github.com/fizzywhizbang/YTGO/functions"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

func showTags() {
	window := widgets.NewQMainWindow(nil, 0)
	window.SetWindowTitle("Edit Settings")
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

	tags := database.GetAllTags(config.Db_name, "id")
	tableWidget.SetColumnCount(2)
	tableWidget.SetRowCount(database.TagCount(config.Db_name) + 1)
	tableWidget.SetHorizontalHeaderLabels([]string{"ID", "Name"})
	tableWidget.SetAlternatingRowColors(true)
	tableColors := "alternate-background-color: #88DD88; background-color:#FFFFFF; color:#000000; font-size: 12px;"
	tableWidget.SetStyleSheet(tableColors)

	rowCounter := 0
	count := database.TagCount(config.Db_name)
	if count >= 1 {
		for tags.Next() {
			var tag database.Tags
			err := tags.Scan(&tag.ID, &tag.Name)
			functions.CheckErr(err, "error getting tags")
			id := widgets.NewQTableWidgetItem2(strconv.Itoa(tag.ID), 0)
			name := widgets.NewQTableWidgetItem2(tag.Name, 0)
			tableWidget.SetItem(rowCounter, 0, id)
			tableWidget.SetItem(rowCounter, 1, name)
			name.SetData(1, core.NewQVariant12(tag.Name))

			id.SetData(0, core.NewQVariant12(strconv.Itoa(tag.ID)))
			id.SetFlags(core.Qt__NoItemFlags)
			rowCounter++
		}
	}

	widgets.NewQTableWidgetItem2("", 0)
	widgets.NewQTableWidgetItem2("", 0)

	tableWidget.SetColumnWidth(1, 300)
	tableWidget.ConnectCellChanged(func(row, column int) {

		// index := tableWidget.IndexFromItem(tableWidget.CurrentItem())
		id := tableWidget.Item(row, 0).Text()
		data := tableWidget.Item(row, column).Text()
		if id == "" {
			//new record
			database.TagInsert(config.Db_name, data)
			window.Close()
			showTags()
		} else {
			//update status in database
			database.TagUpdate(config.Db_name, id, data)
			//reload view
			window.Close()
			showTags()
		}
	})
	tableWidget.SetSortingEnabled(true)

	window.SetCentralWidget(tableWidget)

	// Show the window
	window.Show()

}
