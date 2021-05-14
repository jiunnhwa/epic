package main

import (
	"database/sql"
	"epic/jutil/v1/session"
	"path/filepath"

	//"encoding/json"
	"fmt"
	//"io/ioutil"
	"log"
	"os"
	//"strings"
	"time"

	mylogger "epic/jutil/v1/logger"
)

var p = fmt.Println //custom debug print

var NOW time.Time //TimeSyncVal, single time value

var Configuration Configor //Configs
var DB *sql.DB
var myUser User

func init() {
	//Load configs
	Configuration = *NewConfigor().LoadJSON("config.json")
	p("Configuration.Session.ExpireMins",Configuration.Session.ExpireMins)
	//Configuration.Session.ExpireMins = 60
	//Configuration.SaveJSON("config.json")
	//Load user db
	DB = InitDB(filepath.Join(Configuration.Folders.DB.Base, "user.db"))
	//Load logger
	mylogger.LogDir = "logs"
	mylogger.InitLogger()
	mylogger.GoStartLogRotator()
	//Session housekeep
	go session.AutoDeleteExpiredSessions()
}

//main supervisor block
func main() {
	defer DB.Close()
	if err := run(); err != nil {
		log.Println("error :", err)
		os.Exit(1)
	}
}
//run the serveroutes
func run() (runerr error) {
	return ServeRoutes()
}

