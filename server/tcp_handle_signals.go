package server

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/izikaj/iziproxy/shared"
)

func onTCPRequest(core *Server, conn *shared.Connection, ruuid uuid.UUID) error {
	var pack *ProxyPack
	var ok bool
	core.Lock()
	pack, ok = core.pool[ruuid]
	core.Unlock()
	if !ok {
		return nil
	}

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
			onTCPRequest(core, conn, ruuid)

		case ruuid := <-core.spaceSignal:
			onTCPRequest(core, conn, ruuid)

		case <-time.Tick(60 * time.Second):
			msg, err := shared.Commander.MakePing()
			err = conn.SendMessage(msg)
			if err == nil {
				fmt.Println("PING ERROR!!!", conn.RemoteAddr())
				return
			}

		case err := <-cable.ufoSignal:
			fmt.Printf("ufoSignal [%q], %v\n", err, err.Error())
			return
		}
	}
}
