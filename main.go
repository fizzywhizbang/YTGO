package main

import (
	"fmt"

	"github.com/fizzywhizbang/YTGO/database"
)

var ConfigDir = ""
var ConfigFile = "ytgo.json"

func main() {
	if CkConfig() {
		//config exists so next check for the database
		config := ConfigParser()
		fmt.Println(database.DBCheck(config.Db_name))
	}
}
