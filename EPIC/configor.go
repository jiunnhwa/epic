package main

import (
	"encoding/json"
	"epic/jutil/v1/file"
	"strings"
)

const (
	PrefixEndPoint  = "EP"
	EndPointDir = "./test/endpoints"
	ExtensionEndPoint = "txt"

	PrefixTestCase  = "TC"
	TestCaseDir = "./test/cases"
	ExtensionTestCase = "txt"

	PrefixTestEpic  = "TE"
	TestEpicDir = "./test/epics"
	ExtensionTestEpic = "txt"

	PrefixTestJob  = "TJ"
	TestJobDir = "./jobs/in"
	ExtensionTestJob = "test"

	BaseDir = "./jobs"
	InDir = "in"
	OutDir = "done"
)

type Configor struct {
	Session Session
	Folders Folders
	FileNames FileNames
}
type Session struct {
	ExpireMins	int
}

type Folders struct {
	EndPoint, TestJob, TestCase, TestEpic Folder
	LoadTest Folder
	DB	Folder
}

type Folder struct {
	Base,In,Out, Ext string
}

type FileNames struct {
	TestJob, TestCase, TestEpic FileName
}

type FileName struct {
	Prefix, Suffix string
}

func NewConfigor() *Configor {
	return &Configor{}
}

func NewConfigorSample() *Configor {
	return &Configor{
		Session: Session{
			ExpireMins: 30,
		},
		Folders:      Folders{
			DB:  Folder{
				Base: "./",
				In:   "",
				Out:  "",
				Ext:  "db",
			},
			EndPoint:  Folder{
				Base: "./test",
				In:   "./test/endpoints",
				Out:  "",
				Ext:  "txt",
			},
			TestJob:  Folder{
				Base: "./jobs",
				In:   "./jobs/in",
				Out:  "./jobs/done",
				Ext:  "test",
			},
			TestCase: Folder{
				Base: "./test",
				In:   "./test/cases",
				Out:  "",
				Ext:  "txt",
			},
			TestEpic: Folder{
				Base: "./test",
				In:   "./test/epics",
				Out:  "",
				Ext:  "txt",
			},
		},
		FileNames: FileNames{
			TestJob: FileName{
				Prefix: "TJ",
				Suffix: "",
			},
			TestCase: FileName{
				Prefix: "TC",
				Suffix: "",
			},
			TestEpic: FileName{
				Prefix: "TE",
				Suffix: "",
			},
		},
	}
}

func (c *Configor) SaveJSON(fname string) *Configor {
	bytes, _ :=json.Marshal(c)
	file.WriteFile(fname, string(bytes))
	return c
}

func (c *Configor) LoadJSON(fname string) *Configor {
	data:= file.ReadFile(fname)
	err := json.NewDecoder(strings.NewReader(string(data))).Decode(&c)
	if err != nil {
		panic(err)
	}
	return c
}

