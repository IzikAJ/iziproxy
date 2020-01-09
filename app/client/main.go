package main

import (
	"flag"
	"fmt"

	"github.com/izikaj/iziproxy/client"
	"github.com/izikaj/iziproxy/shared"
)

func main() {
	fmt.Println("Starting client...", flag.Args())
	defer fmt.Println("THE END!")

	flags := shared.AppFlags{
		Addr:  "http://localhost:3001",
		Space: "",
	}

	// flag.StringVar(&(flags.Host), "host", "0.0.0.0", "run as host")
	// flag.IntVar(&(flags.Port), "port", 3000, "run as port")
	flag.StringVar(&(flags.Addr), "addr", "http://localhost:3001", "proxy addr")
	flag.StringVar(&(flags.Space), "space", "", "proxy space")
	// flag.StringVar(&(flags.Host), "host", "0.0.0.0", "run as host")

	flag.Parse()

	fmt.Println("flags", flags)
	client := &client.Client{
		Getaway: "127.0.0.1:2010",
		Host:    flags.Addr,
		Space:   flags.Space,
	}
	client.Init()
	client.Start()
}
