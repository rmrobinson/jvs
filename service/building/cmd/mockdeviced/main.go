package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/rmrobinson/jvs/service/building"
	"github.com/rmrobinson/jvs/service/building/mock"
	"github.com/rmrobinson/jvs/service/building/pb"
	"google.golang.org/grpc"
)

func main() {
	var (
		port = flag.Int("port", 1337, "The port for the deviced process to listen on")
	)
	flag.Parse()

	bm := building.NewHub()

	msb := mock.NewSyncBridge()
	bm.AddBridge(msb, 5*time.Second)
	go msb.Run()

	mab := mock.NewAsyncBridge()
	bm.AddAsyncBridge(mab)
	go mab.Run()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Printf("Error initializing listener: %s\n", err.Error())
		os.Exit(1)
	}
	defer lis.Close()

	api := building.NewAPI(bm)

	grpcServer := grpc.NewServer()
	pb.RegisterBridgeManagerServer(grpcServer, api)
	pb.RegisterDeviceManagerServer(grpcServer, api)

	log.Printf("Listening on :%d\n", *port)
	grpcServer.Serve(lis)
}
