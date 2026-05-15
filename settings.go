package main

import (
	"strconv"

	qt "github.com/mappu/miqt/qt6"
)

func loadSettings() {
	window := qt.NewQMainWindow(nil)
	window.SetWindowTitle("Edit Settings")
	window.SetMinimumSize2(800, 300)

	window.OnKeyPressEvent(func(super func(*qt.QKeyEvent), e *qt.QKeyEvent) {
		if int32(e.Key()) == int32(qt.Key_Escape) {
			//close window
			window.Close()
		}
	})
	// Create main widget and set the layout
	mainWidget := qt.NewQWidget(nil)
	mainWidget.SetContentsMargins(0, 2, 0, 0)
	config := ConfigParser()
	//create form layout
	layout := qt.NewQFormLayout(nil)
	layout.SetFieldGrowthPolicy(qt.QFormLayout__ExpandingFieldsGrow)

	dbname := qt.NewQLineEdit(nil)
	dbname.SetText(config.Db_name)
	layout.InsertRow3(2, "Database Name: ", dbname.QWidget)

	baseDL := qt.NewQLineEdit(nil)
	baseDL.SetText(config.BaseDL)
	layout.InsertRow3(5, "Base Download Directory: ", baseDL.QWidget)

	defbrowser := qt.NewQLineEdit(nil)
	defbrowser.SetText(config.Defbrowser)
	layout.InsertRow3(6, "Default Browser: ", defbrowser.QWidget)

	folderwatch := qt.NewQLineEdit(nil)
	folderwatch.SetText(config.FolderWatch)
	layout.InsertRow3(7, "FolderWatch Loc: ", folderwatch.QWidget)

	monitor := qt.NewQComboBox(nil)
	items := []string{"true", "false"}
	monitor.AddItems(items)

	monitor.SetCurrentText(strconv.FormatBool(config.Monitor))
	layout.InsertRow3(8, "Enable Monitor:", monitor.QWidget)

	buttonGroup := qt.NewQHBoxLayout(nil)
	save := qt.NewQPushButton3("Save")
	buttonGroup.AddWidget(save.QWidget)
	cancel := qt.NewQPushButton3("Cancel")
	buttonGroup.AddWidget(cancel.QWidget)
	layout.AddItem(buttonGroup.QLayoutItem)

	cancel.OnClicked(func() {
		window.Close()
	})

	save.OnClicked(func() {
		//if write config returns true then it saved the json file and the user will be notified to
		//restart the program for the new settings to take effect
		//meow
		if writeConfig(dbname.Text(), baseDL.Text(), defbrowser.Text(), folderwatch.Text(), monitor.CurrentText()) {
			qt.QMessageBox_Information(nil, "OK", "Please restart the program for new settings to take effect")
		}
	})

	mainWidget.SetLayout(layout.QLayout)
	// mainWidget.Layout().QLayoutItem.SetAlignment(qt.Qt__AlignLeft)
	window.SetCentralWidget(mainWidget)

	// Show the window
	window.Show()
}
