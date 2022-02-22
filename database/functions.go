package database

import (
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func CheckErr(err error, msg string) {
	if err != nil {
		log.Println(msg, err.Error())
	}

}

func DateConvert(unixtime int) string {
	ut := int64(unixtime)
	time := time.Unix(ut, 0)
	return time.String()
}

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

func DateConvertTrim(unixtime int, limit int) string {
	ut := int64(unixtime)
	time := time.Unix(ut, 0)
	rs := []rune(time.String())
	return string(rs[:limit])
}

func DateConvertToUnix(d string) string {
	thetime, e := time.Parse(time.RFC3339, d)
	if e != nil {
		panic("Can't parse time format")
	}
	epoch := thetime.Unix()
	return strconv.Itoa(int(epoch))

}
