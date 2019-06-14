package shared

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
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
	// fmt.Printf("  command: %q\n", msg.Command)
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
	case CommandAuth:
		fmt.Printf("  command: %q - %v\n", msg.Command, "auth")
	case CommandRequest:
		fmt.Printf("  command: %q - %v\n", msg.Command, "CommandRequest")
		req, err := RequestFromDump(msg.Data)
		// req, err := RequestFromDump(msg.Data)
		if err != nil {
			fmt.Printf("ERR \n %q\n", err)
			return
		}
		fmt.Printf("  data: %q\n", req)
	case CommandResponse:
		fmt.Printf("  command: %q - %v\n", msg.Command, "CommandResponse")
		req, err := RequestFromDump(msg.Data)
		if err != nil {
			fmt.Printf("ERR \n %q\n", err)
			return
		}
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

// ParseMessage - return persed message
func ParseMessage(data []byte) (msg Message, err error) {
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

type msgManager struct{}

func (m msgManager) ReciveMessage(conn *Connection) (msg Message, err error) {
	data, err := conn.ReadRaw()
	// data, err := ReadConnectionData(conn)
	if err != nil {
		if err == io.EOF {
			fmt.Println("LOST CONNECTION!", data)
		}
		fmt.Println("!! reciveMessage ERROR !!", err)
		return
	}

	msg, err = Commander.Parse(data)
	if err != nil {
		fmt.Printf("DATA PARSE ERROR\n%q\n", err)
		time.Sleep(time.Minute)
	}
	return
}

func (m msgManager) GetRequest(msg Message) (Request, error) {
	return RequestFromDump(msg.Data)
}

func (m msgManager) SendMessage(msg Message, conn *Connection) (err error) {
	data, err := msg.dump()
	if err != nil {
		return
	}
	err = conn.WriteRaw(data)
	return
}

// MsgManager - MsgManager instance
var MsgManager = msgManager{}
