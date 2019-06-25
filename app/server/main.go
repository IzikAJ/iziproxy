package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/izikaj/iziproxy/server"
	"github.com/izikaj/iziproxy/shared"
)

func main() {
	fmt.Println("Starting...")
	defer fmt.Println("THE END!")

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	if port <= 0 {
		port = 3000
	}
	flags := shared.AppFlags{
		Host:   "0.0.0.0",
		Port:   port,
		Single: false,
	}

	flag.StringVar(&(flags.Host), "host", "0.0.0.0", "run with host")
	flag.IntVar(&(flags.Port), "port", port, "run with port")
	flag.BoolVar(&(flags.Single), "single", false, "run as single proxy")

	flag.Parse()

	fmt.Println("host", flags.Host)
	fmt.Println("port", flags.Port)
	if flags.Single {
		fmt.Println("RUN IN SINGLE INSTANCE MODE")
	}

	conf := server.Config{
		Host:   flags.Host,
		Port:   flags.Port,
		Single: flags.Single,
	}

	server.Server(conf.Initialize())
}
