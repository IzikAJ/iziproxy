package shared

import (
	"encoding/json"
	"fmt"
)

// ConnectionError - dump of http request
type ConnectionError struct {
	Code    string `json:"code"`
	Message string `json:"message,omitempty"`
}

func (conn *ConnectionError) Error() string {
	return conn.Code
}

func (conn ConnectionError) printTab(prefix string) {
	fmt.Printf("%vConnectionError: {\n", prefix)
	fmt.Printf("%v%10v: %v\n", prefix, "Code", conn.Code)
	if conn.Message != "" {
		fmt.Printf("%v%10v: %q\n", prefix, "Message", conn.Message)
	}
	fmt.Printf("%v}\n", prefix)
}

func (conn ConnectionError) print() {
	conn.printTab("")
}

func (conn ConnectionError) dump() ([]byte, error) {
	return json.Marshal(conn)
}

// ConnectionErrorFromDump - return data from dump
func ConnectionErrorFromDump(dump []byte) (conn ConnectionError, err error) {
	err = json.Unmarshal(dump, &conn)
	return
}
