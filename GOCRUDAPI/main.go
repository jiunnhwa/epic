package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"database/sql"

	"github.com/gorilla/mux"

	"gocrudapi/model/course"
	"gocrudapi/model/trainer"
	mylogger "gocrudapi/service/logger"
	"gocrudapi/service/request"
	"gocrudapi/service/response"
	"gocrudapi/service/session"

	_ "github.com/go-sql-driver/mysql"
)

//home - landing page.
func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the GO-CRUD-REST-API! Please register an APIKEY to begin.")
}

//courses respond as json all courses
//this is a core func, so each call is captured in a job record.
//logging is also done.
func courses(w http.ResponseWriter, r *http.Request) {
	job := NewJob(r.Method + "/" + request.GetAction(r))
	mylogger.Trace.Println("START:", job.ID)
	defer func() {
		mylogger.Trace.Println("END:", job.ID)
		job.End()
		job = nil
	}()

	results, err := job.Capture(course.GetCourses(DB))

	if err != nil {
		mylogger.Error.Println(err)
		response.AsJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	response.AsJSON(w, 200, results)
}

//getCcourse respond as json course info by courseid
//this is a core func, so each call is captured in a job record.
//logging is also done.
func getCourse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseid := vars["courseid"]

	job := NewJob(r.Method + "/" + request.GetAction(r) + "/" + courseid)
	mylogger.Trace.Println("START:", job.ID)
	defer func() {
		mylogger.Trace.Println("END:", job.ID)
		job.End()
		job = nil
	}()

	//GET
	if r.Method == http.MethodGet {
		results, err := job.Capture(course.GetCourseByID(DB, courseid))
		if err != nil {
			mylogger.Error.Println(err)
			response.AsJSON(w, http.StatusBadRequest, &Status{Code: -1, Message: err.Error()})
			return
		}
		response.AsJSON(w, 200, results)
	}
	//POST: Add
	if r.Method == http.MethodPost {
		var info course.CourseInfo
		(json.NewDecoder(r.Body).Decode(&info))
		fmt.Fprintf(os.Stdout, "Req Json: %+v", info)
		if courseid != string(info.ID) {
			response.AsJSON(w, http.StatusBadRequest, "courseid != string(info.ID) ")
			return
		}
		results, err := job.Capture(course.AddCourseByID(DB, string(info.ID), info.Title, info.Description))
		if err != nil {
			mylogger.Error.Println(err)
			response.AsJSON(w, http.StatusBadRequest, &Status{Code: -1, Message: err.Error()})
			return
		}
		obj := struct {
			RecordsAffected interface{}
		}{
			RecordsAffected: results,
		}
		response.AsJSON(w, 200, obj)
	}
	//PUT: Update
	if r.Method == http.MethodPut {
		var info course.CourseInfo
		(json.NewDecoder(r.Body).Decode(&info))
		fmt.Fprintf(os.Stdout, "Req Json: %+v", info)
		if courseid != string(info.ID) {
			response.AsJSON(w, http.StatusBadRequest, "PUT:"+"courseid != string(info.ID) "+courseid+" / "+fmt.Sprint(info.ID))
			return
		}
		results, err := job.Capture(course.UpdateCourseByID(DB, string(info.ID), info.Title, info.Description))
		if err != nil {
			mylogger.Error.Println(err)
			//response.AsJSON(w, http.StatusBadRequest, err.Error())
			response.AsJSON(w, http.StatusBadRequest, &Status{Code: -1, Message: err.Error()})
			return
		}
		obj := struct {
			RecordsAffected interface{}
		}{
			RecordsAffected: results,
		}
		response.AsJSON(w, 200, obj)
	}
	//DELETE
	if r.Method == http.MethodDelete {
		results, err := job.Capture(course.DeleteCourseByID(DB, courseid))
		if err != nil {
			mylogger.Error.Println(err)
			response.AsJSON(w, http.StatusBadRequest, &Status{Code: -1, Message: err.Error()})
			return
		}
		obj := struct {
			RecordsAffected interface{}
		}{
			RecordsAffected: results,
		}
		response.AsJSON(w, 200, obj)
	}
}

