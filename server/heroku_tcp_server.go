package server

import (
	"fmt"
	"net"

	"github.com/izikaj/iziproxy/shared"
	"github.com/izikaj/iziproxy/shared/names"
)

// HerokuTCPServer - heroku server instance
type HerokuTCPServer struct {
	core *Server

	// include default command handlers
	defaultTCPCommands
}

// Start - start HerokuTCPServer daemon
func (server *HerokuTCPServer) Start() {
	fmt.Println("Starting HerokuTCPServer...")
	defer fmt.Println("HerokuTCPServer stopped")
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

func (server *HerokuTCPServer) handleServerConnection(conn *shared.Connection) {
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

func (server *HerokuTCPServer) resolveConnectionSpace(data shared.ConnectionSetup, cable *Cable) (err error) {
	if server.core.Single {
		fmt.Println("spaceSignal 3", cable.spaceSignal)
		server.core.spaceSignal = cable.spaceSignal
		return nil
	}
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

func (server *HerokuTCPServer) onSetup(conn *shared.Connection, cable *Cable, data shared.ConnectionSetup) (err error) {
	fmt.Println("spaceSignal 3", cable.spaceSignal)
	server.core.spaceSignal = cable.spaceSignal
	return nil
}

func (server *HerokuTCPServer) onResponse(conn *shared.Connection, cable *Cable, data shared.Request) (err error) {
	if req, ok := server.core.pool[data.ID]; ok {
		req.Response = data

		req.signal <- data.Status
	} else {
		fmt.Println("POOL ERROR")
	}
	return nil
}

// NewHerokuTCPServer - create new HerokuTCPServer with confguration
func NewHerokuTCPServer(core *Server) *HerokuTCPServer {
	return &HerokuTCPServer{
		core: core,
	}
}
