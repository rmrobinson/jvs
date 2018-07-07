package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	hue "github.com/rmrobinson/hue-go"
	"github.com/rmrobinson/jvs/service/building"
	"github.com/rmrobinson/jvs/service/building/pb"
	"google.golang.org/grpc"
)

func main() {
	var (
		port      = flag.Int("port", 1337, "The port for the deviced process to listen on")
		hueDBPath = flag.String("hueDBPath", "", "The path to the Hue pairing DB")
	)

	flag.Parse()

	bm := building.NewHub()

	// Bottlerocket setup is done via the API
	// Proxy setup is done via the API

	// Hue setup is done here since we don't configure bridges directly
	if len(*hueDBPath) > 0 {
		hueDB := &building.HueDB{}
		err := hueDB.Open(*hueDBPath)

		if err != nil {
			log.Printf("Error initializing the Hue DB: %s\n", err.Error())
			os.Exit(1)
		}
		defer hueDB.Close()

		// We don't have any bridges specified explicitly; they are autodiscovered.
		// Start the autodiscovery process
		go hueAutodiscover(bm, hueDB)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Printf("Error initializing listener: %s\n", err.Error())
		os.Exit(1)
	}
	defer lis.Close()

	api := building.NewAPI(bm)

	grpcServer := grpc.NewServer()
	pb.RegisterBridgeManagerServer(grpcServer, api)
	grpcServer.Serve(lis)
}

func hueAutodiscover(bm *building.Hub, p building.HuePersister) {
	bridges := make(chan hue.Bridge)

	locator := hue.NewLocator()
	go locator.Run(bridges)

	for {
		bridge := <-bridges

		username, err := p.Profile(context.Background(), bridge.Id())

		if err != nil {
			log.Printf("Unable to get pairing for ID: %s\n", err)
		} else {
			bridge.Username = username
		}

		bm.AddBridge(building.NewHueBridge(bm, &bridge))
	}
}
