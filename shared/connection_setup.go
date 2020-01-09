package shared

import (
	"encoding/json"
	"fmt"
)

// ConnectionSetup - dump of http request
type ConnectionSetup struct {
	Key      string `json:"key"`
	Scope    string `json:"scope,omitempty"`
	Fallback bool   `json:"fallback"`
}

func (conn ConnectionSetup) printTab(prefix string) {
	fmt.Printf("%vConnectionSetup: {\n", prefix)
	fmt.Printf("%v  Key:  %q\n", prefix, conn.Key)
	if conn.Scope != "" {
		fmt.Printf("%v  Scope:    %q\n", prefix, conn.Scope)
	}
	fmt.Printf("%v  Fallback:    %v\n", prefix, conn.Fallback)
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
