/*

	Package implements a file logging service.
	Log file name is of the format YYYYMMDD.log (eg. 20210325.log)
	Log files is rotated on each new day
	A separate goroutine is used to initialize and manage the log rotations, using the generator pattern

	A keep-alive timer periodically writes to the log file.

	Trace logs just about anything
	Info  logs important information
	Warning logs items to be concerned
	Error   logs critical problem

*/

package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"time"
)

//LogLevel type
type LogLevel int

//LogLevel Enums
const (
	TRACE LogLevel = iota + 1
	INFO
	WARNING
	ERROR
)

//Specialised loggers
var (
	Trace   *log.Logger // Just about anything
	Info    *log.Logger // Important information
	Warning *log.Logger // Be concerned
	Error   *log.Logger // Critical problem
)

var LogDir string
var fname string //file name of log (eg. 20210323.log)
var fp *os.File  //file pointer

//FileOpen for logging
func FileOpen(filename string) (*os.File, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open file:", filename, "Error:", err)
	}
	return file, nil
}

//InitLogger bootstrapping the loggers
func InitLogger() {
	fname = path.Join(LogDir, time.Now().Format("20060102")+".log")
	fp, _ = FileOpen(fname)

	format := log.Ldate | log.Ltime
	Trace = log.New(io.MultiWriter(fp, os.Stdout), "TRACE: ", format)
	Info = log.New(io.MultiWriter(fp, os.Stdout), "INFO: ", format)
	Warning = log.New(io.MultiWriter(fp, os.Stdout), "WARNING: ", format)
	Error = log.New(io.MultiWriter(fp, os.Stderr), "ERROR: ", format)

	Trace.Println("TRACE: inited.")
	Info.Println("INFO: inited.")
	Warning.Println("WARNING: inited.")
	Error.Println("ERROR: inited.")
}

//GoStartLogRotator starts various goroutines to manage log rotation
func GoStartLogRotator() {
	//a separate goroutine to init and manage log rotations
	go generator()

	//keep-alive timer, trace prints uptime every X min
	go func() {
		startTime := time.Now().Local()
		for {
			time.Sleep(5 * time.Minute) //sleeps 5min
			str := fmt.Sprintf("%s\tUptime(days):%.2f", time.Now().Format("2006-01-02 15:04"), (time.Since(startTime).Hours() / 24))
			Trace.Println(str)
		}
		//https://gobyexample.com/tickers
		//https://qvault.io/golang/range-over-ticker-in-go-with-immediate-first-tick/
		// ticker := time.NewTicker(time.Minute)
		// for ; true; <-ticker.C {
		// 	fmt.Println("hi")
		// }
	}()
}

//generator uses the generator pattern to do daily log rotation
//a new log file is initialised on each date change
//a keep-alive timer writes to the file every 5 minutes
func generator() {
	//init channel
	c := onDateChange()

	//initLogger onDateChange
	for {
		select {
		case <-c:
			InitLogger()
		}
	}
}

//onDateChange makes a receive-only channel of int to hold the day number
//sends the new day number on each new day
func onDateChange() <-chan int { //
	c := make(chan int)
	go func() {
		start := 0
		for i := 0; ; i++ {
			day := time.Now().Day()
			if start != day {
				start = day
				c <- day
			}
			time.Sleep(time.Duration(100 * time.Millisecond))
		}
	}()
	return c // Return the channel to the caller.
}
