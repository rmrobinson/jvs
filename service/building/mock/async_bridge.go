package mock

import (
	"context"
	"math/rand"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/rmrobinson/jvs/service/building"
	"github.com/rmrobinson/jvs/service/building/pb"
)

type AsyncBridge struct {
	b *pb.Bridge

	availDevices []*pb.Device
	devices      map[string]*device

	notifier building.Notifier
}

func NewAsyncBridge() *AsyncBridge {
	ret := &AsyncBridge{
		b: &pb.Bridge{
			Id:               uuid.New().String(),
			Type:             pb.BridgeType_Loopback,
			Mode:             pb.BridgeMode_Active,
			ModeReason:       "",
			ModelId:          "ABTest1",
			ModelName:        "Test Async Bridge",
			ModelDescription: "Bridge for testing async operations",
			Manufacturer:     "Faltung Systems",
			State: &pb.BridgeState{
				IsPaired: true,
			},
		},
		devices: map[string]*device{},
	}

	count := rand.Intn(5)
	for i := 0; i < count; i++ {
		d := newDevice()
		ret.devices[d.d.Id] = d
	}

	return ret
}

func (ab *AsyncBridge) Run() {
	t := time.NewTicker(7 * time.Second)
	for {
		select {
		case <-t.C:
			for id, d := range ab.devices {
				d.update()
				ab.devices[id] = d
				ab.notifier.DeviceUpdated(ab.b.Id, d.d)
			}
		}
	}
}

func (ab *AsyncBridge) SetNotifier(n building.Notifier) {
	ab.notifier = n
}

func (ab *AsyncBridge) Bridge(context.Context) (*pb.Bridge, error) {
	return ab.b, nil
}
func (ab *AsyncBridge) SetBridgeConfig(ctx context.Context, config *pb.BridgeConfig) error {
	ab.b.Config = config
	return nil
}
func (ab *AsyncBridge) SetBridgeState(ctx context.Context, state *pb.BridgeState) error {
	return ErrReadOnly
}

func (ab *AsyncBridge) SearchForAvailableDevices(context.Context) error {
	if len(ab.availDevices) < 1 {
		count := rand.Intn(5)

		for i := 0; i < count; i++ {
			ab.availDevices = append(ab.availDevices, newDevice().d)
		}
	}

	return nil
}
func (ab *AsyncBridge) AvailableDevices(context.Context) ([]*pb.Device, error) {
	return ab.availDevices, nil
}

func (ab *AsyncBridge) Devices(context.Context) ([]*pb.Device, error) {
	var ret []*pb.Device
	for _, d := range ab.devices {
		ret = append(ret, d.d)
	}
	return ret, nil
}
func (ab *AsyncBridge) Device(ctx context.Context, id string) (*pb.Device, error) {
	if d, ok := ab.devices[id]; ok {
		return d.d, nil
	}
	return nil, ErrDeviceNotPresent
}

func (ab *AsyncBridge) SetDeviceConfig(ctx context.Context, dev *pb.Device, config *pb.DeviceConfig) error {
	var d *device
	var ok bool
	if d, ok = ab.devices[dev.Id]; !ok {
		return ErrDeviceNotPresent
	}

	d.d.Config = proto.Clone(config).(*pb.DeviceConfig)
	return nil
}
func (ab *AsyncBridge) SetDeviceState(ctx context.Context, dev *pb.Device, state *pb.DeviceState) error {
	var d *device
	var ok bool
	if d, ok = ab.devices[dev.Id]; !ok {
		return ErrDeviceNotPresent
	}

	d.d.State = proto.Clone(state).(*pb.DeviceState)
	return nil

}
func (ab *AsyncBridge) AddDevice(ctx context.Context, id string) error {
	var d *pb.Device
	found := false
	for idx, availDevice := range ab.availDevices {
		if availDevice.Id == id {
			d = availDevice
			ab.availDevices = append(ab.availDevices[:idx], ab.availDevices[idx+1:]...)
			found = true
			break
		}
	}

	if !found {
		return ErrDeviceNotPresent
	}

	ab.devices[d.Id] = &device{
		d: d,
	}
	return nil
}
func (ab *AsyncBridge) DeleteDevice(ctx context.Context, id string) error {
	if _, ok := ab.devices[id]; !ok {
		return ErrDeviceNotPresent
	}
	delete(ab.devices, id)
	return nil
}
