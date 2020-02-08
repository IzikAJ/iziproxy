package server

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/izikaj/iziproxy/shared"
)

func (server *TCPServer) onRequest(conn *shared.Connection, pack *ProxyPack) error {
	msg, err := shared.Commander.MakeRequest(pack.Request)
	err = conn.SendMessage(msg)
	if err != nil {
		if err == io.EOF {
			fmt.Println("W: CLIENT DISCONNECTED (EOF)", err)
			return err
		}
		switch err.(type) {
		case net.Error:
			fmt.Println("W: CLIENT DISCONNECTED (EPIPE)", err)
			return err
		}
		fmt.Println("SEND REQUEST ERROR", err)
	}
	return nil
}

func onTCPRequest(conn *shared.Connection, pack *ProxyPack) error {
	msg, err := shared.Commander.MakeRequest(pack.Request)
	err = conn.SendMessage(msg)
	if err != nil {
		if err == io.EOF {
			fmt.Println("W: CLIENT DISCONNECTED (EOF)", err)
			return err
		}
		switch err.(type) {
		case net.Error:
			fmt.Println("W: CLIENT DISCONNECTED (EPIPE)", err)
			return err
		}
		fmt.Println("SEND REQUEST ERROR", err)
	}
	return nil
}

func handleTCPSignals(core *Server, conn *shared.Connection, cable *Cable) {
	for {
		select {
		case ruuid := <-cable.spaceSignal:
			if req, ok := core.pool[ruuid]; ok {
				if err := onTCPRequest(conn, req); err != nil {
					return
				}
			}

		case ruuid := <-core.spaceSignal:
			fmt.Println("spaceSignal 2", core.spaceSignal)
			if req, ok := core.pool[ruuid]; ok {
				if err := onTCPRequest(conn, req); err != nil {
					return
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

func (server *TCPServer) handleSignals(conn *shared.Connection, cable *Cable) {
	for {
		select {
		case ruuid := <-cable.spaceSignal:
			if req, ok := server.core.pool[ruuid]; ok {
				if err := server.onRequest(conn, req); err != nil {
					return
				}
			}

		case ruuid := <-server.core.spaceSignal:
			fmt.Println("spaceSignal 2", server.core.spaceSignal)
			if req, ok := server.core.pool[ruuid]; ok {
				if err := server.onRequest(conn, req); err != nil {
					return
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
