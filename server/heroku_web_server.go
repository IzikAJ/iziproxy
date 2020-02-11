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
	commonWebResponses
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
	router.PathPrefix("/").HandlerFunc(server.subdomainHandler())

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(server.core.Port), router))
}

func (server *HerokuWEBServer) subdomainHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		req, _ := shared.RequestFromRequest(r)

		signal := make(CodeSignal)
		server.core.place(&ProxyPack{
			Request: req,
			signal:  signal,
		})
		server.core.Stats.start()
		server.core.spaceSignal <- req.ID

		select {
		case <-signal:

			if d, ok := server.core.pool[req.ID]; ok {
				resp := d.Response

				if resp.Status == 0 {
					server.core.Stats.fail()
					server.writeFailResponse(&w, http.StatusBadGateway, "EMPTY RESPONSE FROM CLIENT")
					return
				}
				fmt.Printf("> [%d] %s\n", resp.Status, (*d).Request.Path)

				for _, header := range resp.Headers {
					for _, value := range header.Value {
						w.Header().Set(header.Name, value)
					}
				}

				w.WriteHeader(resp.Status)
				w.Write(resp.Body)
				delete(server.core.pool, req.ID)
				server.core.Stats.complete()
			} else {
				server.core.Stats.fail()
				server.writeFailResponse(&w, http.StatusBadGateway, "NO RESPONSE FROM CLIENT")
			}

		case <-time.After(server.packetTimeout):
			server.core.Stats.timeout()
			server.writeFailResponse(&w, http.StatusGatewayTimeout, "TIMEOUT ERROR")
		}
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
