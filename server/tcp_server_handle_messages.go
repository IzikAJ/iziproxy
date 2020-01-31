package server

import (
	"fmt"
	"io"
	"net"

	"github.com/izikaj/iziproxy/shared"
)

func (server *TCPServer) handleMessages(conn *shared.Connection, cable *Cable) (err error) {
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
			//

			fmt.Printf("Connection resolving...: %v\n", server.core.space)

			err = server.resolveConnectionSpace(data, cable)
			if err != nil {
				fmt.Println("ConnectionSpace ERROR?", conn.RemoteAddr())

				msg, _ = shared.Commander.MakeFailed(shared.ConnectionError{
					Code:    "namespace_resolve_error",
					Message: err.Error(),
				})
				conn.SendMessage(msg)
				return
			}
			fmt.Printf("Connection resolved: %v\n", server.core.space)

			cable.Authorized = true
			cable.Owner = "Tester"

			msg, err = shared.Commander.MakeReady(shared.ConnectionResult{
				Scope:   cable.Scope,
				Status:  "connected",
				Message: "connected successfully",
			})
			err = conn.SendMessage(msg)
			if err != nil {
				fmt.Println("Ready ERROR?", conn.RemoteAddr())
				return
			}

		case shared.CommandResponse:
			var resp shared.Request
			resp, err = shared.MessageManager.GetRequest(msg)
			if err != nil {
				fmt.Println("getRequest ERROR", err, msg.Data)
				return
			}
			if req, ok := server.core.pool[resp.ID]; ok {
				req.Response = resp

				req.signal <- resp.Status
			} else {
				fmt.Println("POOL ERROR")
			}

		case shared.CommandPong:
			fmt.Print("<")

		case shared.CommandPing:
			msg, err = shared.Commander.MakePong()
			err = conn.SendMessage(msg)
			if err != nil {
				fmt.Println("PONG ERROR?", conn.RemoteAddr())
				return
			}

			fmt.Print(">")

		default:
			fmt.Println("RECIVED UNHANDLED MESSAGE")
			msg.Print()
		}
	}
}
