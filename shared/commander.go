package shared

import (
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
	return MessageFromDump(data)
}

type allowDump interface {
	dump() ([]byte, error)
}

func (cmd *commander) NewMessage(code int, items ...allowDump) (msg Message, err error) {
	var raw []byte
	for _, item := range items {
		if dumped, err := item.dump(); err == nil {
			raw = dumped
			break
		}
	}
	return Message{Command: code, Data: raw}, nil
}

// MakePing - Ping message
func (cmd *commander) MakePing() (msg Message, err error) {
	return cmd.NewMessage(CommandPing)
}

// MakePong - Pong message
func (cmd *commander) MakePong() (msg Message, err error) {
	return cmd.NewMessage(CommandPong)
}

// MakeRequest - Request message
func (cmd *commander) MakeRequest(req Request) (msg Message, err error) {
	return cmd.NewMessage(CommandRequest, req)
}

// MakeResponse - Response message
func (cmd *commander) MakeResponse(req Request) (msg Message, err error) {
	return cmd.NewMessage(CommandResponse, req)
}

// MakePing - Ping message
func (cmd *commander) MakeSetup(data ConnectionSetup) (msg Message, err error) {
	return cmd.NewMessage(CommandSetup, data)
}

// MakeReady - Pong message
func (cmd *commander) MakeReady(data ConnectionResult) (msg Message, err error) {
	return cmd.NewMessage(CommandReady, data)
}

// MakeFailed - Pong message
func (cmd *commander) MakeFailed(data ConnectionError) (msg Message, err error) {
	return cmd.NewMessage(CommandFailed, data)
}
