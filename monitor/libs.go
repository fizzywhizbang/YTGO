package monitor

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

func CheckErr(err error, msg string) {
	if err != nil {
		log.Println(msg, err.Error())
	}

}
func dateConvertTrim(unixtime int, limit int) string {
	ut := int64(unixtime)
	time := time.Unix(ut, 0)
	rs := []rune(time.String())
	return string(rs[:limit])
}

func friendlyDate(d string) string {
	fmt.Println(d)
	date, e := time.Parse(time.RFC3339, d)
	if e != nil {
		panic("Can't parse time format")
	}
	return date.Format("2006-01-02")
}

func convertYMDtoUnix(ymd string) string {
	layout := "2006-01-02"
	t, err := time.Parse(layout, ymd)
	if err != nil {
		panic(err.Error())
	}
	return strconv.Itoa(int(t.Unix()))
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("something went wrong: " + err.Error())
	}
}
