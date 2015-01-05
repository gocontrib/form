package form

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// NewRequest creates new POST http request with specified form data
func NewRequest(URL string, form map[string]interface{}) *http.Request {
	var data = url.Values{}
	for k, v := range form {
		// TODO support arrays
		data.Add(k, fmt.Sprintf("%v", v))
	}
	var str = data.Encode()

	var r, _ = http.NewRequest("POST", URL, bytes.NewBufferString(str))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(str)))

	return r
}
