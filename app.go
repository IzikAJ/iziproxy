package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"server"
	"shared"

	"github.com/gorilla/mux"
)

func statsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, time.Now().String())
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "TODO\n\n\n")
	fmt.Fprintln(w, time.Now().String())
}

func serve(port int) {
	router := mux.NewRouter()

	router.HandleFunc("/__stats", statsHandler)
	router.PathPrefix("/").HandlerFunc(proxyHandler)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), router))
}

func main() {
	fmt.Println("Starting...")
	defer fmt.Println("THE END!")

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	if port <= 0 {
		port = 3000
	}
	flags := shared.AppFlags{
		Host: "0.0.0.0",
		Port: port,
	}

	flag.StringVar(&(flags.Host), "host", "0.0.0.0", "run as host")
	flag.IntVar(&(flags.Port), "port", port, "run as port")

	flag.Parse()

	fmt.Println("host", flags.Host)
	fmt.Println("port", flags.Port)

	conf := server.Config{}

	go server.TCPServer(conf.Initialize())

	serve(port)
}
