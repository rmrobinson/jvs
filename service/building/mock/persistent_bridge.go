package mock

import (
	"context"
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/rmrobinson/jvs/service/building"
	"github.com/rmrobinson/jvs/service/building/pb"
)

var (
	basePersistentBridge = &pb.Bridge{
		Type:             pb.BridgeType_Loopback,
		ModelId:          "PBTest1",
		ModelName:        "Test Persistent Bridge",
		ModelDescription: "Bridge for testing persistent operations",
		Manufacturer:     "Faltung Systems",
	}
	basePersistentDevice = &pb.Device{
		ModelId:          "PBDevice1",
		ModelName:        "Test Persistent Device",
		ModelDescription: "Device for testing persistent operations",
		Manufacturer:     "Faltung Systems",
	}
)

type PersistentBridge struct {
	bridgeID  string
	persister building.BridgePersister
}

func NewPersistentBridge(persister building.BridgePersister) *PersistentBridge {
	ret := &PersistentBridge{
		persister: persister,
	}
	return ret
}

func (b *PersistentBridge) setup() {
	bridgeConfig := &pb.BridgeConfig{
		Name: "Test Persistent Bridge",
	}
	id, err := b.persister.CreateBridge(context.Background(), bridgeConfig)
	if err != nil {
		log.Printf("Error creating bridge: %s\n", err.Error())
		return
	}
	b.bridgeID = id

	for i := 0; i < 5; i++ {
		d := newDevice()
		d.d.IsActive = false
		d.d.Address = fmt.Sprintf("/test/%d", i)

		b.persister.CreateDevice(context.Background(), d.d)
	}
	for i := 5; i < 7; i++ {
		d := newDevice()
		d.d.IsActive = true
		d.d.Address = fmt.Sprintf("/test/%d", i)

		b.persister.CreateDevice(context.Background(), d.d)
	}

	log.Printf("Created bridge")
}
func (b *PersistentBridge) Run() {
	_, err := b.persister.Bridge(context.Background())
	if err == building.ErrDatabaseNotSetup {
		b.setup()
	}
}

func (b *PersistentBridge) Bridge(ctx context.Context) (*pb.Bridge, error) {
	bridge, err := b.persister.Bridge(ctx)
	if err != nil {
		return nil, err
	}
	proto.Merge(bridge, basePersistentBridge)
	return bridge, nil
}
func (b *PersistentBridge) SetBridgeConfig(ctx context.Context, config *pb.BridgeConfig) error {
	return b.persister.SetBridgeConfig(ctx, config)
}
func (b *PersistentBridge) SetBridgeState(ctx context.Context, state *pb.BridgeState) error {
	return b.persister.SetBridgeState(ctx, state)
}

func (b *PersistentBridge) SearchForAvailableDevices(ctx context.Context) error {
	return b.persister.SearchForAvailableDevices(ctx)
}
func (b *PersistentBridge) AvailableDevices(ctx context.Context) ([]*pb.Device, error) {
	devices, err := b.persister.AvailableDevices(ctx)
	if err != nil {
		return nil, err
	}
	for _, device := range devices {
		proto.Merge(device, basePersistentDevice)
	}
	return devices, nil
}

func (b *PersistentBridge) Devices(ctx context.Context) ([]*pb.Device, error) {
	devices, err := b.persister.Devices(ctx)
	if err != nil {
		return nil, err
	}
	for _, device := range devices {
		proto.Merge(device, basePersistentDevice)
	}
	return devices, nil
}
func (b *PersistentBridge) Device(ctx context.Context, id string) (*pb.Device, error) {
	device, err := b.persister.Device(ctx, id)
	if err != nil {
		return nil, err
	}
	proto.Merge(device, basePersistentDevice)
	return device, nil
}

func (b *PersistentBridge) SetDeviceConfig(ctx context.Context, dev *pb.Device, config *pb.DeviceConfig) error {
	return b.persister.SetDeviceConfig(ctx, dev, config)
}
func (b *PersistentBridge) SetDeviceState(ctx context.Context, dev *pb.Device, state *pb.DeviceState) error {
	return b.persister.SetDeviceState(ctx, dev, state)
}
func (b *PersistentBridge) AddDevice(ctx context.Context, id string) error {
	return b.persister.AddDevice(ctx, id)
}
func (b *PersistentBridge) DeleteDevice(ctx context.Context, id string) error {
	return b.persister.DeleteDevice(ctx, id)
}
