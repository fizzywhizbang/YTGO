package main

import (
	"fmt"

	"github.com/fizzywhizbang/YTGO/functions"
	qt "github.com/mappu/miqt/qt6"
)

func downloadVideoForm(channame string, dldir string, chanid string) *qt.QWidget {
	//if channel id passed add to database otherwise just download

	//create widget to be returned
	layoutWidget := qt.NewQWidget(nil)

	layout := qt.NewQFormLayout(nil)
	layout.SetFieldGrowthPolicy(qt.QFormLayout__ExpandingFieldsGrow)

	//videourl https://www.youtube.com/watch?v=
	yturl := qt.NewQLineEdit(nil)
	yturl.SetPlaceholderText("Only enter video id")
	yturl.SetToolTip("Press enter after you insert the video ID to fetch the details")
	layout.AddRow3("Vidoe ID: ", yturl.QWidget)

	//channel name
	chanName := qt.NewQLineEdit(nil)
	chanName.SetText(channame)
	layout.AddRow3("Channel Name: ", chanName.QWidget)
	//channel directory
	chanDIR := qt.NewQLineEdit(nil)
	chanDIR.SetText(dldir)
	layout.AddRow3("Directory: ", chanDIR.QWidget)
	//download and cancel buttons

	//video title
	videoTitle := qt.NewQLineEdit(nil)
	layout.AddRow3("Video Title: ", videoTitle.QWidget)
	//video description
	videoDesc := qt.NewQTextEdit(nil)
	videoDesc.SetReadOnly(true)
	layout.AddRow3("Description: ", videoDesc.QWidget)
	//video date
	videoDate := qt.NewQLineEdit(nil)
	layout.AddRow3("Date Published: ", videoDate.QWidget)

	startButton := qt.NewQPushButton(nil)
	startButton.SetText("Start Download")

	layout.AddRow3(" ", startButton.QWidget)

	msgbox := qt.NewQLabel(nil)
	layout.AddRow3(" ", msgbox.QWidget)
	videoTitleText := ""
	yturl.OnKeyReleaseEvent(func(super func(event *qt.QKeyEvent), event *qt.QKeyEvent) {
		// 00H8gY69PKo
		if int32(event.Key()) == int32(qt.Key_Return) || int32(event.Key()) == int32(qt.Key_Enter) {
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
	startButton.OnClicked(func() {
		if len(videoTitle.Text()) >= 1 && len(videoDate.Text()) >= 1 {
			functions.MkCrawljob(config.Db_name, config.FolderWatch, chanid, videoTitle.Text(), yturl.Text(), videoDate.Text(), 1)
		}

		msgbox.SetText(videoTitleText + " added to queue")
	})

	instructionsLabel := qt.NewQLabel3("After inserting the Video ID press enter and I'll fetch the video details")
	instructionsLabel.SetFont(qt.NewQFont6("Times", 12))
	layout.AddRowWithWidget(instructionsLabel.QWidget)
	layoutWidget.SetLayout(layout.QLayout)

	return layoutWidget
}
