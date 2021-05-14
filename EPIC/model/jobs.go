package model

import (

	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"epic/jutil/v1/file"
	"epic/jutil/v1/gen"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

//Job is a record to hold job details
type Job struct {
	//Basic
	ID                 string
	Name               string
	StartTime, EndTime time.Time
	time.Duration

	//Body
	TestCases []TestCase

	LogLines []string //holds the printlines

	//Request            string
	//Response           string
	//Result             string
	//Error              Status
}

//NewJob constructor
func NewJob(name string) *Job {
	return &Job{ID: CreateJobID(), Name: name, StartTime: time.Now()}
}


func NewJobSample(name string) *Job {
	j := &Job{ID: CreateJobID(), Name: name, StartTime: time.Now()}
	j.TestCases = append(j.TestCases,*NewTestCaseSample("TC1","Test Case Sample 1"))
	j.TestCases = append(j.TestCases,*NewTestCaseSample("TC2","Test Case Sample 2"))
	return j
}

func (j *Job) WriteFileJSON(ext string) *Job {
	bytes, _ :=json.Marshal(j)
	file.WriteFile(j.Name + "." + ext, string(bytes))
	return j
}

//Sequence generator
var nextNum = gen.GetNextNum(0)

//CreateJobID creates a formatted jobID string. eg: 20210410-0001
func CreateJobID() string {
	//return fmt.Sprintf("%s", time.Now().Format("20060102-150405.000"))  //OptionA: 20210410-225106.481.json

	//Using alternate method with running sequence for each day.
	return fmt.Sprintf("%s-%04d", time.Now().Format("20060102"), nextNum()) //OptionB: 20210410-0001.json
}




//+------------------------------------------------------------------+
//
//
//+------------------------------------------------------------------+


//ReadFile the test job file which defines the test definitions and actions
func (j *Job ) ReadFile(fname string) *Job {
	//data, err := ioutil.ReadFile(fname)
	data := file.ReadFile(fname)
	//if err != nil {
	//	fmt.Println("ReadFile error for:", fname )
	//	panic(err)
	//}
	err := json.NewDecoder(strings.NewReader(string(data))).Decode(&j)
	if err != nil {
		fmt.Println(err)
		//panic(err)
	}
	return j
}

//Run will perform load, exec,  saves, and archive the job file
func (j *Job) Run(baseDir,inDir,outDir,fname string) *Job {
	inFile, outFile:= filepath.Join(baseDir,inDir,fname),filepath.Join(baseDir,outDir,fname)
	j.ReadFile(inFile) //Load
	j.Exec()//Run
	j.Save(inFile)//Save Job with updated loglines
	//bytes, _ :=json.Marshal(j)
	//file.WriteFile(inFile, string(bytes))//Save Job with updated loglines
	file.RenameFile(inFile,outFile)
	return j
}

//Save saves the job file
func (j *Job) Save(fname string) *Job {
	bytes, _ :=json.Marshal(j)
	//file.WriteFile(filepath.Join(TestJobDir, TID +  "." + ExtensionTestJob), string(bytes))
	file.WriteFile(fname, string(bytes))//Save Job
	return j
}

func  (j *Job) ToList(fnames []string) *[]Job {
	var list []Job
	var obj Job
	for _, fname := range fnames {
		bytes:= file.ReadFile(fname)
		json.Unmarshal(bytes, &obj)
		list = append(list, obj)
	}
	return &list
}



//Exec executes the sequences of Actions.
func (j *Job) Exec() *Job {
	fmt.Println("=========== JobExec ", j.ID, " =============")
	ApikeyValue := "" //Replacer, allows get apikey to be dynamically extracted, not inputted.
	for k, v := range j.TestCases {
		j.LogLines = append(j.LogLines, "-------------------------------------------")
		j.LogLines = append(j.LogLines, v.TID + " " + v.Action.Name)
		fmt.Println("========================")
		fmt.Println(k,v, ApikeyValue)
		if ApikeyValue != "" {
				v.Action.APIKEY.AsParam = ApikeyValue //forced Replace
		}
		func (action Action)  {
			j.LogLines = append(j.LogLines, action.Command + " " + action.URL)
			j.LogLines = append(j.LogLines, "Payload: "+ " " + action.Payload)
			var dur time.Duration
			if strings.Contains(strings.ToUpper(strings.TrimSpace(action.Command)), "WAIT") {
				dur, _ = time.ParseDuration(strings.ToLower(strings.Replace(strings.Replace(strings.Replace(strings.ToUpper(strings.TrimSpace(action.Command)), "WAIT", "", -1), "(", "", -1), ")", "", -1)))
				time.Sleep(dur)
				return
			}
			if action.Command == "" {
				fmt.Println("emptiness cannot be tested.")
				return
			}
			//CourseID
			if strings.TrimSpace(action.Payload) != "" {
				action.URL = strings.TrimSpace( strings.TrimSuffix(action.URL, "/")) + "/" + ParseValue(action.Payload, "ID")
				fmt.Println("Adding CourseID", action.URL)
			}
			//Cat the key
			//if !strings.Contains(action.URL,"APIKEY")  {
			//	action.URL = action.URL + "?APIKEY="+action.AsParam
			//}
			if action.ResponseHasKey==false {
				action.AsParam = ApikeyValue
				action.URL = action.URL + "?APIKEY="+ApikeyValue //action.AsParam
			}
			//Fetch Response
			fmt.Println("Fetch:", action.URL)
			j.LogLines = append(j.LogLines, "Actual URL: "+ " " + action.URL)
			wc := NewWebclient(true).Fetch(action)
			if action.ResponseHasKey==true { //Get APIKEY command, so parse the response for apikey
				ApikeyValue = fmt.Sprintf("%v", wc.Result.JSON["APIKEY"])
				return
			}
			//else {//Action Command, append APIKEY
			//	action.AsParam = ApikeyValue
			//}
			//Test
			fmt.Println(wc.Result.String)
			fmt.Println(wc.Result.Status.Text,wc.TimeSpan.Duration)
			fmt.Println(wc.Result.Header)
			j.LogLines = append(j.LogLines, wc.Result.Status.Text,wc.TimeSpan.Duration.String())
			j.LogLines = append(j.LogLines, wc.Result.String)

			if len(action.BodyHasText)>0  {
				if strings.Contains(wc.Result.String, action.BodyHasText) {
					fmt.Println("PASS")
					action.PassStatus = "PASS"
				} else {
					fmt.Println("FAIL")
					action.PassStatus = "FAIL"
				}
			}
			j.LogLines = append(j.LogLines, action.PassStatus)

		}(v.Action)
	}

	for i, action := range j.TestCases {
		if action.PassStatus == "FAIL" {
			j.TestCases[i].PassStatus += fmt.Sprintf("FAIL #%d",i)
		}
	}
	//j.TestStatus = fmt.Sprintf("PASS ALL %d TESTS.",len(j.Actions ))
	return j
}


//ParseValue parses for the field ID
func ParseValue(payload, name string) string {
	fmt.Println("ParseValue", payload)
	var result map[string]string
	err := json.NewDecoder(strings.NewReader(payload)).Decode(&result)
	if err != nil {
		panic(err)
	}
	if val,ok := result[name]; ok {
		return val
	}
	return ""
}


//+------------------------------------------------------------------+
//
//
//+------------------------------------------------------------------+



type Webclient struct {
	Headers map[string]string
	//ResponseData data
	//errors
	client *http.Client
	Result
	TimeSpan
}

func NewWebclient(bSkipVerify bool) *Webclient {
	////True will allow self-signed ssl certs, trusting the augmented cert pool in our client
	config := &tls.Config{InsecureSkipVerify: bSkipVerify}
	tr := &http.Transport{TLSClientConfig: config}
	return &Webclient{client : &http.Client{Transport: tr}}
}

func (wc *Webclient)Fetch(action Action) *Webclient {
	wc.TimeSpan.Start = time.Now()
	defer func() {
		wc.TimeSpan.End = time.Now()
		wc.TimeSpan.Duration = time.Since(wc.TimeSpan.Start)
	}()

	//Append APIKEY as URL,
	if action.APIKEY.AsURL != "" {
		//action.URL  = path.Join(action.URL ,action.APIKEY.AsURL) //path.join gives https:///... resulting in err http: no Host in request URL. WTF?
		action.URL = strings.TrimSpace( strings.TrimSuffix(action.URL, "/")) + "/" + action.APIKEY.AsURL  //DIY
	}

	//Make Request
	req, _ := http.NewRequest(action.Command, action.URL, bytes.NewBuffer([]byte(action.Payload)))
	req.Header.Set("user-agent", "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.2; .NET CLR 1.0.3705;)")
	for k, v := range wc.Headers {
		req.Header.Set(k,v)
	}

	//Append APIKEY if in Param
	q := req.URL.Query()
	if action.APIKEY.AsParam != "" {
		q.Add(action.APIKEY.Name , action.APIKEY.AsParam )
		req.URL.RawQuery = q.Encode()
	}

	//DO
	resp, err := wc.client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	//Store in Result the following: Response Status,Headers,Body
	wc.Result.Status.Text =resp.Status
	//Copy Headers
	//for k, vv := range resp.Header {
	//	for _, v := range vv {
	//		wc.Result.Headers.Add(k, v)
	//	}
	//}

	wc.Result.Header = resp.Header

	wc.Result.Bytes, err = ioutil.ReadAll(bufio.NewReader(resp.Body))
	wc.Result.String = string(wc.Result.Bytes)

	json.Unmarshal(wc.Result.Bytes, &wc.Result.JSON)//maps[string]interface{}

	return wc
}

func  (wc *Webclient) Perform(method, URL, body, contentType string)  *Webclient {
	req, _ := http.NewRequest(method, URL, bytes.NewBuffer([]byte(body)))
	req.Header.Set("user-agent", "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.2; .NET CLR 1.0.3705;)")
	req.Header.Set("Content-Type", contentType)
	resp, err := wc.client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	//ioutil.ReadAll(bufio.NewReader(resp.Body))
	return wc
}
