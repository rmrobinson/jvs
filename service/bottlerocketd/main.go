package main

import (
	"flag"
	"log"
)

var (
	port       = flag.Int("port", 10001, "The listening RPC port")
	deviceFile = flag.String("device", "/dev/firecracker", "The path of the device to use")
	configFile = flag.String("config", "config.db", "The path of the config file to use")
)

func main() {
	flag.Parse()

	s := newServer()

	log.Printf("Listening for gRPC requests on %d\n", *port)

	err := s.run(*deviceFile, *configFile, *port)

	if err != nil {
		log.Printf("Error: %s\n", err)
	}
}
