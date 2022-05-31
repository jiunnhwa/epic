package course

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"gocrudapi/model/trainer"
)

//Primary Key
type CourseID string

//CourseInfo table definition
type CourseInfo struct {
	ID           CourseID          `json:"ID"`
	Title        string            `json:"Title"`
	Description  string            `json:"Description"`
	PreRequsites []CourseID        `json:"PreRequsites"`
	Trainers     []trainer.Trainer `json:"Trainers"`
}

//GetCourseByID returns a row by ID
func GetCourseByID(db *sql.DB, ID string) ([]CourseInfo, error) {

	if err := SQLReject(ID); err != nil {
		return nil, err
	}

	sql := fmt.Sprintf("Select * FROM goschool.courses WHERE ID='%s'", ID)
	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	CourseInfos := []CourseInfo{}
	fmt.Println(rows)
	for rows.Next() {
		var p CourseInfo
		if err := rows.Scan(&p.ID, &p.Title, &p.Description /*, &p.PreRequsites, &p.Trainers*/); err != nil {
			return nil, err
		}
		CourseInfos = append(CourseInfos, p)
	}

	return CourseInfos, nil
}

//GetCourses returns rows limited to 30
func GetCourses(db *sql.DB) ([]CourseInfo, error) {
	rows, err := db.Query("Select * FROM goschool.courses LIMIT 30")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	CourseInfos := []CourseInfo{}
	fmt.Println(rows)
	for rows.Next() {
		var p CourseInfo
		if err := rows.Scan(&p.ID, &p.Title, &p.Description /*, &p.PreRequsites, &p.Trainers*/); err != nil {
			return nil, err
		}
		CourseInfos = append(CourseInfos, p)
	}

	return CourseInfos, nil
}

//DeleteCourseByID a course by ID
func DeleteCourseByID(db *sql.DB, ID string) (int64, error) {

	if err := SQLReject(ID); err != nil {
		return -1, err
	}

	sql := fmt.Sprintf("DELETE FROM goschool.courses WHERE ID='%s'", ID)
	result, err := db.Exec(sql)
	if err != nil {
		return -1, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ID = %s, affected = %d\n", ID, rowsAffected)

	return rowsAffected, nil
}

//AddCourseByID a course by ID
func AddCourseByID(db *sql.DB, ID, title, desc string) (int64, error) {

	if err := SQLReject(ID); err != nil {
		return -1, err
	}

	stmt, err := db.Prepare("INSERT INTO goschool.courses (`ID`, `Title`, `Description`) VALUES (?,?,?)")
	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(ID, title, desc)
	if err != nil {
		//Error 1062: Duplicate entry 'YODA123' for key 'PRIMARY'
		return -1, err
		//log.Fatal(err)
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ID = %d, affected = %d\n", lastId, rowsAffected)
	return rowsAffected, nil
}

//UpdateCourseByID the course details by the given ID
func UpdateCourseByID(db *sql.DB, ID, title, desc string) (int64, error) {
	if err := SQLReject(ID); err != nil {
		return -1, err
	}

	sql := "UPDATE goschool.courses SET "
	if len(title) > 0 {
		sql += fmt.Sprintf("Title='%s' ", StringEscape(title))
	}
	if len(desc) > 0 {
		sql += ", " + fmt.Sprintf("Description='%s' ", StringEscape(desc))
	}
	sql += fmt.Sprintf("WHERE ID='%s'", ID)
	fmt.Println(sql)
	stmt, err := db.Prepare(sql)
	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec()
	if err != nil {
		//Error 1064: You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near 'WHERE ID='3333'' at line 1
		log.Fatal(err)
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ID = %d, affected = %d\n", lastId, rowsAffected)
	return rowsAffected, nil
}

//*************************************************************
// SQL Sanitizer
//*************************************************************
func StringEscape(str string) string {
	str = strings.Replace(str, "'", "\\'", -1)
	return str
}

//SQLReject guards against SQL injection
func SQLReject(str string) error {
	fields := strings.Fields(str)
	if len(fields) > 1 {
		return errors.New("wrong number of fields in input")
	}
	return nil
}
