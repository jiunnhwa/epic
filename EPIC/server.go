package main

import (
	"encoding/json"
	"epic/jutil/v1/convert"
	"epic/jutil/v1/file"
	"epic/jutil/v1/http/response"
	"epic/jutil/v1/maps"
	"epic/jutil/v1/session"
	"epic/jutil/v1/text"
	"epic/model"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func ServeRoutes() error{

	//VIEWS
	http.HandleFunc("/", home)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)

	http.HandleFunc("/endpoints", endpoints)
	http.HandleFunc("/cases", cases)
	http.HandleFunc("/epics", epics)
	http.HandleFunc("/jobs", jobs)
	http.HandleFunc("/loadtester", loadtester)

	//API
	http.HandleFunc("/api/v1/epic/endpoints", apiEndpoints)
	http.HandleFunc("/api/v1/epic/cases", apiCases)
	http.HandleFunc("/api/v1/epic/epics", apiEpics)
	http.HandleFunc("/api/v1/epic/epicscases", apiEpicsCases)
	http.HandleFunc("/api/v1/epic/jobs", apijobs)

	//WebSocket
	http.HandleFunc("/ws", wsHandler)

	//TODO
	//https works fine with the forms
	//but while trying to upgrade from ws to wss, there was some conflict
	//with errors like wsasend: an established connection was aborted by the software in your host machine.
	//this seems to be a socket level error but is quite hard to trace
	//tried troubleshoot in different ways, like changing ports, checking firewall, etc
	//but given the time constraint will have to trace later
	//or better still, rewrite websocket portion as a separate dedicated websocket server
	//as websocket uses GET and that interferes with POST forms
	//step down to using http

	//path := "certs\\"
	//log.Fatal(http.ListenAndServeTLS(":8080", path+"ssl.cert", path+"ssl.key", nil))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		return err
	}
	return nil
}

//Views are globals, loaded once, and reused.
var ViewHome = NewView([]string{filepath.Clean(filepath.Join(tplDir, "home.gohtml")), filepath.Join(tplDir, "base.gohtml")})
var ViewLogin = NewView([]string{filepath.Clean(filepath.Join(tplDir, "login.gohtml")), filepath.Join(tplDir, "base.gohtml")})
var ViewLogout = NewView([]string{filepath.Clean(filepath.Join(tplDir, "logout.gohtml")), filepath.Join(tplDir, "base.gohtml")})
var ViewCases = NewView([]string{filepath.Clean(filepath.Join(tplDir, "cases.gohtml")), filepath.Join(tplDir, "base.gohtml")})
var ViewEndPoints = NewView([]string{filepath.Clean(filepath.Join(tplDir, "endpoints.gohtml")), filepath.Join(tplDir, "base.gohtml")})
var ViewJobs = NewView([]string{filepath.Clean(filepath.Join(tplDir, "jobs.gohtml")), filepath.Join(tplDir, "base.gohtml")})
var ViewEpics = NewView([]string{filepath.Clean(filepath.Join(tplDir, "epics.gohtml")), filepath.Join(tplDir, "base.gohtml")})
var ViewLoadTester = NewView([]string{filepath.Clean(filepath.Join(tplDir, "loadtester.gohtml")), filepath.Join(tplDir, "base.gohtml")})

//handles home page, updates the view data and serve
func home(w http.ResponseWriter, r *http.Request) {
	viewData := &ViewData{PageTitle: "Home", Msg : "Welcome to API-Testing Droid."}
	ViewHome.SetViewData(viewData).ServeTemplate(w,r)
}

//handles login func
func login(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		if AuthenthicateUser(r.FormValue("username"),  r.FormValue("password")) <= 0 {
			w.Write([]byte(`<body><script>alert('Please login')</script><h4><button onclick="location.href = '/';" class="float-left submit-button" >Home</button></h4></body>`))
			return
		}
		myUser.UserName = r.FormValue("username")
		myUser.IsLoggedIn = true
		session.NewUserSession(myUser.UserName,Configuration.Session.ExpireMins)
		viewData := &ViewData{PageTitle: "Home", Msg : "Welcome " + myUser.UserName +  " to API-Testing Droid." , HasSessionID: myUser.IsLoggedIn}
		ViewHome.SetViewData(viewData).ServeTemplate(w,r)
		return
	}
	if r.Method == http.MethodGet {
		viewData := &ViewData{PageTitle: "Login",HasSessionID: myUser.IsLoggedIn }
		ViewLogin.SetViewData(viewData).ServeTemplate(w,r)
		return
	}
	http.Error(w, "Invalid action.", http.StatusBadRequest)
}

