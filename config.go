package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/fizzywhizbang/YTGO/functions"
)

type YTGO struct {
	Db_name     string `json:"db_name"`
	BaseDL      string `json:"base_download_dir"`
	Defbrowser  string `json:"defbrowser"`
	FolderWatch string `json:"folderwatch"`
	Monitor     bool   `json:"monitor"`
}

const (
	YtVideoInfoURL = "https://www.youtube.com/get_video_info?video_id="
	YtFeedURL      = "https://www.youtube.com/feeds/videos.xml?channel_id="
	YtWatchPrefix  = "https://www.youtube.com/watch?v="
	YtChanPrefix   = "https://www.youtube.com/channel/"
	YtSearchPrefix = "https://www.youtube.com/results?search_query="
)

//check for configurations
func CkConfig() bool {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Println("Unable to get home dir")
		return false
	}
	ConfigDir = homedir + "/.config/ytgo"
	if _, err := os.Stat(ConfigDir); os.IsNotExist(err) {
		//if not exist create it
		err = os.Mkdir(ConfigDir, 0755)
		if err != nil {
			log.Println("Error creating config dir (line 29)")
			return false
		}
		if !functions.Exists(ConfigDir + "/" + ConfigFile) {
			return createConfigFile(homedir)
		}
		return true
	}
	return true
}

//config file create in JSON
func createConfigFile(homedir string) bool {
	base := homedir + "/YTGOVideos/"
	folderWatch := homedir + "/bin/JDownloader 2.0/folderwatch/"
	return writeConfig(ConfigDir+"/ytgo.db", base, "", folderWatch, "false")
}

func writeConfig(dbname, basedl, defbrowser, folderwatch, monitor string) bool {

	file, err := os.Create(ConfigDir + "/" + ConfigFile)
	if err != nil {
		log.Println("Unable to create config file (line 63)")
		return false
	}
	defer file.Close()
	fmt.Fprintln(file, "{")
	fmt.Fprintln(file, "\t\"db_name\":\""+dbname+"\",")
	fmt.Fprintln(file, "\t\"base_download_dir\":\""+basedl+"\",")
	fmt.Fprintln(file, "\t\"defbrowser\":\""+defbrowser+"\",")
	fmt.Fprintln(file, "\t\"folderwatch\":\""+folderwatch+"\",")
	fmt.Fprintln(file, "\t\"monitor\":"+monitor+"")
	fmt.Fprintln(file, "}")

	return true
}

func ConfigParser() YTGO {
	var config YTGO

	configFile, err := os.Open(ConfigDir + "/" + ConfigFile)
	if err != nil {
		log.Println("Unable to read config file (line 59)")
	}

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		log.Fatal("Can't decode your json", err)
	}
	defer configFile.Close()

	return config
}
