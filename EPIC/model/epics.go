package model

import (
	"encoding/json"
	"epic/jutil/v1/file"
)

//TestEpic holds a collection of test cases
type TestEpic struct {
	TID         string
	Name       string
	Description string
	Notes       string
	PassStatus  string //PASS,FAIL,UNKNOWN,NEW
	TestCases []TestCase

}

//NewTestEpic constructs a new test epic object with TID, Title
func NewTestEpic(TID, Name string) *TestEpic {
	return &TestEpic{TID: TID, Name:Name}
}

func (tc *TestEpic) ToList(fnames []string) *[]TestEpic {
	var list []TestEpic
	var obj TestEpic
	for _, fname := range fnames {
		bytes:= file.ReadFile(fname)
		json.Unmarshal(bytes, &obj)
		list = append(list, obj)
	}
	return &list
}
