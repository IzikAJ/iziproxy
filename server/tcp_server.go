package server

import (
	"fmt"
	"net"

	"github.com/izikaj/iziproxy/shared"
	"github.com/izikaj/iziproxy/shared/names"
)

// TCPServer - server instance
type TCPServer struct {
	core *Server

	// include default command handlers
	defaultTCPCommands
}

// Start - start TCPServer daemon
func (server *TCPServer) Start() {
	fmt.Println("Starting TCPServer...")
	defer fmt.Println("TCPServer stopped")
	defer server.core.locker.Done()

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

		server.core.Stats.connected()
		fmt.Println("CONNECTION ACCEPTED", conn.RemoteAddr())
		go server.handleServerConnection(&shared.Connection{Conn: conn})
	}
}

func (server *TCPServer) handleServerConnection(conn *shared.Connection) {
	cable := &Cable{
		Connected: true,

		spaceSignal: make(SpaceSignal),
		ufoSignal:   make(UfoSignal),
	}

	defer func() {
		conn.Close()
		fmt.Println("CLOSED CONNECTION")
		if cable.Scope != "" {
			delete(server.core.space, cable.Scope)
		}

		server.core.Stats.disconnected()
	}()

	conn.Init()

	go handleTCPMessages(server, server.core, conn, cable)
	handleTCPSignals(server.core, conn, cable)
}

func (server *TCPServer) resolveConnectionSpace(data shared.ConnectionSetup, cable *Cable) (err error) {
	cable.Scope = data.Scope
	if _, ok := server.core.space[cable.Scope]; ok || cable.Scope == "" {
		// scope already owned / not passed
		if data.Fallback {
			gen := names.ShortNameGenerator(func(name string) bool {
				_, ok := server.core.space[name]
				return !ok
			})
			if cable.Scope, err = gen.Next(); err != nil {
				return
			}
		} else {
			return names.NewGenerationError("no fallback, sorry")
		}
	}
	server.core.space[cable.Scope] = cable.spaceSignal
	return
}

func (server *TCPServer) onSetup(conn *shared.Connection, cable *Cable, data shared.ConnectionSetup) (err error) {
	fmt.Printf("Connection resolving...: %v\n", server.core.space)
	var msg shared.Message

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
	}
	return
}

func (server *TCPServer) onResponse(conn *shared.Connection, cable *Cable, data shared.Request) (err error) {
	server.core.Lock()
	req, ok := server.core.pool[data.ID]
	server.core.Unlock()
	if ok {
		req.Response = data
		req.signal <- data.Status
	} else {
		fmt.Println("POOL ERROR")
	}
	return nil
}

// NewTCPServer - create new TCPServer with confguration
func NewTCPServer(core *Server) *TCPServer {
	return &TCPServer{
		core: core,
	}
}
