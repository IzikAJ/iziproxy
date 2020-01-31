package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/izikaj/iziproxy/server"
)

func main() {
	fmt.Println("Starting...")
	defer fmt.Println("THE END!")

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
	flag.BoolVar(&(conf.Single), "single", conf.Single, "run as single proxy")

	flag.Parse()

	fmt.Println("host", conf.Host)
	fmt.Println("port", conf.Port)
	if conf.Single {
		fmt.Println("RUN IN SINGLE INSTANCE MODE")
	}

	server.NewServer(conf).Start()
}
