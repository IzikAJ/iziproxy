package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/izikaj/iziproxy/shared"
)

// HerokuWEBServer - server instance
type HerokuWEBServer struct {
	core     *Server
	hostName string

	packetTimeout time.Duration
	commonWebHelpers
}

// Start - start HerokuWEBserver daemon
func (server *HerokuWEBServer) Start() {
	fmt.Println("Starting HerokuWEBServer...")
	defer fmt.Println("HerokuWEBServer stopped")
	defer server.core.locker.Done()

	server.listen()
}

func (server *HerokuWEBServer) listen() {
	router := mux.NewRouter()

	router.Path("/__stats").Methods("GET").HandlerFunc(server.statsHandler(server.core))
	router.PathPrefix("/").HandlerFunc(server.clientRequestHandler())

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(server.core.Port), router))
}

func (server *HerokuWEBServer) clientRequestHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		req, _ := shared.RequestFromRequest(r)

		signal := make(CodeSignal)
		server.core.place(&ProxyPack{
			Request: req,
			signal:  signal,
		})
		server.core.Stats.start()

		server.core.spaceSignal <- req.ID

		server.waitForResponse(waitForResponseParams{
			core:    server.core,
			req:     &req,
			signal:  &signal,
			w:       &w,
			timeout: server.packetTimeout,
		})
	}
}

// NewHerokuWEBServer - create new HerokuWEBServer with confguration
func NewHerokuWEBServer(core *Server) *HerokuWEBServer {
	return &HerokuWEBServer{
		core:     core,
		hostName: "proxy.me",

		packetTimeout: 120 * time.Second,
	}
}
