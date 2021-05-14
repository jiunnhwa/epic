package gen

import (
	"fmt"
	"time"
)

//GetNextID uses a closure to generate ID sequences
func GetNextID(prefix string, startNum int) func() string {
	currNum := startNum-1
	return func() string {
		currNum += 1
		return fmt.Sprintf("%s-%0004d",prefix,currNum)
	}
}

//GetNextNum uses a closure to generate sequences, with reset each day
func GetNextNum(startNum int) func() int {
	Day := time.Now().Day()
	currNum := startNum
	return func() int {
		if time.Now().Day() != Day {
			Day = time.Now().Day()
			return currNum
		}
		currNum += 1
		return currNum
	}
}