//handles logout func
func logout(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		session.DeleteUserSession(myUser.UserName)
		myUser = User{}//reset
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method == http.MethodGet {
		viewData := &ViewData{PageTitle: "Logout"}
		ViewLogout.SetViewData(viewData).ServeTemplate(w,r)
		return
	}
	http.Error(w, "Invalid action.", http.StatusBadRequest)
}

//endpoints page handles the crud of endpoints
func endpoints(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		model.NewEndPoint(r.FormValue("name"),r.FormValue("url")).Save(Configuration.Folders.EndPoint.In, Configuration.Folders.EndPoint.Ext)
	}
	viewData := &ViewData{PageTitle: "Manage EndPoints", DataURL: "/api/v1/epic/endpoints",HasSessionID: myUser.IsLoggedIn}
	ViewEndPoints.SetViewData(viewData).ServeTemplate(w,r)
}

//apiEndpoints respond as json the list of endpoints saved
func apiEndpoints(w http.ResponseWriter, r *http.Request) {
	var dt DataTable
	if r.Method == http.MethodGet {
		for _, obj := range *model.NewEndPoint("","").ToList(file.GetFileNames( Configuration.Folders.EndPoint.In, Configuration.Folders.EndPoint.Ext)) {
			dt.Data = append(dt.Data, []string{obj.Name,obj.URL,""})
		}
		response.AsJSON(w,200, dt)
		return
	}
	response.AsJSON(w, http.StatusBadRequest, "Invalid Action.")
}

//cases page handles the crud of testcases
func cases(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		tc := model.NewTestCase(file.GetMaxNextID(PrefixTestCase,TestCaseDir,ExtensionTestCase), r.FormValue("name"))
		tc.Description = r.FormValue("description")
		tc.Notes = r.FormValue("note")
		tc.Action.Command = r.FormValue("actionverb")
		tc.Action.URL = r.FormValue("actionurl")
		dur, _ := time.ParseDuration(r.FormValue("responsetime"))
		tc.Action.PassValue = *model.NewPassValue(r.FormValue("statustext"),r.FormValue("bodytext"),convert.ToInt(r.FormValue("statuscode")),dur)
		tc.Action.Payload = r.FormValue("payload")
		tc.Save(TestCaseDir,ExtensionTestCase)
	}
	viewData := &ViewData{PageTitle: "Manage Cases", DataURL: "/api/v1/epic/cases", Msg: r.URL.Path}
	ViewCases.SetViewData(viewData).ServeTemplate(w, r)
}

//apiCases respond as json the list of testcases
func apiCases(w http.ResponseWriter, r *http.Request) {
	var dt DataTable
	if r.Method == http.MethodGet {
		for _, obj := range *model.NewTestCase("","").ToList(file.GetFileNames( Configuration.Folders.TestCase.In, Configuration.Folders.TestCase.Ext)) {
				dt.Data = append(dt.Data, []string{obj.TID,obj.Name,""})
			}
		response.AsJSON(w,200, dt)
		return
	}
	response.AsJSON(w, http.StatusBadRequest, "Invalid Action.")
}

func GetFileBaseNames(filelist []string) (map[string]interface{} ){
	collections := map[string]interface{}{}
	for _, fname := range filelist {
		fn :=text.RemoveSpaces(strings.Replace(filepath.Base(fname)  ,".txt","",-1))
		collections[fn]=fn
	}
	return collections
}

//ListReplace
func ListReplace(list []string,old,new string) []string{
	var nl []string
	for _, item := range list {
		str :=text.RemoveSpaces(strings.Replace(filepath.Base(item),old,new,-1))
		nl = append(nl,str)
	}
	return nl
}

