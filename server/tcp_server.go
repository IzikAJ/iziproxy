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

	go server.handleMessages(conn, cable)
	server.handleSignals(conn, cable)
}

func (server *TCPServer) resolveConnectionSpace(data shared.ConnectionSetup, cable *Cable) (err error) {
	if server.core.Single {
		fmt.Println("spaceSignal 3", cable.spaceSignal)
		server.core.globalSpaceSignal = cable.spaceSignal
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

// NewTCPServer - create new TCPServer with confguration
func NewTCPServer(core *Server) *TCPServer {
	return &TCPServer{
		core: core,
	}
}
