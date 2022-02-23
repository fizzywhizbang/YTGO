package main

import (
	"os"
	"strconv"

	"github.com/fizzywhizbang/YTGO/database"
	"github.com/fizzywhizbang/YTGO/functions"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

var ConfigDir = ""
var ConfigFile = "ytgo.json"
var config YTGO

//global variables to pass through different forms
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

var App *widgets.QApplication
var Window *widgets.QMainWindow

func main() {
	//initial startup to check the config and if it doesn't exist create the config file and database for our program

	if CkConfig() {
		//config exists so next check for the database
		config = ConfigParser()
		database.DBCheck(config.Db_name)

	}

	// Create application
	App = widgets.NewQApplication(len(os.Args), os.Args)
	Window = widgets.NewQMainWindow(nil, 0)
	Window.SetWindowTitle("YTGO (Youtube Channel Monitor)")
	Window.SetMinimumSize2(900, 600)
	statuses := database.GetAllStatus(config.Db_name)
	menu := Window.MenuBar()

	selectMenu := menu.AddMenu2("&Select View")
	subOpts := menu.AddMenu2("&Channel Opts")
	moveMenu := menu.AddMenu2("&Move Sub")
	systemSettings := menu.AddMenu2("&System")

	subStatus := systemSettings.AddAction("Categories")
	subStatus.ConnectTriggered(func(checked bool) {
		showStatus()
	})

	tags := systemSettings.AddAction("Edit Tags")
	tags.ConnectTriggered(func(checked bool) {
		showTags()
	})

	searches := systemSettings.AddAction("Favorite Searches")
	searches.SetShortcut(gui.NewQKeySequence2("Ctrl+e", gui.QKeySequence__NativeText))
	searches.ConnectTriggered(func(checked bool) {
		showSearches()
	})

	ss := systemSettings.AddAction("System Settings")
	ss.ConnectTriggered(func(checked bool) {
		GlobalStatus = ""
		loadSettings()
	})

	addChan := subOpts.AddAction("Add Channel")
	addChan.SetShortcuts2(gui.QKeySequence__New)
	addChan.ConnectTriggered(func(checked bool) {
		addChannel("")
	})

	updateChanName := subOpts.AddAction("Update Channel Name")
	updateChanName.SetShortcut(gui.NewQKeySequence2("Meta+U", gui.QKeySequence__NativeText))
	// updateChanName.ConnectTriggered(func(checked bool) {
	// })

	vs := subOpts.AddAction("View Settings")
	vs.SetShortcuts2(gui.QKeySequence__Open)
	vs.ConnectTriggered(func(checked bool) {
		if GlobalChannelID != "" {
			ChannelSettings(GlobalChannelID)
		}
	})

	sf := subOpts.AddAction("Show Feed")
	sf.SetShortcuts2(gui.QKeySequence__Find)
	sf.ConnectTriggered(func(checked bool) {
		if GlobalChannelID != "" {
			feedWindow(GlobalChannelID)
		}
	})
	gu := subOpts.AddAction("GoTo URL")
	gu.SetShortcuts2(gui.QKeySequence__Bold)
	gu.ConnectTriggered(func(checked bool) {
		if GlobalChannelID != "" {

			functions.Openbrowser(config.Defbrowser, GlobalChannelID)
		}
	})

	dlu := subOpts.AddAction("Download New Vids")
	dlu.SetShortcuts2(gui.QKeySequence__Save)
	dlu.ConnectTriggered(func(checked bool) {
		//check if sub window is open
		count := functions.UpdateChan(config.Db_name, config.FolderWatch, GlobalChannelID, true, true)
		chaninfo := database.GetChanInfo(config.Db_name, GlobalChannelID)
		Window.StatusBar().ShowMessage("Subscription Selected: "+chaninfo.Displayname+" Added: "+strconv.Itoa(count), 0)
	})

	ud := subOpts.AddAction("Update Database")
	ud.SetShortcuts2(gui.QKeySequence__Underline)
	ud.ConnectTriggered(func(checked bool) {
		chaninfo := database.GetChanInfo(config.Db_name, GlobalChannelID)
		Window.StatusBar().ShowMessage("Updating: "+chaninfo.Displayname+" "+GlobalChannelID, 0)
		functions.UpdateChan(config.Db_name, config.FolderWatch, GlobalChannelID, false, true)
	})
	delChan := subOpts.AddAction("Delete Channel")
	delChan.SetShortcut(gui.NewQKeySequence2("Meta+D", gui.QKeySequence__NativeText))

	main := selectMenu.AddAction("Main")
	main.SetShortcut(gui.NewQKeySequence2("Ctrl+M", gui.QKeySequence__NativeText))
	main.ConnectTriggered(func(checked bool) {
		GlobalStatus = ""
		createHomeWindow()
	})
	refresh := selectMenu.AddAction("Refresh View")
	refresh.SetShortcuts2(gui.QKeySequence__Refresh)
	refresh.ConnectTriggered(func(checked bool) {
		refreshFunc(Window, App)
	})
	// status menu
	for statuses.Next() {
		var status database.Category
		err := statuses.Scan(&status.ID, &status.Name)
		functions.CheckErr(err, "Unable to retrieve statuses (main.go)")

		a := selectMenu.AddAction(status.Name)
		modifier := "CTRL+" + strconv.Itoa(status.ID)
		a.SetShortcut(gui.NewQKeySequence2(modifier, gui.QKeySequence__NativeText))

		b := moveMenu.AddAction("Move to " + status.Name)
		modifier2 := "META+" + strconv.Itoa(status.ID)
		b.SetShortcut(gui.NewQKeySequence2(modifier2, gui.QKeySequence__NativeText))
	}

	createHomeWindow()

	App.Exec()
}

func refreshFunc(window *widgets.QMainWindow, app *widgets.QApplication) {
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
	verticalLayout := widgets.NewQVBoxLayout()

	mainWidget := widgets.NewQWidget(nil, 0)

	toolbar := toolbarInit(widgets.NewQToolBar2(nil))

	//set menubar
	verticalLayout.SetMenuBar(toolbar)

	//add latest to vertical layout

	info := monitorWindow()
	verticalLayout.AddWidget(info, 0, 0)

	mainWidget.SetLayout(verticalLayout)
	statusBar := Window.StatusBar()
	statusBar.SetObjectName("Status Bar")

	// // Set main widget as the central widget of the window
	Window.SetCentralWidget(mainWidget)

	// // Show the window
	Window.Show()

}
