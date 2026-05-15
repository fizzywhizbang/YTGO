package main

import (
	"strconv"

	"github.com/fizzywhizbang/YTGO/database"
	"github.com/fizzywhizbang/YTGO/functions"
	qt "github.com/mappu/miqt/qt6"
)

func showTags() {
	window := qt.NewQMainWindow(nil)
	window.SetWindowTitle("Edit Settings")
	window.SetMinimumSize2(800, 300)
	mainWidget := qt.NewQWidget(nil)
	mainWidget.SetContentsMargins(0, 2, 0, 0)
	window.OnKeyPressEvent(func(super func(e *qt.QKeyEvent), e *qt.QKeyEvent) {
		if int32(e.Key()) == int32(qt.Key_Escape) {
			//close window
			window.Close()
		}
	})
	tableWidget := qt.NewQTableWidget(mainWidget)

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
			id := qt.NewQTableWidgetItem2(strconv.Itoa(tag.ID))
			name := qt.NewQTableWidgetItem2(tag.Name)
			tableWidget.SetItem(rowCounter, 0, id)
			tableWidget.SetItem(rowCounter, 1, name)
			name.SetData(1, qt.NewQVariant11(tag.Name))

			id.SetData(0, qt.NewQVariant11(strconv.Itoa(tag.ID)))
			id.SetFlags(qt.NoItemFlags)
			rowCounter++
		}
	}

	qt.NewQTableWidgetItem2("")
	qt.NewQTableWidgetItem2("")

	tableWidget.SetColumnWidth(1, 300)
	tableWidget.OnCellChanged(func(row, column int) {

		// index := tableWidget.IndexFromItem(tableWidget.CurrentItem())

		data := tableWidget.Item(row, column).Text()
		if row == count {
			//new record
			database.TagInsert(config.Db_name, data)
			window.Close()
			showTags()
		} else {
			id := tableWidget.Item(row, 0).Text()
			//update status in database
			database.TagUpdate(config.Db_name, id, data)
			//reload view
			window.Close()
			showTags()
		}
	})
	tableWidget.SetSortingEnabled(true)

	window.SetCentralWidget(tableWidget.QWidget)

	// Show the window
	window.Show()

}
