package csv

import (
	"encoding/csv"
	"log"
	"os"
)

//ReadCSV a csv files and returns as an array of csv records
func ReadCSV(filename string) [][]string {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := csv.NewReader(f).ReadAll()
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	return rows
}


