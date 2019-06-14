package shared

import (
	"encoding/json"
	"fmt"
)

type commander struct {
}

// Commander - simple commander
var Commander = &commander{}

func (cmd *commander) print() {
	fmt.Printf("CMD: {\n")
	fmt.Printf("  cmd: %q\n", *cmd)
	fmt.Printf("}\n")
}

// Parse - return persed message
func (cmd *commander) Parse(data []byte) (msg Message, err error) {
	err = json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println("")
		fmt.Println("")
		if len(data) > 200 {
			fmt.Println(string(data[:100]))
			fmt.Println(".........")
			fmt.Println(string(data[len(data)-100:]))
		}
		fmt.Println(err)
		fmt.Println("")
		fmt.Println("")
	}
	return
}

// MakePing - ping message
func (cmd *commander) MakePing() (msg Message, err error) {
	return Message{Command: CommandPing}, nil
}

// MakeOK - OK message
func (cmd *commander) MakeOK() (msg Message, err error) {
	return Message{Command: CommandOK}, nil
}

// MakeFail - Fail message
func (cmd *commander) MakeFail() (msg Message, err error) {
	return Message{Command: CommandFail}, nil
}

// MakeRequest - req message
func (cmd *commander) MakeRequest(req Request) (msg Message, err error) {
	raw, err := json.Marshal(req)
	if err != nil {
		return
	}
	msg = Message{Command: CommandRequest, Data: raw}
	return
}

// MakeResponse - req message
func (cmd *commander) MakeResponse(req Request) (msg Message, err error) {
	raw, err := json.Marshal(req)
	if err != nil {
		return
	}
	msg = Message{Command: CommandResponse, Data: raw}
	return
}
