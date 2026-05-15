package main

import (
	"strconv"

	"github.com/fizzywhizbang/YTGO/database"
	"github.com/fizzywhizbang/YTGO/functions"
	qt "github.com/mappu/miqt/qt6"
)

func showStatus() {
	window := qt.NewQMainWindow(nil)
	window.SetWindowTitle("Edit Settings")
	window.SetMinimumSize2(800, 300)
	mainWidget := qt.NewQWidget(nil)
	mainWidget.SetContentsMargins(0, 2, 0, 0)
	window.OnKeyPressEvent(func(super func(*qt.QKeyEvent), e *qt.QKeyEvent) {
		if int32(e.Key()) == int32(qt.Key_Escape) {
			//close window
			window.Close()
		}
	})
	tableWidget := qt.NewQTableWidget(mainWidget)

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
		id := qt.NewQTableWidgetItem2(strconv.Itoa(status.ID))
		name := qt.NewQTableWidgetItem2(status.Name)
		tableWidget.SetItem(rowCounter, 0, id)
		tableWidget.SetItem(rowCounter, 1, name)
		name.SetData(1, qt.NewQVariant11(status.Name))

		id.SetData(0, qt.NewQVariant11(strconv.Itoa(status.ID)))
		id.SetFlags(qt.NoItemFlags)
		rowCounter++
	}

	qt.NewQTableWidgetItem2("")
	qt.NewQTableWidgetItem2("")

	tableWidget.SetColumnWidth(1, 300)
	tableWidget.OnCellChanged(func(row, column int) {
		id := ""
		// index := tableWidget.IndexFromItem(tableWidget.CurrentItem())

		if row == rowCounter {
			// new record
			id = ""
		} else {
			id = tableWidget.Item(row, 0).Text()
		}

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

		qt.QMessageBox_Information(nil, "App Restart", "Settings saved.\nRestart application for new settings to be visible")
		// fmt.Println(index.Data(int(qt.Qt__UserRole)).ToString())
	})
	tableWidget.SetSortingEnabled(true)

	window.SetCentralWidget(tableWidget.QWidget)

	// Show the window
	window.Show()

}
