package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"server"
	"shared"
)

func main() {
	fmt.Println("Starting...", flag.Args())
	defer fmt.Println("THE END!")

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	flags := shared.AppFlags{
		Host: "0.0.0.0",
		Port: port,
	}

	flag.StringVar(&(flags.Host), "host", "0.0.0.0", "run as host")
	flag.IntVar(&(flags.Port), "port", port, "run as port")
	flag.StringVar(&(flags.Addr), "addr", "http://localhost:3001", "proxy addr")

	flag.Parse()

	fmt.Println("host", flags.Host)
	fmt.Println("port", flags.Port)

	server.Server(server.Conf.Initialize())
}
