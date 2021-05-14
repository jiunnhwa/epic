package url

import (
	"log"
	"net/url"
	"strings"
)

// www.google.com -> google.com
func GetHostName(urlstr string) string {
	u, err := url.Parse(urlstr)
	if err != nil {
		log.Fatal(err)
	}
	s := strings.Split(strings.Replace(u.Hostname(), "www.", "", -1), ".")
	return (s[0])
}

//x.com -> com
func GetHostSuffix(urlstr string) string {
	u, err := url.Parse(urlstr)
	if err != nil {
		log.Fatal(err)
	}
	s := strings.Split(u.Hostname(), ".")
	return (s[len(s)-1])
}
