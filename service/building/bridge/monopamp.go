package bridge

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/rmrobinson/jvs/service/building"
	"github.com/rmrobinson/jvs/service/building/pb"
	monopamp "github.com/rmrobinson/monoprice-amp-go"
)

var (
	maxZoneID          = 6
	baseMonopAmpBridge = &pb.Bridge{
		ModelId:          "10761",
		ModelName:        "Monoprice Amp",
		ModelDescription: "6 Zone Home Audio Multizone Controller",
		Manufacturer:     "Monoprice",
	}
	baseMonopAmpDevice = &pb.Device{
		ModelId:          "10761",
		ModelName:        "Zone",
		ModelDescription: "Monoprice Amp Zone",
		Manufacturer:     "Monoprice",
	}
)

// MonopAmpBridge is an implementation of a bridge for the Monoprice amp/stero output device.
type MonopAmpBridge struct {
	amp *monopamp.SerialAmplifier

	persister building.BridgePersister
}

// NewMonopAmpBridge takes a previously set up MonopAmp handle and exposes it as a MonopAmp bridge.
func NewMonopAmpBridge(amp *monopamp.SerialAmplifier, persister building.BridgePersister) *MonopAmpBridge {
	ret := &MonopAmpBridge{
		amp:       amp,
		persister: persister,
	}

	return ret
}

func (b *MonopAmpBridge) setup(ctx context.Context) error {
	// Populate the devices
	for zoneID := 1; zoneID <= maxZoneID; zoneID++ {
		d := &pb.Device{
			// Id is populated by CreateDevice
			IsActive: true,
			Address:  fmt.Sprintf("/zone/%d", zoneID),
			Config: &pb.DeviceConfig{
				Name:        fmt.Sprintf("Amp Zone %d", zoneID),
				Description: "Amplifier output for the specified zone",
			},
		}
		proto.Merge(d, baseMonopAmpDevice)
		if err := b.persister.CreateDevice(ctx, d); err != nil {
			return err
		}
	}

	return nil
}

func (b *MonopAmpBridge) Bridge(ctx context.Context) (*pb.Bridge, error) {
	bridge, err := b.persister.Bridge(ctx)
	if err != nil {
		return nil, err
	}

	ret := &pb.Bridge{
		Config: &pb.BridgeConfig{
			Address: &pb.Address{
				Usb: &pb.Address_Usb{
					Path: "",
				},
			},
			Timezone: "UTC",
		},
	}
	proto.Merge(ret, baseMonopAmpBridge)
	proto.Merge(ret, bridge)
	return ret, nil
}

func (b *MonopAmpBridge) SetBridgeConfig(ctx context.Context, config *pb.BridgeConfig) error {
	return b.persister.SetBridgeConfig(ctx, config)
}
func (b *MonopAmpBridge) SetBridgeState(ctx context.Context, state *pb.BridgeState) error {
	return b.persister.SetBridgeState(ctx, state)
}

func (b *MonopAmpBridge) SearchForAvailableDevices(context.Context) error {
	return nil
}
func (b *MonopAmpBridge) AvailableDevices(ctx context.Context) ([]*pb.Device, error) {
	return b.persister.AvailableDevices(ctx)
}
func (b *MonopAmpBridge) Devices(ctx context.Context) ([]*pb.Device, error) {
	devices, err := b.persister.Devices(ctx)
	if err != nil {
		return nil, err
	}
	for _, device := range devices {
		proto.Merge(device, baseMonopAmpDevice)
	}
	return devices, nil
}
func (b *MonopAmpBridge) Device(ctx context.Context, id string) (*pb.Device, error) {
	device, err := b.persister.Device(ctx, id)
	if err != nil {
		return nil, err
	}
	proto.Merge(device, baseMonopAmpDevice)
	return device, nil
}

func (b *MonopAmpBridge) SetDeviceConfig(ctx context.Context, dev *pb.Device, config *pb.DeviceConfig) error {
	return b.persister.SetDeviceConfig(ctx, dev, config)
}
func (b *MonopAmpBridge) SetDeviceState(ctx context.Context, dev *pb.Device, state *pb.DeviceState) error {
	addr, err := strconv.ParseInt(strings.Trim("/zone/", dev.Address), 10, 32)
	if err != nil {
		return err
	}

	zone := b.amp.Zone(int(addr))
	if zone == nil {
		return errors.New("invalid address")
	}

	err = zone.SetPower(state.Binary.IsOn)
	if err != nil {
		return err
	}

	// TODO: we will have other things to set here eventually.
	return nil
}
func (b *MonopAmpBridge) AddDevice(ctx context.Context, id string) error {
	// Move the device from available to in use
	return b.persister.AddDevice(ctx, id)
}
func (b *MonopAmpBridge) DeleteDevice(ctx context.Context, id string) error {
	// Move the device from in use to available, and remove the saved values
	return b.persister.DeleteDevice(ctx, id)
}
