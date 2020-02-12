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

type spaceParams struct {
	subdomain string
}

// WEBServer - server instance
type WEBServer struct {
	core     *Server
	hostName string

	packetTimeout time.Duration
	commonWebHelpers
}

// Start - start WEBserver daemon
func (server *WEBServer) Start() {
	fmt.Println("Starting WEBServer...")
	defer fmt.Println("WEBServer stopped")
	defer server.core.locker.Done()

	server.listen()
}

func (server *WEBServer) listen() {
	router := mux.NewRouter()

	router.Host(
		fmt.Sprintf("{subdomain:.+}.%v", server.hostName),
	).HandlerFunc(server.clientRequestHandler())

	router.HandleFunc("/__stats", server.statsHandler(server.core))

	router.Methods("GET").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "//"+server.hostName+"/__stats", 302)
		},
	)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(server.core.Port), router))
}

func (server *WEBServer) clientRequestHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		req, _ := shared.RequestFromRequest(r)

		vars := mux.Vars(r)
		params := spaceParams{
			subdomain: vars["subdomain"],
		}

		signal := make(CodeSignal)
		server.core.place(&ProxyPack{
			Request: req,
			signal:  signal,
		})
		server.core.Stats.start()

		if spaceSignal, err := server.core.findSpaceSignal(params); err == nil {
			fmt.Println("spaceSignal 1", spaceSignal)
			spaceSignal <- req.ID
		} else {
			writeFailResponse(&w, http.StatusBadGateway, "NO CLIENT CONNECTED")
			return
		}

		server.waitForResponse(waitForResponseParams{
			core:    server.core,
			req:     &req,
			signal:  &signal,
			w:       &w,
			timeout: server.packetTimeout,
		})
	}
}

// NewWEBServer - create new WEBServer with confguration
func NewWEBServer(core *Server) *WEBServer {
	return &WEBServer{
		core:     core,
		hostName: "proxy.me",

		packetTimeout: 120 * time.Second,
	}
}
