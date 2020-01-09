package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/izikaj/iziproxy/shared"
)

const (
	hostName      = "proxy.me"
	packetTimeout = 120 * time.Second
)

// Web - simple web Web
type Web struct {
	port int
	host string
}

func (server Web) start(conf *Config) {
	defer (*conf).locker.Done()

	fmt.Println("TODO: WEB SERVER")

	if (*conf).Single {
		serveSingle(conf)
	} else {
		serve(conf)
	}
}

func placePack(conf *Config, pack *ProxyPack) {
	(*conf).Lock()
	defer (*conf).Unlock()
	(*conf).pool[(*pack).Request.ID] = pack
}

func failResp(w *http.ResponseWriter, status int, msg string) {
	(*w).WriteHeader(status)
	(*w).Write([]byte(msg))
}

func bindSubdomainHandler(conf *Config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		req, _ := shared.RequestFromRequest(r)

		vars := mux.Vars(r)
		subdomain := vars["subdomain"]

		signal := make(chan int)

		placePack(conf, &ProxyPack{
			Request: req,
			signal:  signal,
		})
		(*conf).Stats.start()

		if spaceSignal, ok := (*conf).space[subdomain]; ok {
			spaceSignal <- req.ID
		} else {
			failResp(&w, http.StatusBadGateway, "NO CLIENT CONNECTED")
			return
		}

		select {
		case <-signal:

			if d, ok := (*conf).pool[req.ID]; ok {
				resp := (*d).Response

				if resp.Status == 0 {
					(*conf).Stats.fail()
					failResp(&w, http.StatusBadGateway, "EMPTY RESPONSE FROM CLIENT")
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
				(*conf).Stats.complete()
			} else {
				conf.Stats.fail()
				failResp(&w, http.StatusBadGateway, "NO RESPONSE FROM CLIENT")
			}

		case <-time.Tick(packetTimeout):
			conf.Stats.timeout()
			failResp(&w, http.StatusGatewayTimeout, "TIMEOUT ERROR")
		}
	}
}

func bindStatsHandler(conf *Config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, time.Now().String())
		raw, _ := json.Marshal((*conf).Stats)
		fmt.Fprintln(w, string(raw))
	}
}

func serve(conf *Config) {
	router := mux.NewRouter()

	router.Host(
		"{subdomain:.+}." + hostName,
	).HandlerFunc(
		bindSubdomainHandler(conf),
	)

	router.HandleFunc("/stats", bindStatsHandler(conf))

	router.Methods("GET").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "//"+hostName+"/stats", 302)
		},
	)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(conf.Port), router))
}

func serveSingle(conf *Config) {
	fmt.Println("STARTING IN SINGLE MODE")
	router := mux.NewRouter()

	router.HandleFunc("/stats", bindStatsHandler(conf))

	router.PathPrefix("/").HandlerFunc(
		bindSubdomainHandler(conf),
	)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(conf.Port), router))
}
