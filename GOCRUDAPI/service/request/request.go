/*

	Package centralises funcs related to http request.

*/

package request

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/gorilla/mux"
)

// GetUserIP gets a requests IP address by reading off the forwarded-for
// header (for proxies) and falls back to use the remote address.
//According to developer.mozilla.org, The X-Forwarded-For (XFF) header is a de-facto standard header for identifying the originating IP address of a client connecting to a web server through an HTTP proxy or a load balancer.
func GetUserIP(r *http.Request) string {
	if forwarded := r.Header.Get("X-FORWARDED-FOR"); forwarded != "" {
		//host, port, _ := net.SplitHostPort(forwarded)
		//return host, port
		return forwarded
	}
	//host, port, _ := net.SplitHostPort(r.RemoteAddr)
	//return host, port
	return r.RemoteAddr
}

//GetAction extracts the 'action' fragment from the url path.
func GetAction(r *http.Request) string {
	vars := mux.Vars(r)
	courseid := vars["courseid"]
	if len(courseid) > 0 {
		fields := strings.Split(strings.TrimSuffix(r.URL.Path, "/"), "/")
		fmt.Println(fields)
		return fields[len(fields)-2] //action is 2nd last word
	}
	_, action := path.Split(r.URL.Path) //action is at last posn
	return action
}

//IsIPLocalHost returns true if RemoteAddr/Header Forwarded is 127.0.0.1 or localhost
func IsIPLocalHost(r *http.Request) bool {
	userIP := GetUserIP(r)
	if userIP == "127.0.0.1" || userIP == "localhost" {
		return true
	}
	return false
}
