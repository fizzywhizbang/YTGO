package main

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

func loadSettings() {
	window := widgets.NewQMainWindow(nil, 0)
	window.SetWindowTitle("Edit Settings")
	window.SetMinimumSize2(800, 300)

	window.ConnectKeyPressEvent(func(e *gui.QKeyEvent) {
		if int32(e.Key()) == int32(core.Qt__Key_Escape) {
			//close window
			window.Close()
		}
	})
	// Create main widget and set the layout
	mainWidget := widgets.NewQWidget(nil, 0)
	mainWidget.SetContentsMargins(0, 2, 0, 0)
	config := ConfigParser()
	//create form layout
	layout := widgets.NewQFormLayout(nil)
	layout.SetFieldGrowthPolicy(widgets.QFormLayout__ExpandingFieldsGrow)

	dbname := widgets.NewQLineEdit(nil)
	dbname.SetText(config.Db_name)
	layout.InsertRow3(2, "Database Name: ", dbname)

	baseDL := widgets.NewQLineEdit(nil)
	baseDL.SetText(config.BaseDL)
	layout.InsertRow3(5, "Base Download Directory: ", baseDL)

	defbrowser := widgets.NewQLineEdit(nil)
	defbrowser.SetText(config.Defbrowser)
	layout.InsertRow3(6, "Default Browser: ", defbrowser)

	folderwatch := widgets.NewQLineEdit(nil)
	folderwatch.SetText(config.FolderWatch)
	layout.InsertRow3(8, "FolderWatch Loc: ", folderwatch)

	monitor := widgets.NewQComboBox(nil)
	items := []string{"True", "False"}
	monitor.AddItems(items)
	monitor.SetCurrentText(config.Monitor)

	buttonGroup := widgets.NewQHBoxLayout()
	save := widgets.NewQPushButton2("Save", nil)
	buttonGroup.AddWidget(save, 0, 0)
	cancel := widgets.NewQPushButton2("Cancel", nil)
	buttonGroup.AddWidget(cancel, 0, 0)
	layout.AddItem(buttonGroup)

	cancel.ConnectClicked(func(checked bool) {
		window.Close()
	})

	save.ConnectClicked(func(checked bool) {
		//if write config returns true then it saved the json file and the user will be notified to
		//restart the program for the new settings to take effect
		//meow
		if writeConfig(dbname.Text(), baseDL.Text(), defbrowser.Text(), folderwatch.Text(), monitor.CurrentText()) {
			widgets.QMessageBox_Information(nil, "OK", "Please restart the program for new settings to take effect", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		}
	})

	mainWidget.SetLayout(layout)
	// mainWidget.Layout().QLayoutItem.SetAlignment(core.Qt__AlignLeft)
	window.SetCentralWidget(mainWidget)

	// Show the window
	window.Show()
}
