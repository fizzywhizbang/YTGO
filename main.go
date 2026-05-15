package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/fizzywhizbang/YTGO/database"
	"github.com/fizzywhizbang/YTGO/functions"
	"github.com/fizzywhizbang/YTGO/monitor"
	qt "github.com/mappu/miqt/qt6"
)

var ConfigDir = ""
var ConfigFile = "ytgo.json"
var config YTGO

// global variables to pass through different forms
var GlobalStatus = ""
var orderby = "displayname"
var SelectedChannel = ""
var FormSelected = "main"
var GlobalChannelID = ""
var globalSearchTags = ""
var indexSelected = 0
var sectionClicked = 0
var GlobalSearchType = ""
var subCount = 0

var App *qt.QApplication
var Window *qt.QMainWindow

func main() {
	//initial startup to check the config and if it doesn't exist create the config file and database for our program

	if CkConfig() {
		//config exists so next check for the database
		config = ConfigParser()
		database.DBCheck(config.Db_name)

	}
	if config.Monitor {
		cfile := ConfigDir + "/" + ConfigFile
		//begin go routine
		go monitor.MonitorStart(cfile)
	}
	// Create application
	App = qt.NewQApplication(os.Args)
	Window = qt.NewQMainWindow(nil)
	Window.SetWindowTitle("YTGO (Youtube Channel Monitor)")
	Window.SetMinimumSize2(900, 600)
	statuses := database.GetAllStatus(config.Db_name)
	menu := Window.MenuBar()

	selectMenu := menu.AddMenuWithTitle("&Select View")
	subOpts := menu.AddMenuWithTitle("&Channel Opts")
	moveMenu := menu.AddMenuWithTitle("&Move Sub")
	systemSettings := menu.AddMenuWithTitle("&System")

	subStatus := systemSettings.AddActionWithText("Categories")
	subStatus.OnTriggered(func() {
		showStatus()
	})

	tags := systemSettings.AddActionWithText("Edit Tags")
	tags.OnTriggered(func() {
		showTags()
	})

	searches := systemSettings.AddActionWithText("Favorite Searches")
	searches.SetShortcut(qt.NewQKeySequence2("Ctrl+e"))
	searches.OnTriggered(func() {
		showSearches()
	})

	ss := systemSettings.AddActionWithText("System Settings")
	ss.OnTriggered(func() {
		GlobalStatus = ""
		loadSettings()
	})

	binClean := systemSettings.AddActionWithText("Clean Folderwatch")
	binClean.OnTriggered(func() {
		if functions.Cleanfwatch(config.FolderWatch) {
			qt.QMessageBox_Information(nil, "OK", "FolderWatch DIR Cleaned ")
		} else {
			qt.QMessageBox_Information(nil, "OK", "Error cleaning FolderWatch DIR ")
		}
	})

	addChan := subOpts.AddActionWithText("Add Channel")
	addChan.SetShortcutsWithShortcuts(qt.QKeySequence__New)
	addChan.OnTriggered(func() {
		addChannel("")
	})

	updateChanName := subOpts.AddActionWithText("Update Channel Name")
	updateChanName.SetShortcut(qt.NewQKeySequence2("Meta+U"))
	// updateChanName.OnTriggered(func() {
	// })

	vs := subOpts.AddActionWithText("View Settings")
	vs.SetShortcutsWithShortcuts(qt.QKeySequence__Open)
	vs.OnTriggered(func() {
		if GlobalChannelID != "" {
			ChannelSettings(GlobalChannelID)
		}
	})

	sf := subOpts.AddActionWithText("Show Feed")
	sf.SetShortcutsWithShortcuts(qt.QKeySequence__Find)
	sf.OnTriggered(func() {
		if GlobalChannelID != "" {
			feedWindow(GlobalChannelID)
		}
	})
	gu := subOpts.AddActionWithText("GoTo URL")
	gu.SetShortcutsWithShortcuts(qt.QKeySequence__Bold)
	gu.OnTriggered(func() {
		if GlobalChannelID != "" {
			fmt.Println(config.Defbrowser)
			functions.Openbrowser(GlobalChannelID, config.Defbrowser)
		}
	})

	dlu := subOpts.AddActionWithText("Download New Vids")
	dlu.SetShortcutsWithShortcuts(qt.QKeySequence__Save)
	dlu.OnTriggered(func() {
		//check if sub window is open
		count := functions.UpdateChan(config.Db_name, config.FolderWatch, GlobalChannelID, true, true)
		chaninfo := database.GetChanInfo(config.Db_name, GlobalChannelID)
		Window.StatusBar().ShowMessage("Subscription Selected: " + chaninfo.Displayname + " Added: " + strconv.Itoa(count))
	})

	ud := subOpts.AddActionWithText("Update Database")
	ud.SetShortcutsWithShortcuts(qt.QKeySequence__Underline)
	ud.OnTriggered(func() {
		chaninfo := database.GetChanInfo(config.Db_name, GlobalChannelID)
		Window.StatusBar().ShowMessage("Updating: " + chaninfo.Displayname + " " + GlobalChannelID)
		functions.UpdateChan(config.Db_name, config.FolderWatch, GlobalChannelID, false, true)
	})
	delChan := subOpts.AddActionWithText("Delete Channel")
	delChan.SetShortcut(qt.NewQKeySequence2("Meta+D"))

	main := selectMenu.AddActionWithText("Main")
	main.SetShortcut(qt.NewQKeySequence2("Ctrl+M"))
	main.OnTriggered(func() {
		GlobalStatus = ""
		createHomeWindow()
	})
	refresh := selectMenu.AddActionWithText("Refresh View")
	refresh.SetShortcutsWithShortcuts(qt.QKeySequence__Refresh)
	refresh.OnTriggered(func() {
		refreshFunc(Window, App)
	})
	// status menu
	for statuses.Next() {
		var status database.Category
		err := statuses.Scan(&status.ID, &status.Name)
		functions.CheckErr(err, "Unable to retrieve statuses (main.go)")

		a := selectMenu.AddActionWithText(status.Name)
		modifier := "CTRL+" + strconv.Itoa(status.ID)
		a.SetShortcut(qt.NewQKeySequence2(modifier))

		b := moveMenu.AddActionWithText("Move to " + status.Name)
		modifier2 := "META+" + strconv.Itoa(status.ID)
		b.SetShortcut(qt.NewQKeySequence2(modifier2))
	}

	createHomeWindow()

	qt.QApplication_Exec()
}

func refreshFunc(window *qt.QMainWindow, app *qt.QApplication) {
	if GlobalStatus == "" && globalSearchTags == "" {

		createHomeWindow()
	} else {
		if GlobalStatus == "" && globalSearchTags != "" {
			showSubsSearch(globalSearchTags, GlobalSearchType, GlobalStatus)
		} else {
			showSubs(GlobalStatus)
		}

	}
}
func createHomeWindow() {
	verticalLayout := qt.NewQVBoxLayout(nil)

	mainWidget := qt.NewQWidget(nil)

	toolbar := toolbarInit(qt.NewQToolBar2("toolbar"))

	//set menubar
	verticalLayout.SetMenuBar(toolbar.QWidget)

	//add latest to vertical layout

	info := monitorWindow()
	verticalLayout.AddWidget(info.QWidget)

	mainWidget.SetLayout(verticalLayout.QLayout)
	statusBar := Window.StatusBar()
	statusBar.SetObjectName(*qt.NewQAnyStringView3("Status Bar"))

	// // Set main widget as the central widget of the window
	Window.SetCentralWidget(mainWidget)

	// // Show the window
	Window.Show()

}
