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

func showStatus() {
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

	statuses := database.GetAllStatus(config.Db_name)
	tableWidget.SetColumnCount(2)
	tableWidget.SetRowCount(database.StatusCount(config.Db_name) + 1)
	tableWidget.SetHorizontalHeaderLabels([]string{"ID", "Name"})
	tableWidget.SetAlternatingRowColors(true)
	tableColors := "alternate-background-color: #88DD88; background-color:#FFFFFF; color:#000000; font-size: 12px;"
	tableWidget.SetStyleSheet(tableColors)

	rowCounter := 0
	for statuses.Next() {
		var status database.Category
		err := statuses.Scan(&status.ID, &status.Name)
		functions.CheckErr(err, "Unable to get statuses")
		id := widgets.NewQTableWidgetItem2(strconv.Itoa(status.ID), 0)
		name := widgets.NewQTableWidgetItem2(status.Name, 0)
		tableWidget.SetItem(rowCounter, 0, id)
		tableWidget.SetItem(rowCounter, 1, name)
		name.SetData(1, core.NewQVariant12(status.Name))

		id.SetData(0, core.NewQVariant12(strconv.Itoa(status.ID)))
		id.SetFlags(core.Qt__NoItemFlags)
		rowCounter++
	}

	widgets.NewQTableWidgetItem2("", 0)
	widgets.NewQTableWidgetItem2("", 0)

	tableWidget.SetColumnWidth(1, 300)
	tableWidget.ConnectCellChanged(func(row, column int) {
		fmt.Println(row)
		// index := tableWidget.IndexFromItem(tableWidget.CurrentItem())
		id := tableWidget.Item(row, 0).Text()
		data := tableWidget.Item(row, column).Text()
		if id == "" {
			//new record
			database.StatusInsert(config.Db_name, data)

			window.Close()
			showStatus()
		} else {
			//update status in database
			database.StatusUpdate(config.Db_name, id, data)
			//reload view
			window.Close()
			showStatus()
		}

		widgets.QMessageBox_Information(nil, "App Restart", "Settings saved.\nRestart application for new settings to be visible", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		// fmt.Println(index.Data(int(core.Qt__UserRole)).ToString())
	})
	tableWidget.SetSortingEnabled(true)

	window.SetCentralWidget(tableWidget)

	// Show the window
	window.Show()

}
