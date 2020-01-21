package shared

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

// Message - tcp client message
type Message struct {
	Command int    `json:"cmd"`
	Data    []byte `json:"-"`
}

// Print is used for debug messages in console
func (msg Message) Print() {
	fmt.Printf("MESSAGE: {\n")
	msg.explain()
	fmt.Printf("}\n")
}

func encode(data []byte) (raw []byte) {
	return []byte(base64.StdEncoding.EncodeToString(data))
}

func decode(data []byte) (raw []byte) {
	raw, _ = base64.StdEncoding.DecodeString(string(data))
	return raw
}

func (msg Message) explain() {
	switch msg.Command {
	case CommandPing:
		fmt.Printf("  command: %q - %v\n", msg.Command, "ping")
	case CommandPong:
		fmt.Printf("  command: %q - %v\n", msg.Command, "pong")
	case CommandSetup:
		fmt.Printf("  command: %q - %v\n", msg.Command, "setup")
		data, _ := ConnectionSetupFromDump(msg.Data)
		data.printTab("  ")
	case CommandReady:
		fmt.Printf("  command: %q - %v\n", msg.Command, "ready")
		data, _ := ConnectionResultFromDump(msg.Data)
		data.print()
	case CommandFailed:
		fmt.Printf("  command: %q - %v\n", msg.Command, "failed")
	case CommandRequest:
		fmt.Printf("  command: %q - %v\n", msg.Command, "request")
		req, _ := RequestFromDump(msg.Data)
		fmt.Printf("  data: %q\n", req)
	case CommandResponse:
		fmt.Printf("  command: %q - %v\n", msg.Command, "response")
		req, _ := RequestFromDump(msg.Data)
		fmt.Printf("  data: %q\n", req)
	default:
		fmt.Printf("  command: %q - %v\n", msg.Command, "UNKNOWN CMD")
	}
	// return time.Now().String()
}

// Time - message tiemstamp
func (msg Message) Time() string {
	return time.Now().String()
}

func (msg Message) dump() ([]byte, error) {
	return json.Marshal(msg)
}

// recognizeError is a trivial implementation of error.
type recognizeError struct {
	Message string
}

func (e *recognizeError) Error() string {
	return e.Message
}

func (msg Message) result() (conn ConnectionResult, err error) {
	if msg.Command != CommandReady {
		return conn, &recognizeError{"invalid"}
	}
	return ConnectionResultFromDump(msg.Data)
}

func (msg Message) error() (conn ConnectionError, err error) {
	return ConnectionErrorFromDump(msg.Data)
}

// MessageFromDump - return persed message
func MessageFromDump(data []byte) (msg Message, err error) {
	return msg, json.Unmarshal(data, &msg)
}

// MarshalJSON - custom json dump
func (msg Message) MarshalJSON() ([]byte, error) {
	type Alias Message
	raw, err := json.Marshal(&struct {
		RawData string `json:"data"`
		*Alias
	}{
		RawData: string(encode(msg.Data)),
		Alias:   (*Alias)(&msg),
	})

	return raw, err
}

// UnmarshalJSON - custom json parse
func (msg *Message) UnmarshalJSON(data []byte) error {
	type Alias Message
	aux := &struct {
		RawData string `json:"data"`
		*Alias
	}{
		Alias: (*Alias)(msg),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		fmt.Println("!!!!!!!!!!!")
		fmt.Println("Unmarshal ERROR", err)
		if len(data) > 200 {
			fmt.Println(string(data[:100]))
			fmt.Println(".........")
			fmt.Println(string(data[len(data)-100:]))
		}

		return err
	}
	(*msg).Data = decode([]byte((*aux).RawData))
	return nil
}
