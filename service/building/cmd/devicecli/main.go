package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/rmrobinson/jvs/service/building/pb"
	"google.golang.org/grpc"
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

func setBridgeName(bc pb.BridgeManagerClient, id string, name string) {
	req := &pb.SetBridgeConfigRequest{
		Id: id,
		Config: &pb.BridgeConfig{
			Name: name,
		},
	}

	setResp, err := bc.SetBridgeConfig(context.Background(), req)
	if err != nil {
		fmt.Printf("Unable to set bridge name: %s\n", err.Error())
		return
	}

	fmt.Printf("Set bridge name\n")
	fmt.Printf("%+v\n", setResp.Bridge)
}


func getDevices(dc pb.DeviceManagerClient) {
	getResp, err := dc.GetDevices(context.Background(), &pb.GetDevicesRequest{})
	if err != nil {
		fmt.Printf("Unable to get devices: %s\n", err.Error())
		return
	}

	fmt.Printf("Got devices\n")
	for _, device := range getResp.Devices {
		fmt.Printf("%+v\n", device)
	}
}

func setDeviceName(dc pb.DeviceManagerClient, id string, name string) {
	req := &pb.SetDeviceConfigRequest{
		Id: id,
		Config: &pb.DeviceConfig{
			Name: name,
			Description: "Manually set",
		},
	}

	setResp, err := dc.SetDeviceConfig(context.Background(), req)
	if err != nil {
		fmt.Printf("Unable to set device name: %s\n", err.Error())
		return
	}

	fmt.Printf("Set device name\n")
	fmt.Printf("%+v\n", setResp.Device)
}

func setDeviceIsOn(dc pb.DeviceManagerClient, id string, isOn bool) {
	d, err := dc.GetDevice(context.Background(), &pb.GetDeviceRequest{Id: id})

	if err != nil {
		fmt.Printf("Unable to get device: %s\n", err.Error())
		return
	}

	d.Device.State.Binary.IsOn = isOn
	req := &pb.SetDeviceStateRequest{
		Id: id,
		State: d.Device.State,
	}

	setResp, err := dc.SetDeviceState(context.Background(), req)
	if err != nil {
		fmt.Printf("Unable to set device state: %s\n", err.Error())
		return
	}

	fmt.Printf("Set device isOn\n")
	fmt.Printf("%+v\n", setResp.Device)
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
		id = flag.String("id", "", "The device ID to change")
		name = flag.String("name", "", "The device name to set")
		on = flag.Bool("on", false, "The device ison state to set")
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
	case "getBridges":
		getBridges(bridgeClient)
	case "setBridgeConfig":
		setBridgeName(bridgeClient, *id, *name)
	case "getDevices":
		getDevices(deviceClient)
	case "setDeviceConfig":
		setDeviceName(deviceClient, *id, *name)
	case "setDeviceState":
		setDeviceIsOn(deviceClient, *id, *on)
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
