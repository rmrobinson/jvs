package mock

import (
	"errors"
	"math/rand"

	"github.com/google/uuid"
	"github.com/rmrobinson/jvs/service/building/pb"
)

var (
	ErrDeviceNotPresent = errors.New("device not present")
	ErrReadOnly         = errors.New("read only")
)

type device struct {
	d *pb.Device
}

func newDevice() *device {
	id := uuid.New().String()
	d := &device{
		d: &pb.Device{
			Id:               id,
			IsActive:         true,
			ModelId:          "TD1",
			ModelName:        "Test Device",
			ModelDescription: "Device for testing purposes",
			Manufacturer:     "Faltung Systems",
			Address:          "/testing/" + id,
			Config: &pb.DeviceConfig{
				Name:        "Device X",
				Description: "Device randomly generated",
			},
			State: &pb.DeviceState{
				IsReachable: true,
				Binary:      &pb.DeviceState_BinaryState{},
			},
		},
	}

	c := rand.Intn(10)

	if c%2 == 0 {
		d.d.State.Binary.IsOn = true
	}

	if c > 5 {
		d.d.Range = &pb.Device_RangeDevice{
			Maximum: 10,
			Minimum: 0,
		}
		d.d.State.Range = &pb.DeviceState_RangeState{
			Value: int32(c),
		}
	}

	if c > 7 {
		d.d.State.ColorRgb = &pb.DeviceState_RGBState{
			Red:   int32(c * 2 * 10),
			Green: int32(c * 3 * 9),
			Blue:  int32(c * 5 * 7),
		}
	}

	return d
}

func (d *device) update() {
	c := rand.Intn(10)

	if c > 4 && d.d.State.Binary != nil {
		d.d.State.Binary.IsOn = !d.d.State.Binary.IsOn
	}

	if c > 6 && d.d.State.Range != nil {
		d.d.State.Range.Value = rand.Int31n(d.d.Range.Maximum)
	}

	if c > 8 && d.d.State.ColorRgb != nil {
		d.d.State.ColorRgb.Red = rand.Int31n(255) * 2
		d.d.State.ColorRgb.Green = rand.Int31n(255) * 3
		d.d.State.ColorRgb.Blue = rand.Int31n(255) * 2
	}
}
