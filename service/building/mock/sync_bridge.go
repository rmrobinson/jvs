package mock

import (
	"math/rand"

	"github.com/google/uuid"
	"github.com/rmrobinson/jvs/service/building/pb"
	"time"
)

type SyncBridge struct {
	b *pb.Bridge

	availDevices []*pb.Device
	devices      map[string]*device
}

func NewSyncBridge() *SyncBridge {
	ret := &SyncBridge{
		b: &pb.Bridge{
			Id:               uuid.New().String(),
			Type:             pb.BridgeType_Loopback,
			Mode:             pb.BridgeMode_Active,
			ModeReason:       "",
			ModelId:          "SBTest1",
			ModelName:        "Test Sync Bridge",
			ModelDescription: "Bridge for testing sync operations",
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

func (sb *SyncBridge) Run() {
	t := time.NewTicker(time.Second * 6)
	for {
		select {
		case <- t.C:
			for id, d := range sb.devices {
				d.update()
				sb.devices[id] = d
			}
		}
	}
}

func (sb *SyncBridge) Bridge() (*pb.Bridge, error) {
	return sb.b, nil
}
func (sb *SyncBridge) SetBridgeConfig(config *pb.BridgeConfig) error {
	sb.b.Config = config
	return nil
}
func (sb *SyncBridge) SetBridgeState(state *pb.BridgeState) error {
	return ErrReadOnly
}

func (sb *SyncBridge) SearchForAvailableDevices() error {
	if len(sb.availDevices) < 1 {
		count := rand.Intn(5)

		for i := 0; i < count; i++ {
			sb.availDevices = append(sb.availDevices, newDevice().d)
		}
	}

	return nil
}
func (sb *SyncBridge) AvailableDevices() ([]*pb.Device, error) {
	return sb.availDevices, nil
}

func (sb *SyncBridge) Devices() ([]*pb.Device, error) {
	var ret []*pb.Device
	for _, d := range sb.devices {
		ret = append(ret, d.d)
	}
	return ret, nil
}
func (sb *SyncBridge) Device(id string) (*pb.Device, error) {
	if d, ok := sb.devices[id]; ok {
		return d.d, nil
	}
	return nil, ErrDeviceNotPresent
}

func (sb *SyncBridge) SetDeviceConfig(id string, config *pb.DeviceConfig) error {
	var d *device
	var ok bool
	if d, ok = sb.devices[id]; !ok {
		return ErrDeviceNotPresent
	}

	d.d.Config = config
	return nil
}
func (sb *SyncBridge) SetDeviceState(id string, state *pb.DeviceState) error {
	var d *device
	var ok bool
	if d, ok = sb.devices[id]; !ok {
		return ErrDeviceNotPresent
	}

	d.d.State = state
	return nil

}
func (sb *SyncBridge) AddDevice(id string) error {
	var d *pb.Device
	found := false
	for idx, availDevice := range sb.availDevices {
		if availDevice.Id == id {
			d = availDevice
			sb.availDevices = append(sb.availDevices[:idx], sb.availDevices[idx+1:]...)
			found = true
			break
		}
	}

	if !found {
		return ErrDeviceNotPresent
	}

	sb.devices[d.Id] = &device{
		d: d,
	}
	return nil
}
func (sb *SyncBridge) DeleteDevice(id string) error {
	if _, ok := sb.devices[id]; !ok {
		return ErrDeviceNotPresent
	}
	delete(sb.devices, id)
	return nil
}
