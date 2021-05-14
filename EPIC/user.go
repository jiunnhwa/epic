package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

//User details
type User struct {
	UserName   string
	Password   []byte
	First      string
	Last       string
	UserRole   int
	IsLoggedIn bool
}

func InitDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		panic(err)
	}
	if db == nil {
		panic("db nil")
	}
	return db
}

func CreateTable(db *sql.DB) {
	// create table if not exists
	sql_table := `
	CREATE TABLE "User" (
		"UserName"	TEXT NOT NULL UNIQUE,
		"Password"	TEXT NOT NULL,
		"First"	TEXT,
		"Last"	TEXT,
		"UserRole"	INTEGER,
		"IsLoggedIn"	INTEGER,
		PRIMARY KEY("UserName")
	)	
	`

	_, err := db.Exec(sql_table)
	if err != nil {
		panic(err)
	}
}

func AddItem(db *sql.DB, items []User) {
	sql_additem := `
	INSERT INTO User(
		UserName,
		Password,
		First,
		Last,
		UserRole
	) values(?, ?, ?, ?, ?)
	`
	stmt, err := db.Prepare(sql_additem)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	for _, item := range items {
		_, err2 := stmt.Exec(item.UserName, item.Password, item.First, item.Last, item.UserRole)
		if err2 != nil {
			panic(err2)
		}
	}
}

func ReadItem(db *sql.DB) []User {
	sql_readall := `
	SELECT 
	UserName,
	Password,
	First,
	Last,
	UserRole FROM User
	ORDER BY datetime(UserName) ASC
	`

	rows, err := db.Query(sql_readall)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var result []User
	for rows.Next() {
		item := User{}
		err2 := rows.Scan(&item.UserName, &item.Password, &item.First, &item.Last, &item.UserRole)
		if err2 != nil {
			panic(err2)
		}
		result = append(result, item)
	}
	return result
}

//FindRecordByName returns a row by name
func FindRecordByName(db *sql.DB, name, pw string) ([]User, error) {
	sql := fmt.Sprintf("SELECT UserName, Password, First, Last, UserRole FROM User WHERE UserName='%s'" /* AND Password='%s'"*/, name /*, pw*/)
	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	UserInfos := []User{}
	for rows.Next() {
		var rec User
		if err := rows.Scan(&rec.UserName, &rec.Password, &rec.First, &rec.Last, &rec.UserRole); err != nil {
			return nil, err
		}
		UserInfos = append(UserInfos, rec)
	}
	return UserInfos, nil
}

//AuthenthicateUser returns -1 if invalid, or user role
func AuthenthicateUser( name, pw string) int {
	users, err := FindRecordByName(DB,  name, pw)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	if len(users)==0 {//user not found
		return -2
	}
	fmt.Println(string(users[0].Password), pw)
	if string(users[0].Password) != pw {
		return -1
	}
	return users[0].UserRole
}