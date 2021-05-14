package main

import (
	"encoding/json"
	"epic/jutil/v1/file"
	"epic/jutil/v1/http/response"
	"log"
	"net/http"
)

type Table struct {
	DataTable
	Err error
	Bytes []byte
}

type DataTable struct {
	Data [][]string `json:"data"`
}

func NewTable()  *Table{
	return &Table{}
}

func (t *Table) GetData(url string) *Table  {
	resp, err  := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	t.Err = json.NewDecoder(resp.Body).Decode(&t.DataTable)
	return t
}

func (t *Table) LoadData(fname string) *Table  {
	if t.Bytes = file.ReadFile(fname); t.Err != nil {
		return t
	}
	t.Err = json.Unmarshal(t.Bytes, &t.DataTable)
	return t
}

func (t *Table) RespondAsJson(w http.ResponseWriter, r *http.Request) *Table  {
	if t.Err !=nil {
		http.Error(w, t.Err.Error(), http.StatusBadRequest)
		return t
	}
	response.AsJSON(w,200, t.DataTable)
	return t
}

