package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/izikaj/iziproxy/server"
)

func main() {
	fmt.Println("Starting...")
	defer fmt.Println("THE END!")
	runtime.GOMAXPROCS(4)

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	if port <= 0 {
		port = 3000
	}
	conf := &server.Config{
		Host:   "0.0.0.0",
		Port:   port,
		Single: false,
	}

	flag.StringVar(&(conf.Host), "host", conf.Host, "run with host")
	flag.IntVar(&(conf.Port), "port", conf.Port, "run with port")

	flag.Parse()

	fmt.Println("host", conf.Host)
	fmt.Println("port", conf.Port)

	server.NewServer(conf).Start()
}
