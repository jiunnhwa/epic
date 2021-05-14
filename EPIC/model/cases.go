package model

import (
	"encoding/json"
	"epic/jutil/v1/file"
	"net/http"
	"path/filepath"
	"time"
)

//TestCase is the base test unit
type TestCase struct {
	TID         string
	Name       string
	Description string
	Notes       string
	PassStatus  string //PASS,FAIL,UNKNOWN,NEW
	Action 		Action
}

//NewTestCase constructs a new test case object with TID, Name
func NewTestCase(TID, Name string) *TestCase {
	//generate PK:TID
	return &TestCase{TID: TID, Name:Name}
}

func NewTestCaseSample(TID, Name string) *TestCase {
	tc := &TestCase{TID: TID, Name:Name, Description: "Test Description - " + TID, Notes: "Test Note - " + TID}
	tc.Action.Command = "GET"
	tc.Action.URL = "api/v1/school/course"
	tc.Action.PassValue = *NewPassValue("OK", ":",200, 10000000)
	return tc
}

func (tc *TestCase)Save(dir, ext string) *TestCase {
	fname := filepath.Join(dir,tc.TID +"." + ext)
	bytes, _ :=json.Marshal(tc)
	file.WriteFile(fname, string(bytes))
	return tc
}

func (tc *TestCase) ToList(fnames []string) *[]TestCase {
	var list []TestCase
	var obj TestCase
	for _, fname := range fnames {
		bytes := file.ReadFile(fname)
		json.Unmarshal(bytes, &obj)
		list = append(list, obj)
	}
	return &list
}



//Action holds the command
type Action struct {
	Command  string //GET... WAIT
	URL string
	Header http.Header //maps[string]string
	Payload string //JsonBody
	PassStatus      string   //PASS,FAIL
	PassValue //Assertions status code, has text contains.
	APIKEY
}

//PassValue are the response values to expect
type PassValue struct {
	StatusHasText string
	StatusHasCode int
	BodyHasText string
	ResponseTime time.Duration
}

//NewPassValue constructs the new object based on the values
func NewPassValue(StatusHasText, BodyHasText string, StatusHasCode int, ResponseTime time.Duration) *PassValue {
	return &PassValue{ StatusHasText: StatusHasText, StatusHasCode: StatusHasCode, BodyHasText:   BodyHasText, ResponseTime:  ResponseTime	}
}

//APIKEY whether the key is in param, header or as url
type APIKEY struct {
	Name string //name used
	AsParam string //https://api.nomics.com/v1/currencies/ticker?key=3990ec554a414b59dd85d29b2286dd85
	AsHeader string  //
	AsURL string //https://mocki.io/v1/e750d778-4861-498e-b00e-213314f799d6
	ResponseHasKey bool //JSON struct keys {"APIKEY":...}
}
