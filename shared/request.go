package shared

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
)

// RequestHeader - one serialized HTTP header
type RequestHeader struct {
	Name  string   `json:"name"`
	Value []string `json:"value"`
}

// Request - dump of http request
type Request struct {
	ID      uuid.UUID       `json:"id"`
	Status  int             `json:"status"`
	Method  string          `json:"method,omitempty"`
	Path    string          `json:"path,omitempty"`
	Headers []RequestHeader `json:"headers,omitempty"`
	Body    []byte          `json:"body,omitempty"`
}

func (req Request) print() {
	fmt.Printf("Request: {\n")
	if req.Method != "" {
		fmt.Printf("  method:  %q\n", req.Method)
	}
	if req.Path != "" {
		fmt.Printf("  path:    %q\n", req.Path)
	}
	if req.Status > 0 {
		fmt.Printf("  status:    %q\n", req.Status)
	}
	fmt.Printf("  headers: {\n")
	for _, phdr := range req.Headers {
		fmt.Printf("    %30.30q: %q\n", phdr.Name, phdr.Value)
	}
	fmt.Println("  }")
	bodySize := len(req.Body)
	if bodySize > 100 {
		fmt.Printf("  body:    %q\n", string([]byte(req.Body)[:100]))
	} else if bodySize > 0 {
		fmt.Printf("  body:    %q\n", req.Body)
	}
	fmt.Println("}")
}

func (req Request) dump() ([]byte, error) {
	return json.Marshal(req)
}

// RequestFromRequest - return data from http request
func RequestFromRequest(r *http.Request) (Request, error) {
	headers := make([]RequestHeader, 0)
	for k, v := range r.Header {
		headers = append(headers, RequestHeader{Name: k, Value: v})
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	uuid, err := uuid.NewUUID()

	pr := Request{
		ID:      uuid,
		Method:  r.Method,
		Path:    r.RequestURI,
		Headers: headers,
		Body:    body,
	}

	return pr, err
}

// RequestFromResponse - return data from http response
func RequestFromResponse(r *http.Response) (pr Request, err error) {
	headers := make([]RequestHeader, 0)
	if r == nil {
		return
	}
	for k, v := range r.Header {
		headers = append(headers, RequestHeader{Name: k, Value: v})
	}
	body := make([]byte, 0)
	if r.Header.Get("Content-Length") != "0" {
		body, err = ioutil.ReadAll(r.Body)
		defer r.Body.Close()
	}

	pr = Request{
		Status:  r.StatusCode,
		Headers: headers,
		Body:    body,
	}

	return
}

// RequestFromDump - return data from dump
func RequestFromDump(data []byte) (req Request, err error) {
	err = json.Unmarshal(data, &req)
	return
}
