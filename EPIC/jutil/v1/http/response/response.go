/*

	Package centralises funcs related to http response.
*/

package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//AsJSON writes out the header and body for a json payload.
func AsJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	fmt.Println(string(response))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
