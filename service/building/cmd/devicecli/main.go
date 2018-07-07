package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"

	"github.com/rmrobinson/jvs/service/building/pb"
	"google.golang.org/grpc"
	"sync"
)

func getBridges(bc pb.BridgeManagerClient) {
	getResp, err := bc.GetBridges(context.Background(), &pb.GetBridgesRequest{})

	if err != nil {
		fmt.Printf("Unable to get bridges: %s\n", err.Error())
		return
	}

	fmt.Printf("Got bridges\n")
	for _, bridge := range getResp.Bridges {
		fmt.Printf("%+v\n", bridge)
	}
}

func monitorBridges(bc pb.BridgeManagerClient) {
	stream, err := bc.WatchBridges(context.Background(), &pb.WatchBridgesRequest{})

	if err != nil {
		return
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("Error while watching bridges: %v", err)
			break
		}

		log.Printf("Change: %v, Bridge: %+v\n", msg.Action, msg.Bridge)
	}
}
func monitorDevices(dc pb.DeviceManagerClient) {
	stream, err := dc.WatchDevices(context.Background(), &pb.WatchDevicesRequest{})

	if err != nil {
		return
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("Error while watching devices: %v", err)
			break
		}

		log.Printf("Change: %v, Device: %+v\n", msg.Action, msg.Device)
	}
}
func main() {
	var (
		addr = flag.String("addr", "", "The address to connect to")
		mode = flag.String("mode", "", "The mode of operation for the client")
	)

	flag.Parse()

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(*addr, opts...)

	if err != nil {
		fmt.Printf("Unable to connect: %s\n", err.Error())
		return
	}

	bridgeClient := pb.NewBridgeManagerClient(conn)
	deviceClient := pb.NewDeviceManagerClient(conn)

	switch *mode {
	case "get":
		getBridges(bridgeClient)
	case "monitor":
		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			monitorBridges(bridgeClient)
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			monitorDevices(deviceClient)
		}()
		wg.Wait()
	default:
		fmt.Printf("Unknown mode specified")
	}
}
