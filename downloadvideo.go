package main

import (
	"fmt"

	"github.com/fizzywhizbang/YTGO/functions"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

func downloadVideoForm(channame string, dldir string, chanid string) *widgets.QWidget {
	//if channel id passed add to database otherwise just download

	//create widget to be returned
	layoutWidget := widgets.NewQWidget(nil, 0)

	layout := widgets.NewQFormLayout(nil)
	layout.SetFieldGrowthPolicy(widgets.QFormLayout__ExpandingFieldsGrow)

	//videourl https://www.youtube.com/watch?v=
	yturl := widgets.NewQLineEdit(nil)
	yturl.SetPlaceholderText("Only enter video id")
	yturl.SetToolTip("Press enter after you insert the video ID to fetch the details")
	layout.AddRow3("Vidoe ID: ", yturl)

	//channel name
	chanName := widgets.NewQLineEdit(nil)
	chanName.SetText(channame)
	layout.AddRow3("Channel Name: ", chanName)
	//channel directory
	chanDIR := widgets.NewQLineEdit(nil)
	chanDIR.SetText(dldir)
	layout.AddRow3("Directory: ", chanDIR)
	//download and cancel buttons

	//video title
	videoTitle := widgets.NewQLineEdit(nil)
	layout.AddRow3("Video Title: ", videoTitle)
	//video description
	videoDesc := widgets.NewQTextEdit(nil)
	videoDesc.SetReadOnly(true)
	layout.AddRow3("Description: ", videoDesc)
	//video date
	videoDate := widgets.NewQLineEdit(nil)
	layout.AddRow3("Date Published: ", videoDate)

	startButton := widgets.NewQPushButton(nil)
	startButton.SetText("Start Download")

	layout.AddRow3(" ", startButton)

	msgbox := widgets.NewQLabel(nil, 0)
	layout.AddRow3(" ", msgbox)
	videoTitleText := ""
	yturl.ConnectKeyReleaseEvent(func(event *gui.QKeyEvent) {
		// 00H8gY69PKo
		if int32(event.Key()) == int32(core.Qt__Key_Return) || int32(event.Key()) == int32(core.Qt__Key_Enter) {
			// if FormSelected == editchannel then only return video info
			if FormSelected == "EditChannel" {
				video := functions.GetVideoInfo(yturl.Text())
				videoTitle.SetText(video.Title)
				if len(video.Description) < 1 {
					//use the title if no description
					videoTitleText = video.Title
					videoDesc.SetText(videoTitleText)
				} else {
					videoDesc.SetText(video.Description)
				}
				videoDate.SetText(functions.DateConvertTrim(video.Publish_date, 10))
				msgbox.SetText(" ") //remove from message box
			}

			fmt.Println(FormSelected)
		}
	})
	startButton.ConnectClicked(func(checked bool) {
		if len(videoTitle.Text()) >= 1 && len(videoDate.Text()) >= 1 {
			functions.MkCrawljob(config.Db_name, config.FolderWatch, chanid, videoTitle.Text(), yturl.Text(), videoDate.Text(), 1)
		}

		msgbox.SetText(videoTitleText + " added to queue")
	})

	instructionsLabel := widgets.NewQLabel2("After inserting the Video ID press enter and I'll fetch the video details", nil, 0)
	instructionsLabel.SetFont(gui.NewQFont2("Times", 12, 1, true))
	layout.AddRow5(instructionsLabel)
	layoutWidget.SetLayout(layout)

	return layoutWidget
}
