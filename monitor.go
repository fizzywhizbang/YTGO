package main

import (
	"strconv"

	"github.com/fizzywhizbang/YTGO/database"
	"github.com/fizzywhizbang/YTGO/functions"
	"github.com/therecipe/qt/widgets"
)

func monitorWindow() *widgets.QTextBrowser {

	results := database.GetLatestVideos(config.Db_name)
	var video database.Video

	textarea := widgets.NewQTextBrowser(nil)
	textarea.SetOpenExternalLinks(true)
	message := ""
	message += "<div class=\"font-weight: bold; font-size:16px; color:green\">#############################################################<br>"
	message += "last update was at " + functions.DateConvert(database.GetLastCheck(config.Db_name)) + "<br>"
	message += "The following channels checked<br>"
	message += "#############################################################</div>"
	message += "<div class=\"font-weight: bold; font-size:16px; color:green\">#############################################################<br>"
	statuses := database.GetAllStatus(config.Db_name)
	statusCounts := "Category Counts: "
	for statuses.Next() {
		var status database.Category
		err := statuses.Scan(&status.ID, &status.Name)
		functions.CheckErr(err, "Could not get statuses for monitor area")
		statusCounts += status.Name + " (" + strconv.Itoa(database.CheckCount(config.Db_name, strconv.Itoa(status.ID))) + ") "
	}

	message += statusCounts
	message += "<br>#############################################################</div><br>"

	for results.Next() {
		//make this a label for now
		err := results.Scan(&video.ID, &video.YT_videoid, &video.Title, &video.Description, &video.Publisher, &video.Publish_date, &video.Downloaded)
		functions.CheckErr(err, "Unable to get latest videos (monitor window)")
		channelInfo := database.GetChanInfo(config.Db_name, video.Publisher)
		message += "<br><span class=\"font-weight:bold; font-size:14px; text-decoration: underline; color:red\">Publisher:</span> <span class=\"font-style:italic; text-decoration:underline; color:green\">" + channelInfo.Displayname + "</span><br>"
		message += "<a href='" + YtWatchPrefix + video.YT_videoid + "'>" + video.Title + "</a> -- published: <span class=\"font-style:italic; text-decoration:underline; color:green\">" + functions.DateConvert(video.Publish_date) + "</span><hr>"
	}

	textarea.SetHtml(message)

	return textarea
}
