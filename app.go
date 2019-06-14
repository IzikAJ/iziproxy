package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"shared"
)

func statsHandler func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, time.Now().String())
	// raw, _ := json.Marshal((*conf).Stats)
	// fmt.Fprintln(w, string(raw))
}

func serve(port: Int) {
	router := mux.NewRouter()

	router.Host(
		"{subdomain:.+}." + hostName,
	).HandlerFunc(
		bindSubdomainHandler(conf),
	)

	router.HandleFunc("/stats", bindStatsHandler(conf))

	router.Methods("GET").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/stats", 302)
		},
	)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), router))
}

func main() {
	fmt.Println("Starting...")
	defer fmt.Println("THE END!")

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	flags := shared.AppFlags{
		Host: "0.0.0.0",
		Port: port,
	}

	flag.StringVar(&(flags.Host), "host", "0.0.0.0", "run as host")
	flag.IntVar(&(flags.Port), "port", port, "run as port")

	flag.Parse()

	fmt.Println("host", flags.Host)
	fmt.Println("port", flags.Port)

	serve(port)
}
