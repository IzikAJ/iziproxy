package client

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/izikaj/iziproxy/shared"
	"github.com/izikaj/iziproxy/shared/names"
)

// ReconnectTimeout - timeout to reconnect attempt in seconds
const ReconnectTimeout = 3

// ReconnectTimes - count of reconnect attempts
const ReconnectTimes = 10

func (client *Client) load(req shared.Request) (resp shared.Request, err error) {
	httpReq, _ := http.NewRequest(
		req.Method,
		client.Host+req.Path,
		bytes.NewReader(req.Body),
	)

	for _, header := range req.Headers {
		httpReq.Header.Del(header.Name)
		for _, value := range header.Value {
			httpReq.Header.Add(header.Name, value)
		}
	}

	httpResp, err := client.http.Do(httpReq)

	resp, err = shared.RequestFromResponse(httpResp)
	fmt.Printf("REQ > [%d] %s:%s (%d)\n", resp.Status, req.Method, req.Path, len(resp.Body))
	resp.ID = req.ID
	return
}

func (client *Client) handleRequest(msg shared.Message) (err error) {
	req, err := shared.MessageManager.GetRequest(msg)
	fmt.Printf("REQ < %s:%s\n", req.Method, req.Path)

	resp, err := client.load(req)

	msg2, err := shared.Commander.MakeResponse(resp)
	err = shared.MessageManager.SendMessage(msg2, client.conn)
	if err != nil {
		fmt.Println("handleRequest ERROR:", err)
	}
	return
}

func (client *Client) handleIncomingMessages() (err error) {
	defer func() {
		client.signal <- err
	}()

	var msg shared.Message
	for {
		msg, err = shared.MessageManager.ReciveMessage(client.conn)
		if err != nil {
			fmt.Println("reciveMessage ERROR", err, msg)
			// client.signal <- err
			return
		}
		switch msg.Command {
		case shared.CommandPong:
			fmt.Print("<")

		case shared.CommandReady:
			fmt.Println("Connection READY")
			msg.Print()

		case shared.CommandFailed:
			fmt.Println("Connection FAILED")
			client.alive = false
			err = &names.GenerationError{S: "TEST"}
			return

		case shared.CommandRequest:
			go client.handleRequest(msg)

		case shared.CommandPing:
			msg, err = shared.Commander.MakePong()
			err = client.conn.SendMessage(msg)
			if err != nil {
				fmt.Println("PONG ERROR?", client.conn.RemoteAddr())
				return
			}

			fmt.Print(">")

		default:
			fmt.Println("RECIVED UNHANDLED MESSAGE")
			msg.Print()
		}
	}
}

func (client *Client) handle() (err error) {
	conn := client.conn
	defer conn.Close()
	conn.Init()
	var msg shared.Message

	// test setup
	msg, err = shared.Commander.MakeSetup(shared.ConnectionSetup{
		Token:    "test_key",
		Scope:    client.Space,
		Fallback: client.Fallback,
	})
	if err != nil {
		return
	}
	err = shared.MessageManager.SendMessage(msg, conn)
	if err != nil {
		return
	}

	go client.handleIncomingMessages()

	for {
		select {
		case err = <-client.signal:
			fmt.Printf("ERROR SIGNAL: %v\n", err)
			return

		case <-time.Tick(10 * time.Second):
			msg, err = shared.Commander.MakePing()
			err = conn.SendMessage(msg)
			if err != nil {
				fmt.Println("PING ERROR?", (*conn).RemoteAddr())
				return
			}
		}
	}
}

func (client *Client) connect() {
	client.wg.Add(1)
	defer client.wg.Done()
	for {
		conn, err := net.Dial("tcp", (*client).Getaway)
		if err != nil {
			fmt.Println("CONNECTION ERROR", err)
			if client.retry > 0 {
				client.retry--
				fmt.Printf("  retry times least %d\n", (*client).retry)
				time.Sleep(ReconnectTimeout * time.Second)
				continue
			} else {
				client.alive = false
				return
			}
		} else {
			(*client).retry = ReconnectTimes
		}
		defer conn.Close()
		client.conn = &shared.Connection{Conn: conn}
		err = client.handle()
		if !client.alive {
			fmt.Printf("??? %v", err)
			return
		}
	}
}
