package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"

	"github.com/rmrobinson/jvs/service/device/pb"
	"google.golang.org/grpc"
)

func addLoopback(bc pb.BridgeManagerClient) {
	addResp, err := bc.AddBridge(context.Background(), &pb.AddBridgeRequest{
		Type: pb.BridgeType_Loopback,
	})

	if err != nil {
		fmt.Printf("Unable to add loopback bridge: %s\n", err.Error())
		return
	}

	fmt.Printf("Created bridge %+v\n", addResp.Bridge)
}

func addProxy(bc pb.BridgeManagerClient) {
	addResp, err := bc.AddBridge(context.Background(), &pb.AddBridgeRequest{
		Type: pb.BridgeType_Proxy,
		Config: &pb.BridgeConfig{
			Id: "68c0987b-62ea-4d16-8347-9340feca87d0",
			Address: &pb.Address{
				Ip: &pb.Address_Ip{
					Host: "127.0.0.1",
					Port: 1338,
				},
			},
		},
	})

	if err != nil {
		fmt.Printf("Unable to add loopback bridge: %s\n", err.Error())
		return
	}

	fmt.Printf("Created bridge %+v\n", addResp.Bridge)
}

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

func monitor(bc pb.BridgeManagerClient) {
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

	switch *mode {
	case "proxy":
		addProxy(bridgeClient)
	case "loopback":
		addLoopback(bridgeClient)
	case "get":
		getBridges(bridgeClient)
	case "monitor":
		monitor(bridgeClient)
	default:
		fmt.Printf("Unknown mode specified")
	}
}