//ToCollections returns a map
func ToCollections(list []string) (map[string]interface{} ){
	collections := map[string]interface{}{}
	for _, item := range list {
		collections[item]=item
	}
	return collections
}

//epics page handles the crud of testepics
func epics(w http.ResponseWriter, r *http.Request) {
	submitVal := r.FormValue("submit")
	collections := maps.Merge (GetFileBaseNames(file.GetFileNames( Configuration.Folders.TestCase.In, Configuration.Folders.TestCase.Ext)),GetFileBaseNames( file.GetFileNames( Configuration.Folders.TestEpic.In, Configuration.Folders.TestEpic.Ext)  ))
	if r.Method == http.MethodPost {
		//CREATE TE, push to test/epic
		if strings.ToUpper(strings.TrimSpace((submitVal))) == "CREATE" { //from search box
			//j := model.NewJob("TestJob")
			TID := file.GetMaxNextID(PrefixTestEpic,TestEpicDir,ExtensionTestEpic)
			te:= model.NewTestEpic(TID, "TestEpic")

			for _, s := range r.Form["lstBox2"] {
				var testcase model.TestCase
				data := file.ReadFile(filepath.Join(Configuration.Folders.TestCase.In, s + "." + Configuration.Folders.TestCase.Ext))
				err := json.Unmarshal(data, &testcase)
				if err != nil {
					log.Println(err)
				}
				te.TestCases = append(te.TestCases,testcase)
			}
			//Save
			bytes, _ :=json.Marshal(te)
			file.WriteFile(filepath.Join(Configuration.Folders.TestEpic.In, TID +  "." + Configuration.Folders.TestEpic.Ext), string(bytes))
		}
	}
	viewData := &ViewData{PageTitle: "Manage Epics", DataURL: "/api/v1/epic/epics", DynamicMap: collections}
	ViewEpics.SetViewData(viewData).ServeTemplate(w,r)
}

//apiEpics returns a list of Epics
func apiEpics(w http.ResponseWriter, r *http.Request) {
	var dt DataTable
	if r.Method == http.MethodGet {
		for _, obj := range *model.NewTestEpic("", "").ToList(file.GetFileNames(Configuration.Folders.TestEpic.In, Configuration.Folders.TestEpic.Ext)) {
			dt.Data = append(dt.Data, []string{obj.TID, obj.Name, ""})
		}
		response.AsJSON(w, 200, dt)
		return
	}
	response.AsJSON(w, http.StatusBadRequest, "Invalid Action.")
}

//apiEpicsCases returns a list of Epics and Cases
func apiEpicsCases(w http.ResponseWriter, r *http.Request) {
	var dt DataTable
	if r.Method == http.MethodGet {
		for _, obj := range *model.NewTestCase("","").ToList(file.GetFileNames( Configuration.Folders.TestCase.In, Configuration.Folders.TestCase.Ext)) {
			dt.Data = append(dt.Data, []string{obj.TID,obj.Name,""})
		}
		for _, obj := range *model.NewTestEpic("", "").ToList(file.GetFileNames(Configuration.Folders.TestEpic.In, Configuration.Folders.TestEpic.Ext)) {
			dt.Data = append(dt.Data, []string{obj.TID, obj.Name, ""})
		}
		response.AsJSON(w,200, dt)
	}
}

