package util

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Params map[string]interface{}

type Request struct {
	URL     string      // Raw URL string
	Method  string      // HTTP method to use
	Params  *Params     // URL query parameters
	Payload interface{} // Data to JSON-encode and POST

	// Result is a pointer to a data structure. On success (HTTP status < 300),
	// response from server is unmarshaled into Result.
	Result interface{}

	// Error is a pointer to a data structure. On error (HTTP status >= 300),
	// response from server is unmarshaled into Error.
	Error interface{}

	// Optional
	Userinfo *url.Userinfo
	Header   *http.Header

	// The following fields are populated by Send().
	timestamp time.Time      // Time when HTTP request was sent
	status    int            // HTTP status for executed request
	response  *http.Response // Response object from http package
	body      []byte         // Body of server's response (JSON or otherwise)
}

// A Response is a Request object that has been executed.
type Response Request

// Timestamp returns the time when HTTP request was sent.
func (r *Response) Timestamp() time.Time {
	return r.timestamp
}

// RawText returns the body of the server's response as raw text.
func (r *Response) RawText() string {
	return strings.TrimSpace(string(r.body))
}

// Status returns the HTTP status for the executed request, or 0 if request has
// not yet been sent.
func (r *Response) Status() int {
	return r.status
}

// HttpResponse returns the underlying Response object from http package.
func (r *Response) HttpResponse() *http.Response {
	return r.response
}

// Unmarshal parses the JSON-encoded data in the server's response, and stores
// the result in the value pointed to by v.
func (r *Response) Unmarshal(v interface{}) error {
	return json.Unmarshal(r.body, v)
}

func (r *Response) String() string {
	if r == nil {
		return ""
	}
	return strings.TrimSpace(string(r.body))
}
