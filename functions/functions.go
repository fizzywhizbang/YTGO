package functions

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

//check if a file exists used for startup
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

//generic error notification instead of writing it 10000000 times
func CheckErr(err error, msg string) {
	if err != nil {
		log.Println(msg, err.Error())
	}

}

//convert a unix date to a human readable date
func DateConvert(unixtime int) string {
	ut := int64(unixtime)
	time := time.Unix(ut, 0)
	return time.String()
}

//open your favorite web browser
func Openbrowser(url, defbrowser string) {
	//leading space there intentionally
	urlstring := url
	if !strings.Contains(url, "https://") {
		urlstring = " https://www.youtube.com/channel/" + url
	}

	cmd := exec.Command(defbrowser, "-new-tab", urlstring)
	if err := cmd.Start(); err != nil {
		log.Fatalln("can't open browser", err)

	}

}

//convert a date and trim stuff we don't want to see
func DateConvertTrim(unixtime int, limit int) string {
	ut := int64(unixtime)
	time := time.Unix(ut, 0)
	rs := []rune(time.String())
	return string(rs[:limit])
}

//convert standard date to unix for database
func DateConvertToUnix(d string) string {
	thetime, e := time.Parse(time.RFC3339, d)
	if e != nil {
		panic("Can't parse time format")
	}
	epoch := thetime.Unix()
	return strconv.Itoa(int(epoch))

}

func ConvertYMDtoUnix(ymd string) string {
	layout := "2006-01-02"
	t, err := time.Parse(layout, ymd)
	if err != nil {
		panic(err.Error())
	}
	return strconv.Itoa(int(t.Unix()))
}

func Cleanfwatch(fwatch string) bool {
	added := fwatch + "added/"
	dir, err := os.Open(added)
	if err != nil {
		log.Println("Unable to access FolderWatch")
	}
	defer dir.Close()
	names, _ := dir.Readdirnames(-1)
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(added, name))
		if err != nil {
			log.Println("Unable to remove crawljobs")
			return false
		}
	}
	return true
}
