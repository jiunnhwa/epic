package text

import (
	"strconv"
	"strings"
)


//type Text struct {
//	string
//}

//var String = &Text{}

func Enquote(str string) string {
	if strings.HasPrefix(str, "\"") {
		return str
	}
	return strconv.Quote(str)
	//return fmt.Sprintf("%q", str)
}

func RemoveSpaces(str string) string {
	return strings.Replace(str, " ", "", -1)
}
