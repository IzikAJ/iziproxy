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
						return
					}
					switch err.(type) {
					case net.Error:
						fmt.Println("W: CLIENT DISCONNECTED (EPIPE)", err)
						return
					}
					fmt.Println("SEND REQUEST ERROR", err)
					continue
				}
			}
		case <-time.Tick(60 * time.Second):
			msg, err := shared.Commander.MakePing()
			err = conn.SendMessage(msg)
			if err == nil {
				fmt.Println("PING ERROR!!!", conn.RemoteAddr())
				return
			}
		case err := <-cable.ufoSignal:
			// TODO
			// handle ufo signals?
			fmt.Printf("ufoSignal [%q], %v\n", err, err.Error())
			return
		}
	}
}
