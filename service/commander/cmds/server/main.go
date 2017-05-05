package main

import (
	"faltung.ca/jvs/service/commander"
	"flag"
	"log"
	"faltung.ca/jvs/service/commander/modules"
)

var (
	port = flag.Int("port", 10004, "The listening RPC port")
)

func main() {
	flag.Parse()

	s := commander.NewServer()

	log.Printf("Listening for gRPC requests on %d\n", *port)

	var c commander.RootCommand
	c.AddSubCommand(&modules.DeviceCommand{})

	s.SetRootCommand(c)
	err := s.Run(*port)

	if err != nil {
		log.Printf("Error: %s\n", err)
	}
}
