package server

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/izikaj/iziproxy/shared"
)

func handleServerConnection(conf *Config, conn *shared.Connection) {
	cable := Cable{
		Connected: true,

		pool:        make(map[uuid.UUID]*ProxyPack),
		spaceSignal: make(chan uuid.UUID),
		ufoSignal:   make(chan int),
	}

	defer func() {
		(*conn).Close()
		fmt.Println("CLOSED CONNECTION")
		if cable.Scope != "" {
			delete((*conf).space, cable.Scope)
		}

		(*conf).Stats.disconnected()
	}()

	conn.Init()

	go func() {
		for {
			msg, err := shared.MsgManager.ReciveMessage(conn)
			if err != nil {
				conf.Stats.fail()
				if err == io.EOF {
					fmt.Println("R: CLIENT DISCONNECTED (EOF)", err)
					conf.Stats.disconnected()
					cable.ufoSignal <- 1
					return
				}
				switch err.(type) {
				case net.Error:
					fmt.Println("R: CLIENT DISCONNECTED (EPIPE)", err)
					conf.Stats.disconnected()
					cable.ufoSignal <- 1
					return
				}
				fmt.Println("reciveMessage ERROR", err, msg)
				return
			}
			switch msg.Command {
			case shared.CommandSetup:
				fmt.Println("CONNECTION SETUP COMMAND")
				msg.Print()
				//
				data, err := shared.ConnectionSetupFromDump(msg.Data)
				if err != nil {
					fmt.Println("getData ERROR", err, msg.Data)
					conf.Stats.disconnected()
					cable.ufoSignal <- 1
					return
				}
				//
				cable.Scope = data.Scope
				(*conf).space[data.Scope] = cable.spaceSignal
				cable.Authorized = true
				cable.Owner = "Tester"

				msg, err := shared.Commander.MakeReady()
				err = conn.SendMessage(msg)
				if err != nil {
					fmt.Println("Ready ERROR?", (*conn).RemoteAddr())
					return
				}

			case shared.CommandResponse:
				resp, err := shared.MsgManager.GetRequest(msg)
				if err != nil {
					fmt.Println("getRequest ERROR", err, msg.Data)
					return
				}
				if req, ok := (*conf).pool[resp.ID]; ok {
					(*req).Response = resp

					(*req).signal <- resp.Status
				} else {
					fmt.Println("POOL ERROR")
				}

			case shared.CommandPong:
				fmt.Print("<")

			default:
				fmt.Println("RECIVED UNHANDLED MESSAGE")
				msg.Print()
			}
		}
	}()
	for {
		select {
		case ruuid := <-cable.spaceSignal:
			if req, ok := (*conf).pool[ruuid]; ok {
				msg, err := shared.Commander.MakeRequest((*req).Request)
				err = conn.SendMessage(msg)
				if err != nil {
					if err == io.EOF {
						fmt.Println("W: CLIENT DISCONNECTED (EOF)", err)
						conf.Stats.disconnected()
						return
					}
					switch err.(type) {
					case net.Error:
						fmt.Println("W: CLIENT DISCONNECTED (EPIPE)", err)
						conf.Stats.disconnected()
						return
					}
					fmt.Println("SEND REQUEST ERROR", err)
					continue
				}
			}
		case <-time.Tick(30 * time.Second):
			msg, err := shared.Commander.MakePing()
			err = conn.SendMessage(msg)
			if err != nil {
				fmt.Println("PING ERROR?", (*conn).RemoteAddr())
				return
			}
		case <-cable.ufoSignal:
			// TODO
			// handle ufo signals
			return
		}
	}
}

// TCPServer - run tcp server
func TCPServer(conf *Config) {
	(*conf).locker.Add(1)
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
