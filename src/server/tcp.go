package server

import (
	"fmt"
	"net"
	"time"

	"shared"

	"github.com/google/uuid"
)

func handleServerConnection(conf *Config, conn *shared.Connection) {
	defer func() {
		(*conn).Close()
		fmt.Println("CLOSED CONNECTION")
		(*conf).Stats.disconnected()
	}()

	conn.Init()

	spaceSignal := make(chan uuid.UUID)
	earthSignal := make(chan uuid.UUID)
	go func() {
		for {
			msg, err := shared.MsgManager.ReciveMessage(conn)
			if err != nil {
				conf.Stats.fail()
				fmt.Println("reciveMessage ERROR", err, msg)
				return
			}
			resp, err := shared.MsgManager.GetRequest(msg)
			if err != nil {
				fmt.Println("getRequest ERROR", err, msg.Data)
				return
			}

			if req, ok := (*conf).pool[resp.ID]; ok {
				(*req).Response = resp

				earthSignal <- resp.ID
				(*req).signal <- resp.Status
			} else {
				fmt.Println("POOL ERROR")
			}
		}
	}()
	(*conf).space["test"] = spaceSignal
	for {
		select {
		case <-earthSignal:
			// fmt.Println("EARTH SIGNAL?", ruuid)

		case ruuid := <-spaceSignal:
			// recived request uuid
			// fmt.Println("SPACE SIGNAL?")

			if req, ok := (*conf).pool[ruuid]; ok {
				msg, err := shared.Commander.MakeRequest((*req).Request)
				err = conn.SendMessage(msg)
				if err != nil {
					fmt.Println("SEND REQUEST ERROR?", err)
					continue
				}
			} else {
				fmt.Println("NO RECORD IN POOL!", ruuid)
			}
		case <-time.Tick(60 * time.Second):
			msg, err := shared.Commander.MakePing()
			err = conn.SendMessage(msg)
			if err != nil {
				fmt.Println("PING ERROR?", (*conn).RemoteAddr())
				return
			}
		}
	}
}

// TCPServer - run tcp server
func TCPServer(conf *Config) {
	defer (*conf).locker.Done()

	listener, err := net.Listen("tcp", ":2010")
	if err != nil {
		fmt.Println("CANT LISTEN", err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("CANT ACCEPT", err)
			continue
		}

		(*conf).Stats.connected()
		fmt.Println("CONNECTION ACCEPTED", conn.RemoteAddr())
		go handleServerConnection(conf, &shared.Connection{Conn: conn})
	}
}
