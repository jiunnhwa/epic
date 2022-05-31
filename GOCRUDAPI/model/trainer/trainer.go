package trainer

import (
	"database/sql"
	"fmt"
)

//Trainer table definition
type Trainer struct {
	ID        int
	FirstName string
	LastName  string
	Bio       string
	Age       int
}

//GetTrainers returns rows limited to 30
func GetTrainers(db *sql.DB) ([]Trainer, error) {
	rows, err := db.Query("Select * FROM goschool.trainers LIMIT 30")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	Trainers := []Trainer{}
	fmt.Println(rows)
	for rows.Next() {
		var p Trainer
		if err := rows.Scan(&p.ID, &p.FirstName, &p.LastName, &p.Age); err != nil {
			return nil, err
		}
		Trainers = append(Trainers, p)
	}

	return Trainers, nil
}
