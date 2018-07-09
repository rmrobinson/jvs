package bridge

import (
	"context"
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"
	br "github.com/rmrobinson/bottlerocket-go"
	"github.com/rmrobinson/jvs/service/building"
	"github.com/rmrobinson/jvs/service/building/pb"
)

var (
	houses = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P"}
	maxDeviceID = 16
	baseX10Bridge = &pb.Bridge{
		ModelId: "CM17A",
		ModelName: "Firecracker",
		ModelDescription: "Serial-X10 bridge",
		Manufacturer: "x10.com",
		State: &pb.BridgeState{
			IsPaired: true,
			Version: &pb.BridgeState_Version{
				Api: "1.0.0",
				Sw:  "0.05b3",
			},
		},
	}
	baseX10Device = &pb.Device{
		ModelId: "1",
		ModelName: "X10 Wall Unit",
		ModelDescription: "Plug-in X10 control unit",
		Manufacturer: "x10.com",
	}
)

// BottlerocketBridge offers the standard bridge capabilities over the Bottlerocket X10 USB/serial interface.
type BottlerocketBridge struct {
	br *br.Bottlerocket

	persister building.BridgePersister
}

// NewBottlerocketBridge takes a previously set up bottlerocket handle and exposes it as a bottlerocket bridge.
func NewBottlerocketBridge(bridge *br.Bottlerocket, persister building.BridgePersister) *BottlerocketBridge {
	return &BottlerocketBridge{
		br: bridge,
		persister: persister,
	}
}

func (b *BottlerocketBridge) setup(ctx context.Context) error {
	// Populate the devices
	for _, houseID := range houses {
		for deviceID := 1; deviceID <= maxDeviceID; deviceID++ {
			d := &pb.Device{
				// Id is populated by CreateDevice
				IsActive: false,
				Address: fmt.Sprintf("/x10/%s%d", houseID, deviceID),
				Config: &pb.DeviceConfig{
					Name: "X10 device",
					Description: "Basic X10 device",
				},
				State: &pb.DeviceState{
					Binary: &pb.DeviceState_BinaryState{
						IsOn: false,
					},
				},
			}
			proto.Merge(d, baseX10Device)
			if err := b.persister.CreateDevice(ctx, d); err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *BottlerocketBridge) Bridge(ctx context.Context) (*pb.Bridge, error) {
	bridge, err := b.persister.Bridge(ctx)
	if err != nil {
		return nil, err
	}

	ret := &pb.Bridge{
		Config: &pb.BridgeConfig{
			Address: &pb.Address{
				Usb: &pb.Address_Usb{
					Path: b.br.Path(),
				},
			},
			Timezone: "UTC",
		},
	}
	proto.Merge(ret, baseX10Bridge)
	proto.Merge(ret, bridge)
	return ret, nil
}

func (b *BottlerocketBridge) SetBridgeConfig(ctx context.Context, config *pb.BridgeConfig) error {
	return b.persister.SetBridgeConfig(ctx, config)
}
func (b *BottlerocketBridge) SetBridgeState(ctx context.Context, state *pb.BridgeState) error {
	return b.persister.SetBridgeState(ctx, state)
}

func (b *BottlerocketBridge) SearchForAvailableDevices(context.Context) error {
	return nil
}
func (b *BottlerocketBridge) AvailableDevices(ctx context.Context) ([]*pb.Device, error) {
	return b.persister.AvailableDevices(ctx)
}
func (b *BottlerocketBridge) Devices(ctx context.Context) ([]*pb.Device, error) {
	devices, err := b.persister.Devices(ctx)
	if err != nil {
		return nil, err
	}
	for _, device := range devices {
		proto.Merge(device, baseX10Device)
	}
	return devices, nil
}
func (b *BottlerocketBridge) Device(ctx context.Context, id string) (*pb.Device, error) {
	device, err := b.persister.Device(ctx, id)
	if err != nil {
		return nil, err
	}
	proto.Merge(device, baseX10Device)
	return device, nil
}

func (b *BottlerocketBridge) SetDeviceConfig(ctx context.Context, dev *pb.Device, config *pb.DeviceConfig) error {
	return b.persister.SetDeviceConfig(ctx, dev, config)
}
func (b *BottlerocketBridge) SetDeviceState(ctx context.Context, dev *pb.Device, state *pb.DeviceState) error {
	var err error

	addr := strings.Trim("/x10/", dev.Address)
	if state.Binary.IsOn {
		err = b.br.SendCommand(addr, "ON")
	} else {
		err = b.br.SendCommand(addr, "OFF")
	}

	if err != nil {
		return err
	}

	return b.persister.SetDeviceState(ctx, dev, state)
}
func (b *BottlerocketBridge) AddDevice(ctx context.Context, id string) error {
	// Move the device from available to in use
	return b.persister.AddDevice(ctx, id)
}
func (b *BottlerocketBridge) DeleteDevice(ctx context.Context, id string) error {
	// Move the device from in use to available, and remove the saved values
	return b.persister.DeleteDevice(ctx, id)
}
