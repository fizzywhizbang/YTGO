package functions

import (
	"log"
	"os"
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
