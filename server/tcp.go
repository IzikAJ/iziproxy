package server

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/izikaj/iziproxy/shared"
	"github.com/izikaj/iziproxy/shared/names"
)

func resolveConnectionSpace(conf *Config, data shared.ConnectionSetup, cable *Cable) (err error) {
	cable.Scope = data.Scope
	if _, ok := conf.space[cable.Scope]; ok || cable.Scope == "" {
		// scope already owned / not passed
		if data.Fallback {
			gen := names.ShortNameGenerator(func(name string) bool {
				_, ok := conf.space[name]
				return !ok
			})
			if cable.Scope, err = gen.Next(); err != nil {
				return
			}
		} else {
			return &names.GenerationError{S: "no fallback, sorry"}
		}
	}
	conf.space[cable.Scope] = cable.spaceSignal
	return
}

func handleIncomingMessages(conf *Config, conn *shared.Connection, cable *Cable) (err error) {
	// defer func() {
	// 	client.signal <- err
	// }()

	var msg shared.Message
	for {
		msg, err = shared.MessageManager.ReciveMessage(conn)
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
			var data shared.ConnectionSetup
			data, err = shared.ConnectionSetupFromDump(msg.Data)
			if err != nil {
				fmt.Println("getData ERROR", err, msg.Data)
				conf.Stats.disconnected()
				cable.ufoSignal <- 1
				return
			}
			//

			fmt.Printf("Connection resolving...: %v\n", conf.space)

			err = resolveConnectionSpace(conf, data, cable)
			if err != nil {
				fmt.Println("ConnectionSpace ERROR?", (*conn).RemoteAddr())

				msg, _ = shared.Commander.MakeFailed(shared.ConnectionError{
					Code:    "namespace_resolve_error",
					Message: err.Error(),
				})
				conn.SendMessage(msg)
				return
			}
			fmt.Printf("Connection resolved: %v\n", conf.space)

			cable.Authorized = true
			cable.Owner = "Tester"

			msg, err = shared.Commander.MakeReady(shared.ConnectionResult{
				Scope:   cable.Scope,
				Status:  "connected",
				Message: "connected successfully",
			})
			err = conn.SendMessage(msg)
			if err != nil {
				fmt.Println("Ready ERROR?", (*conn).RemoteAddr())
				return
			}

		case shared.CommandResponse:
			var resp shared.Request
			resp, err = shared.MessageManager.GetRequest(msg)
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

		case shared.CommandPing:
			msg, err = shared.Commander.MakePong()
			err = conn.SendMessage(msg)
			if err != nil {
				fmt.Println("PONG ERROR?", (*conn).RemoteAddr())
				return
			}

			fmt.Print(">")

		default:
			fmt.Println("RECIVED UNHANDLED MESSAGE")
			msg.Print()
		}
	}
}

func handleServerConnection(conf *Config, conn *shared.Connection) {
	cable := &Cable{
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

	go handleIncomingMessages(conf, conn, cable)

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
		case <-time.Tick(10 * time.Second):
			msg, err := shared.Commander.MakePing()
			attempts := 10
			for {
				err = conn.SendMessage(msg)
				if err == nil {
					break
				} else {
					attempts--
					if attempts < 0 {
						fmt.Println("PING ERROR!!!", (*conn).RemoteAddr())
						return
					}
					fmt.Println("PING ERROR?", (*conn).RemoteAddr())
				}
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
	// (*conf).locker.Add(1)
	defer (*conf).locker.Done()

	listener, err := net.Listen("tcp", ":2010")
	if err != nil {
		fmt.Println("CAN'T LISTEN", err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("CAN'T ACCEPT", err)
			continue
		}

		(*conf).Stats.connected()
		fmt.Println("CONNECTION ACCEPTED", conn.RemoteAddr())
		go handleServerConnection(conf, &shared.Connection{Conn: conn})
	}
}
