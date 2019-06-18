package client

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"time"

	"shared"
)

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

func (client *Client) handle() {
	conn := client.conn
	defer conn.Close()
	conn.Init()
	for {
		msg, err := shared.MsgManager.ReciveMessage(conn)
		if err != nil {
			if err == io.EOF {
				fmt.Println("DISCONNECTED!")
				break
			}
			fmt.Println("MESSAGE ERROR", err, msg)
			continue
		}

		switch msg.Command {
		case shared.CommandPing:
			fmt.Printf(".")
		case shared.CommandRequest:
			go func() {
				req, err := shared.MsgManager.GetRequest(msg)
				fmt.Printf("REQ < %s:%s\n", req.Method, req.Path)

				resp, err := client.load(req)

				msg2, err := shared.Commander.MakeResponse(resp)
				err = shared.MsgManager.SendMessage(msg2, conn)
				if err != nil {
					fmt.Println("ERROR", err)
				}
				time.Sleep((500 + time.Duration(rand.Intn(1000))) * time.Millisecond)
			}()
		default:
			fmt.Println("UNKNOWN COMMAND", msg.Command)
			msg.Print()
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
			if (*client).retry > 0 {
				(*client).retry--
				fmt.Printf("  retry times least %d\n", (*client).retry)
				time.Sleep(time.Second)
				continue
			} else {
				return
			}
		}
		defer conn.Close()
		client.conn = &shared.Connection{Conn: conn}
		client.handle()
	}
}
