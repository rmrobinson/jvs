package main

import (
	"flag"
	"log"
)

var (
	port       = flag.Int("port", 10000, "The listening RPC port")
	configFile = flag.String("config", "config.db", "The path to the configuration database")
)

func main() {
	flag.Parse()

	s := newServer()

	log.Printf("Listening for gRPC requests on %d\n", *port)

	err := s.run(*configFile, *port)

	if err != nil {
		log.Printf("Error: %s\n", err)
	}
}
