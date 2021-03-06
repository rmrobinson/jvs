package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/rmrobinson/jvs/service/building"
	"github.com/rmrobinson/jvs/service/building/pb"
	"google.golang.org/grpc"
)

func main() {
	var (
		port = flag.Int("port", 1338, "Port to listen on")
		proxyAddr  = flag.String("proxy", "", "Address to proxy requests to")
	)
	flag.Parse()

	// Setup the proxy connection first
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(*proxyAddr, opts...)

	if err != nil {
		log.Printf("Error initializing proxy connection to %s: %s\n", *proxyAddr, err.Error())
		os.Exit(1)
	}

	log.Printf("Proxying to %s\n", *proxyAddr)

	// Setup the hub and proxy once we have a connected remote.
	hub := building.NewHub()

	p := building.NewProxyBridge(hub, conn)
	go p.Run()

	connStr := fmt.Sprintf("%s:%d", "", *port)
	lis, err := net.Listen("tcp", connStr)
	if err != nil {
		log.Printf("Error initializing listener: %s\n", err.Error())
		os.Exit(1)
	}
	defer lis.Close()
	log.Printf("Listening on %s\n", connStr)

	api := building.NewAPI(hub)

	grpcServer := grpc.NewServer()
	pb.RegisterBridgeManagerServer(grpcServer, api)
	pb.RegisterDeviceManagerServer(grpcServer, api)
	grpcServer.Serve(lis)
}
