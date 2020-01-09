package shared

import (
	"encoding/json"
	"fmt"
)

// ConnectionSetup - dump of http request
type ConnectionSetup struct {
	Token    string `json:"token"`
	Scope    string `json:"scope,omitempty"`
	Fallback bool   `json:"fallback"`
}

func (conn ConnectionSetup) printTab(prefix string) {
	fmt.Printf("%vConnectionSetup: {\n", prefix)
	fmt.Printf("%v%10v: %v\n", prefix, "Token", conn.Token)
	if conn.Scope != "" {
		fmt.Printf("%v%10v: %q\n", prefix, "Scope", conn.Scope)
	}
	fmt.Printf("%v%10v: %v\n", prefix, "Fallback", conn.Fallback)
	fmt.Printf("%v}\n", prefix)
}

func (conn ConnectionSetup) print() {
	conn.printTab("")
}

func (conn ConnectionSetup) dump() ([]byte, error) {
	return json.Marshal(conn)
}

// ConnectionSetupFromDump - return data from dump
func ConnectionSetupFromDump(data []byte) (conn ConnectionSetup, err error) {
	err = json.Unmarshal(data, &conn)
	return
}
