package main

import (
	"fmt"
	"net"

	"faltung.ca/jvs/lib/devicemanager-go"
)

type server struct {
	manager *devicemanager.Manager
}

func newServer() *server {
	s := &server{
		manager: devicemanager.NewManager(),
	}

	return s
}

func (s *server) run(devicePath string, configPath string, port int) error {
	var b Bridge

	err := b.Setup(devicePath, configPath)
	defer b.Shutdown()

	if err != nil {
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	defer lis.Close()

	if err != nil {
		return err
	}

	var tmp devicemanager.Bridge
	tmp = &b

	s.manager.AddBridge(tmp)

	s.manager.Run(lis)

	return nil
}