//jobs page handles the crud of testjobs
func jobs(w http.ResponseWriter, r *http.Request) {
	p("jobs(L261)", myUser.UserName,session.UserHasSession(myUser.UserName)  )
	if !session.UserHasSession(myUser.UserName) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	submitVal := r.FormValue("submit")
	submitDel := r.FormValue("submitDel")
	if r.Method == http.MethodPost {
		//RUN
		if strings.ToUpper(strings.TrimSpace((submitVal))) == "RUN" { //from search box
			for _, fname := range file.GetFileNames(Configuration.Folders.TestJob.In, Configuration.Folders.TestJob.Ext) {
				job := model.NewJob("").Run(BaseDir,InDir,OutDir,filepath.Base(fname))
				viewData := &ViewData{PageTitle: "Administer Jobs", DataURL: "/api/v1/epic/jobs" , LogLines: strings.Join(job.LogLines, "\n") }
				ViewJobs.SetViewData(viewData).ServeTemplate(w,r)
				return
			}
		}
		//CREATE TJ, push to jobs/Inbox
		if strings.ToUpper(strings.TrimSpace((submitVal))) == "CREATE JOB" {
			collections := maps.Merge (GetFileBaseNames(file.GetFileNames( Configuration.Folders.TestCase.In, Configuration.Folders.TestCase.Ext)),GetFileBaseNames( file.GetFileNames( Configuration.Folders.TestEpic.In, Configuration.Folders.TestEpic.Ext)  ))
			list := GetFileBaseNames( file.GetFileNames( Configuration.Folders.EndPoint.In, Configuration.Folders.EndPoint.Ext)  )
			testItem := r.Form["lstBox2"]
			if len(testItem) >0{
				TID := file.GetMaxNextID(PrefixTestJob,TestJobDir,ExtensionTestJob)
				j := model.NewJob(TID)
				for _, s := range r.Form["lstBox2"] {
					p(r.Form["lstBox2"], s, Configuration.Folders.TestEpic.Ext, strings.HasPrefix(s, "TE") )
					if strings.HasPrefix(s, "TE") {
						var testepic model.TestEpic
						data := file.ReadFile(filepath.Join(Configuration.Folders.TestEpic.In, s +"." + Configuration.Folders.TestEpic.Ext))
						err := json.Unmarshal(data, &testepic)
						if err != nil {
							log.Println(err)
						}
						testcases := testepic.TestCases
						fmt.Println("testcases",testcases)
						for _, testcase := range testcases {
							var endpoint model.EndPoint
							dataEP := file.ReadFile(filepath.Join(Configuration.Folders.EndPoint.In, strings.TrimPrefix( strings.TrimSuffix(r.Form["lstBoxEP"][0], "/"), "/") +"." + Configuration.Folders.TestCase.Ext))
							err1 := json.Unmarshal(dataEP, &endpoint)
							if err1 != nil {
								log.Println(err)
							}
							testcase.Action.URL = strings.TrimPrefix( strings.TrimSuffix(endpoint.URL, "/"), "/")    + "/" + strings.TrimPrefix( testcase.Action.URL, "/")
							j.TestCases = append(j.TestCases,testcase)
						}
					} else {
						var testcase model.TestCase
						data := file.ReadFile(filepath.Join(Configuration.Folders.TestCase.In, s +"." + Configuration.Folders.TestCase.Ext))
						err := json.Unmarshal(data, &testcase)
						if err != nil {
							log.Println(err)
						}
						var endpoint model.EndPoint
						dataEP := file.ReadFile(filepath.Join(Configuration.Folders.EndPoint.In, strings.TrimPrefix( strings.TrimSuffix(r.Form["lstBoxEP"][0], "/"), "/") +"." + Configuration.Folders.TestCase.Ext))
						err1 := json.Unmarshal(dataEP, &endpoint)
						if err1 != nil {
							log.Println(err)
						}

						testcase.Action.URL = strings.TrimPrefix( strings.TrimSuffix(endpoint.URL, "/"), "/")    + "/" + strings.TrimPrefix( testcase.Action.URL, "/")
						j.TestCases = append(j.TestCases,testcase)
					}
				}
				j.Save(filepath.Join(Configuration.Folders.TestJob.In, TID +  "." + Configuration.Folders.TestJob.Ext))
			}
			viewData := &ViewData{PageTitle: "Administer Jobs", DataURL: "/api/v1/epic/jobs", DynamicList: list,  DynamicMap: collections, IsActionCreateJob: true}
			ViewJobs.SetViewData(viewData).ServeTemplate(w,r)
			return
		}
		//DELETE
		if strings.ToUpper(strings.TrimSpace((submitDel))) == "DELETE" {
			//TODO: delete job
			var obj model.Job
			for _, fname := range file.GetFileNames(Configuration.Folders.TestJob.In, Configuration.Folders.TestJob.Ext) {
				bytes:= file.ReadFile(fname)
				if err := json.Unmarshal(bytes, &obj); err != nil {
					fmt.Println(err)
				}
				fdel := filepath.Join(Configuration.Folders.TestJob.In, obj.Name +  "." + Configuration.Folders.TestJob.Ext)
				if err := os.Remove(fdel); err != nil {
					fmt.Println(err)
				}
			}
		}
	}
	viewData := &ViewData{PageTitle: "Administer Jobs", DataURL: "/api/v1/epic/jobs"}
	ViewJobs.SetViewData(viewData).ServeTemplate(w,r)
}