//getKey returns an APIKEY and message with its terms of use, or an error message
//this is a core func, so each call is captured in a job record.
//logging is also done.
func getKey(w http.ResponseWriter, r *http.Request) {
	sess, _ := session.NewSessionKey(request.GetUserIP(r))
	results := struct {
		APIKEY     string
		ClientIP   string
		Message    string
		ExpireTime time.Time
	}{
		APIKEY:     sess.SessID,
		ClientIP:   sess.ClientIP,
		Message:    "APIKEY will expire after " + fmt.Sprint(sess.Duration) + " mins of idle time; or after " + fmt.Sprint(sess.Grant) + " use.",
		ExpireTime: sess.ExpireTime,
	}
	fmt.Println(sess, results)
	response.AsJSON(w, 200, results)
}

//sessions returns as json the list of active sessions.
//this func is for private internal use, and is not to be published to external customer.
func sessions(w http.ResponseWriter, r *http.Request) {
	//SECURE
	mylogger.Trace.Println(r.Method + "/" + request.GetAction(r))
	if ok := request.IsIPLocalHost(r); ok != true { //Base implementation, alternate is to have an allowed list of IPs.
		response.AsJSON(w, http.StatusForbidden, "Internal use only.") //Alternate is to do a response.Redirect(...)
	}
	response.AsJSON(w, 200, session.GetSessions())
}

//sessions returns as json the list of active sessions.
//this func is for private internal use, and is not to be published to external customer.
//allows support to debug by having access to a Job file.
func getLog(w http.ResponseWriter, r *http.Request) {
	//SECURE
	mylogger.Trace.Println(r.Method + "/" + request.GetAction(r))
	if ok := request.IsIPLocalHost(r); ok != true { //Base implementation, alternate is to have an allowed list of IPs.
		response.AsJSON(w, http.StatusForbidden, "Internal use only.") //Alternate is to do a response.Redirect(...)
	}
	vars := mux.Vars(r)
	logid := vars["logid"]
	//GET
	data, err := ioutil.ReadFile("./logs/" + logid + ".json")
	if err != nil {
		mylogger.Error.Println(err)
		response.AsJSON(w, http.StatusBadRequest, "log file not found.")
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") //already json, else will get doubly escaped/encoded
	w.Write(data)
}

//trainers respond as json the lost of trainers
//logging is also done.
func trainers(w http.ResponseWriter, r *http.Request) {
	results, err := trainer.GetTrainers(DB)
	fmt.Println(results, len(results))
	if err != nil {
		mylogger.Error.Println(err)
		response.AsJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	response.AsJSON(w, 200, results)
}

//logPath is a middleware with single responsibility of logging the request path
func logPath(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mylogger.Trace.Println(r.URL.Path)
		f(w, r)
	}
}

//logPath is a middleware with single responsibility of logging the duration of the job
func logDuration(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() { log.Println("Duration:", time.Since(start)) }()
		f(w, r)
	}
}

//checkAPIKey is a middleware that handles as a single 'transaction block' of validating the user key.
//only when all tests are passed, will the func be allowed to be called.
//checks include keys validity, session validity, locational session validity
//session expiration is also extened upon pass
func checkAPIKey(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		userKey := q.Get("APIKEY")

		//APIKEY empty
		if len(strings.Fields(userKey)) == 0 {
			response.AsJSON(w, http.StatusBadRequest, "Empty APIKEY. Please provide a valid key.")
			return
		}

		sess, err := session.FindSession(userKey)
		//APIKEY does not match
		if err != nil {
			response.AsJSON(w, http.StatusBadRequest, "Invalid key:"+err.Error()+" Please provide a valid, non-expired APIKEY.")
			return
		}
		//APIKEY has expired.
		if time.Now().After(sess.ExpireTime) {
			response.AsJSON(w, http.StatusBadRequest, "Expired key. Please renew APIKEY.")
			return
		}
		//APIKEY use exceeds Grant.
		if sess.Grant <= 0 {
			msg := struct {
				Message string
			}{
				Message: "Grant usage exceeded. Please renew APIKEY.",
			}
			response.AsJSON(w, http.StatusBadRequest, msg)
			return
		}
		//APIKEY from multi-location sessions, only single location sessions allowed as this is a temp key.
		hostS, _, _ := net.SplitHostPort(sess.ClientIP)        //in-session ip address
		hostR, _, _ := net.SplitHostPort(request.GetUserIP(r)) //req ip address
		fmt.Println("hostS", hostS, "ClientIP:", sess.ClientIP)
		fmt.Println("hostR", hostR)
		if hostS != hostR && hostS != "" /*NewEmpty*/ {
			response.AsJSON(w, http.StatusBadRequest, "Only single location sessions allowed.")
			return
		}

		//Passed. Renew session for another ExpireMins
		sess, _ = session.RenewSession(string(sess.SessID), hostR)

		//call f
		f(w, r)
	}
}

