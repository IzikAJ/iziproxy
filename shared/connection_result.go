package shared

import (
	"encoding/json"
	"fmt"
)

// ConnectionResult - dump of http request
type ConnectionResult struct {
	Status  string `json:"status"`
	Scope   string `json:"scope"`
	Message string `json:"message,omitempty"`
}

func (conn ConnectionResult) printTab(prefix string) {
	fmt.Printf("%vConnectionResult: {\n", prefix)
	fmt.Printf("%v%10v: %v\n", prefix, "Status", conn.Status)
	fmt.Printf("%v%10v: %v\n", prefix, "Scope", conn.Scope)
	if conn.Message != "" {
		fmt.Printf("%v%10v: %q\n", prefix, "Message", conn.Message)
	}
	fmt.Printf("%v}\n", prefix)
}

func (conn ConnectionResult) print() {
	conn.printTab("")
}

func (conn ConnectionResult) dump() ([]byte, error) {
	return json.Marshal(conn)
}

// ConnectionResultFromDump - return data from dump
func ConnectionResultFromDump(dump []byte) (conn ConnectionResult, err error) {
	err = json.Unmarshal(dump, &conn)
	return
}
