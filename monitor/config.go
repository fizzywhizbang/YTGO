package monitor

import (
	"encoding/json"
	"log"
	"os"
)

const (
	ytVideoInfoURL = "https://www.youtube.com/get_video_info?video_id="
	ytFeedURL      = "https://www.youtube.com/feeds/videos.xml?channel_id="
	ytWatchPrefix  = "https://www.youtube.com/watch?v="
	ytChanPrefix   = "https://www.youtube.com/channel/"
	ytSearchPrefix = "https://www.youtube.com/results?search_query="
)

func loadConfig() YTCM {
	var config YTCM

	configFile, err := os.Open(ytcmConfigFile)

	checkErr(err)
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		log.Fatal("Can't decode your json", err)
	}
	defer configFile.Close()

	return config
}

func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