//apijobs returns a list of Jobs
func apijobs(w http.ResponseWriter, r *http.Request) {
	var dt DataTable
	if r.Method == http.MethodGet {
		for _, obj := range *model.NewJob("").ToList(file.GetFileNames( Configuration.Folders.TestJob.In, Configuration.Folders.TestJob.Ext)) {
			dt.Data = append(dt.Data, []string{obj.ID,obj.Name,""})
		}
		response.AsJSON(w, 200, dt)
		return
	}
	response.AsJSON(w, http.StatusBadRequest, "Invalid Action.")
}

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//loadtester page handles load testing
func loadtester(w http.ResponseWriter, r *http.Request) {
	//p("jobs(L373)", myUser.UserName,session.UserHasSession(myUser.UserName)  )
	//if !session.UserHasSession(myUser.UserName) {
	//	http.Redirect(w, r, "/login", http.StatusSeeOther)
	//	return
	//}
	submitVal := r.FormValue("submit")
	list := ToCollections(ListReplace( file.GetFileNames( Configuration.Folders.LoadTest.In, Configuration.Folders.LoadTest.Ext) ,".test","" ))
	if r.Method == http.MethodPost {
		//CREATE TE, push to test/epic
		if strings.ToUpper(strings.TrimSpace((submitVal))) == "RUN" {
		//Create and Run
			fname:= filepath.Join(Configuration.Folders.LoadTest.In, r.FormValue("lstBoxLoad")  + "." + Configuration.Folders.LoadTest.Ext)
			jobFile := fname
			var urls []string
			if n, err := strconv.Atoi(r.FormValue("numdroids")); err == nil {
				NumLoadTestJob = n
			} else {
				NumLoadTestJob = 0
			}
			for i := 0; i < NumLoadTestJob; i++ {
				urls = append(urls, jobFile )
			}
			go func() {
				resultsC := asyncJobRun(urls)
				//consumer channel
				for _ = range urls {
					result := <-resultsC
					resultText := fmt.Sprintf("%s status: %s\n", result.url, result.response)//extracts
					jobLogs <- resultText
				}
			}()
		}
		http.Redirect(w, r, "/loadtester", http.StatusSeeOther)//WS uses get
		return
	}
	viewData := &ViewData{PageTitle: "Manage Epics", DataURL: "/api/v1/epic/epics", DynamicList: list, WebSocketOutput : ""}
	ViewLoadTester.SetViewData(viewData).ServeTemplate(w,r)
}

var NumLoadTestJob int
var jobLogs =make(chan string, 10) // chan string
type JobResponse struct {
	url      string
	response string
}

//asyncJobRun runs goroutine on the producer channel
func asyncJobRun(urls []string) <-chan *JobResponse {
	ch := make(chan *JobResponse, len(urls)) // buffered
	for i, url := range urls {
		go func(index int, url string) {
			fmt.Printf("Fetching %d: %s \n",i, url)
			job := model.NewJob("").ReadFile(url).Exec().Save(url)
			ch <- &JobResponse{url, strings.Join(job.LogLines,"\n")}
		}(i,url)
	}
	return ch
}

//wsHandler writes out the message
func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}
	defer conn.Close()

	for i := 0; i < NumLoadTestJob; i++ {
		resultText := fmt.Sprintf("%d status: %s\n", i, time.Now())
		if err := conn.WriteMessage(websocket.TextMessage, []byte(resultText)  ); err != nil {
			fmt.Println(err)
		}
		if err := conn.WriteMessage(websocket.TextMessage, []byte(<-jobLogs)  ); err != nil {
			fmt.Println(err)
		}

	}
}