//*************************************************************
// Job Control
//*************************************************************
var NOW time.Time //Sync timestamping, loggers can use global NOW, or time.Now()

//Status is a record to hold status information
type Status struct {
	Code    int
	Message string
}

//Job is a record to hold job details
type Job struct {
	ID                 string
	Name               string
	StartTime, EndTime time.Time
	Request            string
	Response           string
	Result             string
	Error              Status
}

//NewJob constructor
func NewJob(name string) *Job {
	NOW = time.Now()
	return &Job{ID: CreateJobID(), Name: name, StartTime: NOW}
}

//Sequence generator
var nextNum = GetNextNum(0)

//CreateJobID creates a formatted jobID string. eg: 20210410-0001
func CreateJobID() string {
	//return fmt.Sprintf("%s", time.Now().Format("20060102-150405.000"))  //OptionA: 20210410-225106.481.json

	//Using alternate method with running sequence for each day.
	return fmt.Sprintf("%s-%04d", time.Now().Format("20060102"), nextNum()) //OptionB: 20210410-0001.json
}

//GetNextNum uses a closure to generate sequences, with reset each day
func GetNextNum(startNum int) func() int {
	Day := time.Now().Day()
	currNum := startNum
	return func() int {
		if time.Now().Day() != Day {
			Day = time.Now().Day()
			return currNum
		}
		currNum += 1
		return currNum
	}
}

//Capture stores in the Result field information of interest marshalled as json string
func (j *Job) Capture(result interface{}, err error) (interface{}, error) {
	data := struct {
		Result interface{}
		Err    error
	}{
		Result: result,
		Err:    err,
	}
	bytes, _ := json.Marshal(data)
	j.Result = string(bytes)
	return result, err
}

//Save persists the Job File to disk in the log dir. eg: \logs\20210413-0001.json
func (j *Job) Save() {
	bytes, _ := json.Marshal(j)
	err := ioutil.WriteFile(path.Join("logs", j.ID+".json"), bytes, 0666)
	if err != nil {
		log.Print(err)
	}
}

//End marks the end time, and saves the job.
func (j *Job) End() {
	j.EndTime = time.Now()
	j.Save()
}

//*************************************************************
//
//*************************************************************
var DB *sql.DB

var Config = struct {
	APPName string `default:"app name"`

	DB struct {
		Name     string
		User     string `default:"root"`
		Password string `required:"true" env:"DBPassword"`
		Host     string `default:"127.0.0.1"`
		Port     uint   `default:"3306"`
	}
}{}

func main() {
	mylogger.LogDir = "logs"
	mylogger.InitLogger()
	mylogger.GoStartLogRotator()

	DB = OpenDB()
	defer CloseDB()

	go session.AutoDeleteExpiredSessions()

	path := "certs\\"
	log.Fatal(http.ListenAndServeTLS(":5000", path+"ssl.cert", path+"ssl.key", ServeRoutes()))

}

func OpenDB() *sql.DB {
	defer mylogger.Trace.Println("Database opened")
	dbname := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	if dbname == "" || host == "" || port == "" || username == "" || password == "" {
		fmt.Println("set environment issue.")
		os.Exit(1)
	}
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbname)
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func CloseDB() {
	mylogger.Trace.Println("Database closing")
	if err := DB.Close(); err != nil {
		panic(err.Error())
	}
}

func ServeRoutes() *mux.Router {
	router := mux.NewRouter()

	//user
	router.HandleFunc("/api/v1/", home)
	router.HandleFunc("/api/v1/goschool/getkey", getKey).Methods("GET")

	router.HandleFunc("/api/v1/goschool/course/{courseid}", logDuration(logPath(checkAPIKey(getCourse)))).Methods("GET", "PUT", "POST", "DELETE")
	router.HandleFunc("/api/v1/goschool/courses", (logPath((courses)))).Methods("GET")

	router.HandleFunc("/api/v1/goschool/trainers", trainers).Methods("GET")

	//for internal support use
	router.HandleFunc("/api/v1/goschool/log/{logid}", getLog).Methods("GET")
	router.HandleFunc("/api/v1/goschool/sessions", sessions).Methods("GET")

	mylogger.Trace.Println("Listening at port 5000")
	return router
}
