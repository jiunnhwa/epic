package model


import (
	"net/http"
	"time"
)


//Data holds the data in various types
type Data struct {
	Bytes []byte
	JSON map[string]interface{}
	String string
}

//TimeSpan holds fields to calculate duration
type TimeSpan struct {
	Start, End time.Time
	time.Duration
}

//Status fields
type Status struct {
	Text string
	Code int
}

//Result fields
type Result struct {
	//bytes []byte
	//JSON maps[string]interface{}
	//string
	Data
	//Status string
	//StatusCode int
	Status
	Header http.Header //maps[string]string
}

