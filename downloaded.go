package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/aquilax/truncate"
	"github.com/fizzywhizbang/YTGO/database"
	"github.com/fizzywhizbang/YTGO/functions"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

//this is a generic treeview form for displaying feeds and downloaded videos
func contentListDL(chanid string) *widgets.QTreeWidget {

	// results := database.GetChannelVids(config.Db_name, chanid)

	treeWidget := widgets.NewQTreeWidget(nil)
	treeWidget.SetColumnCount(4)
	treeWidget.SetObjectName("treewidget")
	treeWidget.Header().SetStretchLastSection(false)
	treeWidget.Header().SetSectionsClickable(true)
	treeWidget.SetSortingEnabled(true)
	treeWidget.SortByColumn(2, core.Qt__SortOrder(0))
	treeWidget.SetAlternatingRowColors(true)
	treeWidget.HorizontalScrollBar().SetHidden(true)
	tableColors := "alternate-background-color: #88DD88; background-color:#FFFFFF; color:#000000; font-size: 12px;"
	treeWidget.SetStyleSheet(tableColors)
	treeWidget.Header()
	treeWidget.SetHeaderLabels([]string{"VidoeID", "Title", "Date", "Status"})

	contentFill(chanid, treeWidget)

	treeWidget.ResizeColumnToContents(0)
	treeWidget.ResizeColumnToContents(1)
	treeWidget.ResizeColumnToContents(2)
	treeWidget.ConnectDoubleClicked(func(index *core.QModelIndex) {
		data := index.Data(int(core.Qt__UserRole)).ToString()
		item := treeWidget.CurrentColumn()
		fmt.Println(item)
		if item == 0 {
			//re-download the video
			functions.MkCrawljob(config.Db_name, config.FolderWatch, chanid, treeWidget.CurrentItem().Text(1), data, treeWidget.CurrentItem().Text(2), 1)
			widgets.QMessageBox_Information(nil, "OK", "Added to Queue "+treeWidget.CurrentItem().Text(1), widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		}
		if item == 1 {
			// search youtube for the item
			// text := cleanUnicode(treeWidget.CurrentItem().Text(1))
			//seting sort order by date uploaded
			/*
				sp=CAI%253D #search by upload date
				sp=CAASAhAB #search by  relevance
				sp=EgQIAhAB #today and relevance
				sp=CAISBAgCEAE%253D #today & upload date
				sp=CAMSBAgCEAE%253D #today and view count
				sp=CAESBAgCEAE%253D #today and rating
				sp=CAESAhAB #rating only
			*/
			url := YtSearchPrefix + treeWidget.CurrentItem().Text(1) + "&sp=CAI%253D" //order by upload date
			functions.Openbrowser(config.Defbrowser, url)
		}

	})
	return treeWidget
}

func contentFill(chanid string, treeWidget *widgets.QTreeWidget) *widgets.QTreeWidget {
	results := database.GetChannelVids(config.Db_name, chanid)

	for results.Next() {
		var video database.Video
		err := results.Scan(&video.ID, &video.YT_videoid, &video.Title, &video.Description, &video.Publisher, &video.Publish_date, &video.Downloaded)
		functions.CheckErr(err, "Unable to get videos for channel")
		watched := "Queued"
		if video.Downloaded == 2 {
			watched = "Skipped"
		}
		if video.Downloaded == 1 {
			watched = "Downloaded"
		}
		truncated := truncate.Truncate(video.Title, 65, "...", truncate.PositionEnd)
		treewidgetItem := widgets.NewQTreeWidgetItem2([]string{video.YT_videoid, truncated, functions.DateConvertTrim(video.Publish_date, 10), watched}, 0)
		treewidgetItem.SetData(0, int(core.Qt__UserRole), core.NewQVariant12(video.YT_videoid))

		treeWidget.AddTopLevelItem(treewidgetItem)
	}

	return treeWidget
}

func cleanUnicode(str string) string {

	re := regexp.MustCompile("[[:^ascii:]]")
	text := re.ReplaceAllLiteralString(str, "")
	text = strings.Trim(text, " ")
	text = strings.Replace(text, " ", "+", -1)
	return text
}
