package server

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/izikaj/iziproxy/shared"
)

func (server *TCPServer) handleSignals(conn *shared.Connection, cable *Cable) {
	for {
		select {
		case ruuid := <-cable.spaceSignal:
			if req, ok := server.core.pool[ruuid]; ok {
				msg, err := shared.Commander.MakeRequest(req.Request)
				err = conn.SendMessage(msg)
				if err != nil {
					if err == io.EOF {
						fmt.Println("W: CLIENT DISCONNECTED (EOF)", err)
						server.core.Stats.disconnected()
						return
					}
					switch err.(type) {
					case net.Error:
						fmt.Println("W: CLIENT DISCONNECTED (EPIPE)", err)
						server.core.Stats.disconnected()
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
						fmt.Println("PING ERROR!!!", conn.RemoteAddr())
						return
					}
					fmt.Println("PING ERROR?", conn.RemoteAddr())
				}
			}
		case <-cable.ufoSignal:
			// TODO
			// handle ufo signals
			return
		}
	}
}
