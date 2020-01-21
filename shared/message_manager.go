package shared

import (
	"fmt"
	"time"
)

type messageManager struct{}

func (m messageManager) ReciveMessage(conn *Connection) (msg Message, err error) {
	data, err := conn.ReadRaw()
	if err != nil {
		return
	}

	msg, err = Commander.Parse(data)
	if err != nil {
		fmt.Printf("DATA PARSE ERROR\n%q\n", err)
		time.Sleep(time.Minute)
	}
	return
}

func (m messageManager) GetRequest(msg Message) (Request, error) {
	return RequestFromDump(msg.Data)
}

func (m messageManager) SendMessage(msg Message, conn *Connection) (err error) {
	data, err := msg.dump()
	if err != nil {
		return
	}
	err = conn.WriteRaw(data)
	return
}

// messageManager - messageManager instance
var MessageManager = messageManager{}
