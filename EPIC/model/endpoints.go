package model

import (
	"encoding/json"
	"epic/jutil/v1/file"
	"epic/jutil/v1/text"
	"log"
	"path/filepath"
	"strings"
)

//EndPoint holds the the name and endpoint
type EndPoint struct {
	Name  string //localhost
	URL string //https://127.0.0.1:8888
	fname string
}

//NewEndPoint constructs a new EndPoint
func NewEndPoint(Name, URL string) *EndPoint {
	Name = text.RemoveSpaces(strings.Replace(Name ,":","-",-1))
	URL= text.RemoveSpaces(URL)

	return &EndPoint{Name: Name, URL: URL}
}

func (ep *EndPoint) Save(dir, ext string)  {
	fname := filepath.Join(dir,ep.Name +"." + ext)
	bytes, err :=json.Marshal(ep)
	if err != nil {
		log.Println(err)
	}
	file.WriteFile(fname, string(bytes))
}

func (ep *EndPoint) ToList(fnames []string) *[]EndPoint {
	var list []EndPoint
	var obj EndPoint
	for _, fname := range fnames {
		bytes := file.ReadFile(fname)
		json.Unmarshal(bytes, &obj)
		list = append(list, obj)
	}
	return &list
}









//
//
//
//
//
//
////ActionVerb holds the the http verb and endpoint
//type ActionVerb struct {
//	Command  string //GET... WAIT
//	URL string
//}
//
////NewActionVerb constructs a new ActionVerb
//func NewActionVerb(command, url string) *ActionVerb {
//	return &ActionVerb{Command: command, URL: url}
//}
//
//func SaveActionVerb(fname string, av *[]ActionVerb)  {
//	bytes, _ :=json.Marshal(av)
//	file.WriteFile(fname, string(bytes))
//}
//
//func LoadActionVerb(fname string, av *[]ActionVerb)  *[]ActionVerb {
//	file, _ := ioutil.ReadFile(fname)
//	if err := json.Unmarshal([]byte(file), *av); err != nil {
//		log.Fatal(err)
//	}
//	return av
//}
//
