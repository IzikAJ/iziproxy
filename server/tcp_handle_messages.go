package server

import (
	"fmt"
	"io"
	"net"

	"github.com/izikaj/iziproxy/shared"
)

type defaultTCPCommands struct {
}

func (server *defaultTCPCommands) onPing(conn *shared.Connection, cable *Cable) (err error) {
	var data shared.Message
	data, err = shared.Commander.MakePong()
	err = conn.SendMessage(data)
	if err != nil {
		fmt.Println("PONG ERROR?", conn.RemoteAddr())
		return
	}
	fmt.Print(">")
	return
}

func (server *defaultTCPCommands) onPong(conn *shared.Connection, cable *Cable) (err error) {
	fmt.Print("<")
	return
}

func (server *defaultTCPCommands) onUnrecognized(conn *shared.Connection, cable *Cable, data shared.Message) (err error) {
	fmt.Println("RECIVED UNHANDLED MESSAGE")
	data.Print()
	return
}

func (server *defaultTCPCommands) onResponse(conn *shared.Connection, cable *Cable, data shared.Request) (err error) {
	panic("onResponse - should be implemented")
}

func (server *defaultTCPCommands) onSetup(conn *shared.Connection, cable *Cable, data shared.ConnectionSetup) (err error) {
	panic("onSetup - should be implemented")
}

func handleTCPMessages(server AbstractTCPCommands, core *Server, conn *shared.Connection, cable *Cable) (err error) {
	defer func() {
		cable.ufoSignal <- err
	}()

	var msg shared.Message
	for {
		msg, err = shared.MessageManager.ReciveMessage(conn)
		if err != nil {
			if err == io.EOF {
				fmt.Println("R: CLIENT DISCONNECTED (EOF)", err)
				return
			}
			switch err.(type) {
			case net.Error:
				fmt.Println("R: CLIENT DISCONNECTED (EPIPE)", err)
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
				return
			}
			server.onSetup(conn, cable, data)
			//

			// fmt.Printf("Connection resolving...: %v\n", core.space)

			// // err = server.resolveConnectionSpace(data, cable)
			// if err != nil {
			// 	fmt.Println("ConnectionSpace ERROR?", conn.RemoteAddr())

			// 	msg, _ = shared.Commander.MakeFailed(shared.ConnectionError{
			// 		Code:    "namespace_resolve_error",
			// 		Message: err.Error(),
			// 	})
			// 	conn.SendMessage(msg)
			// 	return
			// }
			// fmt.Printf("Connection resolved: %v\n", core.space)

			// cable.Authorized = true
			// cable.Owner = "Tester"

			// msg, err = shared.Commander.MakeReady(shared.ConnectionResult{
			// 	Scope:   cable.Scope,
			// 	Status:  "connected",
			// 	Message: "connected successfully",
			// })
			// err = conn.SendMessage(msg)
			// if err != nil {
			// 	fmt.Println("Ready ERROR?", conn.RemoteAddr())
			// 	return
			// }

		case shared.CommandResponse:
			var resp shared.Request
			resp, err = shared.MessageManager.GetRequest(msg)
			if err != nil {
				fmt.Println("getRequest ERROR", err, msg.Data)
				return
			}
			server.onResponse(conn, cable, resp)

		case shared.CommandPong:
			server.onPong(conn, cable)

		case shared.CommandPing:
			server.onPing(conn, cable)

		default:
			server.onUnrecognized(conn, cable, msg)
			fmt.Println("RECIVED UNHANDLED MESSAGE")
			msg.Print()
		}
	}
}
