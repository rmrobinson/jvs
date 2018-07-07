package bridge

import (
	"context"
	"time"

	hue "github.com/rmrobinson/hue-go"
	"github.com/rmrobinson/jvs/service/device/pb"
)

const (
	stateRefreshInterval = 60 * time.Second
)

// HuePersister is an interface to persisting bridge profiles.
type HuePersister interface {
	Profile(ctx context.Context, bridgeID string) (string, error)
	SaveProfile(ctx context.Context, bridgeID string, username string) error
}

// HueBridge is an implementation of a bridge for the Hue service.
type HueBridge struct {
	bn  BridgeNotifier
	hue *hue.Bridge

	state               hue.BridgeDescription
	lastStateRefreshErr error
	stopStateRefresh    chan struct{}
}

// NewHueBridge takes a previously set up Hue handle and exposes it as a Hue bridge.
func NewHueBridge(notifier BridgeNotifier, bridge *hue.Bridge) *HueBridge {
	ret := &HueBridge{
		bn:               notifier,
		hue:              bridge,
		stopStateRefresh: make(chan struct{}),
	}

	go ret.stateRefresher()

	return ret
}

func (b *HueBridge) Close() {
	close(b.stopStateRefresh)
}

func (b *HueBridge) ID() string {
	return b.state.Device.SerialNumber
}

func (b *HueBridge) ModelID() string {
	return b.state.Device.ModelNumber
}
func (b *HueBridge) ModelName() string {
	return b.state.Device.ModelName
}
func (b *HueBridge) ModelDescription() string {
	return b.state.Device.ModelDescription
}
func (b *HueBridge) Manufacturer() string {
	return b.state.Device.Manufacturer
}
func (b *HueBridge) IconURLs() []string {
	ret := []string{}

	for _, icon := range b.state.Device.Icons {
		ret = append(ret, b.state.UrlBase+"/"+icon.FileName)
	}

	return ret
}
func (b *HueBridge) Name() string {
	return b.state.Device.FriendlyName
}
func (b *HueBridge) SetName(name string) error {
	c := &hue.ConfigArg{}
	c.SetName(name)
	return b.hue.SetConfig(c)
}

func (b *HueBridge) Devices() ([]*pb.Device, error) {
	lights, err := b.hue.Lights()

	if err != nil {
		return nil, err
	}

	sensors, err := b.hue.Sensors()

	if err != nil {
		return nil, err
	}

	var devices []*pb.Device

	for _, light := range lights {
		d := convertLightToDevice(light)
		d.Address = b.ID()
		devices = append(devices, d)
	}
	for _, sensor := range sensors {
		d := convertSensorToDevice(sensor)
		d.Address = b.ID()
		devices = append(devices, d)
	}

	return devices, nil
}

func (b *HueBridge) stateRefresher() {
	ticker := time.NewTicker(stateRefreshInterval)

	for {
		select {
		case <-ticker.C:
			state, err := b.hue.Description()

			if err != nil {
				b.lastStateRefreshErr = err
				continue
			}

			if b.state.UrlBase != state.UrlBase {
				b.state = state
				b.bn.BridgeUpdated(b.ID())
			}
		case <-b.stopStateRefresh:
			ticker.Stop()
			return
		}
	}
}

func convertLightToDevice(l hue.Light) *pb.Device {
	d := &pb.Device{}
	d.Reset()

	d.Id = l.UniqueId
	d.Path = "light/" + l.Id
	d.IsActive = true

	d.Manufacturer = l.ManufacturerName
	d.ModelId = l.ModelId

	config := &pb.DeviceConfig{}
	d.Config = config

	config.Name = l.Name

	state := &pb.DeviceState{}
	d.State = state

	state.IsReachable = l.State.Reachable
	state.Version = l.SwVersion

	state.Binary = &pb.DeviceState_BinaryState{IsOn: l.State.On}

	if l.Model == "Dimmable light" ||
		l.Model == "Color light" ||
		l.Model == "Extended color light" {

		// We are hardcoded to only supporting a uint8 worth of brightness values
		d.Range = &pb.Device_RangeDevice{Minimum: 0, Maximum: 254}

		state.Range = &pb.DeviceState_RangeState{Value: int32(l.State.Brightness)}
	}

	if l.Model == "Color light" ||
		l.Model == "Extended color light" ||
		l.Model == "Color temperature light" {
		if l.State.ColorMode == "xy" {
			xy := hue.XY{X: l.State.XY[0], Y: l.State.XY[1]}
			rgb := hue.RGB{}
			rgb.FromXY(xy, l.ModelId)
			state.ColorRgb = &pb.DeviceState_RGBState{Red: int32(rgb.Red), Blue: int32(rgb.Blue), Green: int32(rgb.Green)}
		} else if l.State.ColorMode == "ct" {
			rgb := hue.RGB{}
			rgb.FromCT(l.State.ColorTemperature)
			state.ColorRgb = &pb.DeviceState_RGBState{Red: int32(rgb.Red), Blue: int32(rgb.Blue), Green: int32(rgb.Green)}
		} else if l.State.ColorMode == "hs" {
			hsb := hue.HSB{Hue: l.State.Hue, Saturation: l.State.Saturation, Brightness: l.State.Brightness}
			rgb := hue.RGB{}
			rgb.FromHSB(hsb)
			state.ColorRgb = &pb.DeviceState_RGBState{Red: int32(rgb.Red), Blue: int32(rgb.Blue), Green: int32(rgb.Green)}
		}
	}

	return d
}

func convertSensorToDevice(s hue.Sensor) *pb.Device {
	d := &pb.Device{}
	d.Reset()

	d.Id = s.UniqueId
	d.Path = "sensor/" + s.Id

	d.Manufacturer = s.ManufacturerName
	d.ModelId = s.ModelId

	config := &pb.DeviceConfig{}
	d.Config = config

	config.Name = s.Name

	state := &pb.DeviceState{}
	d.State = state

	state.IsReachable = s.Config.Reachable
	state.Version = s.SwVersion

	if s.Type == "ZGPSwitch" {
		button := &pb.DeviceState_ButtonState{}
		button.IsOn = true

		switch s.State.ButtonEvent {
		case 34:
			button.Id = 1
		case 16:
			button.Id = 2
		case 17:
			button.Id = 3
		case 18:
			button.Id = 4
		}

		state.Button = append(state.Button, button)
	} else if s.Type == "ZLLSwitch" {
		button := &pb.DeviceState_ButtonState{}

		switch s.State.ButtonEvent {
		case 1000, 1001, 1002, 1003:
			button.Id = 1
			button.IsOn = true
		case 2000, 2001, 2002, 2003:
			button.Id = 2
			button.IsOn = true
		case 3000, 3001, 3002, 3003:
			button.Id = 3
			button.IsOn = true
		case 4000, 4001, 4002, 4003:
			button.Id = 4
			button.IsOn = true
		}

		state.Button = append(state.Button, button)
	} else if s.Type == "ZLLPresence" {
		d.Range = &pb.Device_RangeDevice{Minimum: 0, Maximum: s.SensitivityMax}
		state.Range = &pb.DeviceState_RangeState{Value: s.Sensitivity}

		state.Presence = &pb.DeviceState_PresenceState{IsPresent: s.State.Presence}
	} else if s.Type == "ZLLTemperature" {
		state.Temperature = &pb.DeviceState_TemperatureState{TemperatureCelsius: s.State.Temperature / 100}
	}

	return d
}
