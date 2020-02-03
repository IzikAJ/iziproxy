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
}

// Start - start WEBserver daemon
func (server *WEBServer) Start() {
	fmt.Println("Starting WEBServer...")
	defer fmt.Println("WEBServer stopped")
	defer server.core.locker.Done()

	if server.core.Single {
		server.serveSingle()
	} else {
		server.serveSpaced()
	}
}

func (server *WEBServer) serveSpaced() {
	router := mux.NewRouter()

	router.Host(
		"{subdomain:.+}." + server.hostName,
	).HandlerFunc(server.subdomainHandler())

	router.HandleFunc("/stats", server.statsHandler())

	router.Methods("GET").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "//"+server.hostName+"/stats", 302)
		},
	)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(server.core.Port), router))
}

func (server *WEBServer) serveSingle() {
	router := mux.NewRouter()

	router.Path("/__stats").Methods("GET").HandlerFunc(server.statsHandler())
	router.PathPrefix("/").HandlerFunc(server.subdomainHandler())

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(server.core.Port), router))
}

func (server *WEBServer) subdomainHandler() func(http.ResponseWriter, *http.Request) {
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

		select {
		case <-signal:

			if d, ok := server.core.pool[req.ID]; ok {
				resp := d.Response

				if resp.Status == 0 {
					server.core.Stats.fail()
					writeFailResponse(&w, http.StatusBadGateway, "EMPTY RESPONSE FROM CLIENT")
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
				writeFailResponse(&w, http.StatusBadGateway, "NO RESPONSE FROM CLIENT")
			}

		case <-time.Tick(server.packetTimeout):
			server.core.Stats.timeout()
			writeFailResponse(&w, http.StatusGatewayTimeout, "TIMEOUT ERROR")
		}
	}
}

func (server *WEBServer) statsHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		writeStatsResponse(&w, server.core)
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
