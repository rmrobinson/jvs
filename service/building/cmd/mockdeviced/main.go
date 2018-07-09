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
		port = flag.Int("port", 1337, "The port for the mockdeviced process to listen on")
		dbPath = flag.String("dbPath", "/tmp/mockbridge.db", "The FS path to read for the mock bridge")
	)
	flag.Parse()

	bm := building.NewHub()

//	msb := mock.NewSyncBridge()
//	bm.AddBridge(msb, 5*time.Second)
//	go msb.Run()

//	mab := mock.NewAsyncBridge()
//	bm.AddAsyncBridge(mab)
//	go mab.Run()

	db := &building.BridgeDB{}
	err := db.Open(*dbPath)
	if err != nil {
		log.Printf("Error opening db path %s: %s\n", *dbPath, err.Error())
		os.Exit(1)
	}
	defer db.Close()

	pbb := mock.NewPersistentBridge(db)
	bm.AddBridge(pbb, time.Hour)
	go pbb.Run()

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
