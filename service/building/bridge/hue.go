package bridge

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/rmrobinson/hue-go"
	"github.com/rmrobinson/jvs/service/building"
	"github.com/rmrobinson/jvs/service/building/pb"
)


var (
	ErrHueAddressInvalid = errors.New("hue address invalid")
	ErrHueResponseError = errors.New("hue response error")
	ErrDeviceLacksRangeCapability = errors.New("invalid argument supplied: device lacks range capabilities")
	ErrDeviceRangeLimitExceeded = errors.New("invalid argument supplied: range value outside of allowed values")
	lightAddrPrefix   = "/light/"
	sensorAddrPrefix  = "/sensor/"
)

func addrToLight(addr string) int {
	id, err := strconv.ParseInt(strings.TrimPrefix(addr, lightAddrPrefix), 10, 32)
	if err != nil {
		return 0 // this is an invalid ID
	}
	return int(id)
}
func addrToSensor(addr string) int {
	id, err := strconv.ParseInt(strings.TrimPrefix(addr, sensorAddrPrefix), 10, 32)
	if err != nil {
		return 0 // this is an invalid ID
	}
	return int(id)
}

// HueBridge is an implementation of a bridge for the Friends of Hue system.
type HueBridge struct {
	bridge *hue.Bridge
}

// NewHueBridge takes a previously set up Hue handle and exposes it as a Hue bridge.
func NewHueBridge(bridge *hue.Bridge) *HueBridge {
	return &HueBridge{
		bridge:       bridge,
	}
}

// Setup seeds the persistent store with the proper data
func (b *HueBridge) Setup(ctx context.Context) error {
	return nil
}

// Bridge retrieves the persisted state of the bridge from the backing store.
func (b *HueBridge) Bridge(ctx context.Context) (*pb.Bridge, error) {
	desc, err := b.bridge.Description()
	if err != nil {
		return nil, err
	}
	config, err := b.bridge.Config()
	if err != nil {
		return nil, err
	}

	ret := &pb.Bridge{
		Id: config.ID,
		ModelId: desc.Device.ModelNumber,
		ModelName: desc.Device.ModelName,
		ModelDescription: desc.Device.ModelDescription,
		Manufacturer: desc.Device.Manufacturer,
		Config: &pb.BridgeConfig{
			Name: desc.Device.FriendlyName,
			Address: &pb.Address{
				Ip: &pb.Address_Ip{
					Host: config.IPAddress,
					Netmask: config.SubnetMask,
					Gateway: config.GatewayAddress,
				},
			},
			Timezone: config.Timezone,
		},
		State: &pb.BridgeState{
			Zigbee: &pb.BridgeState_Zigbee{
				Channel: config.ZigbeeChannel,
			},
			Version: &pb.BridgeState_Version{
				Sw: config.SwVersion,
				Api: config.APIVersion,
			},
		},
	}

	for _, icon := range desc.Device.Icons {
		ret.IconUrl = append(ret.IconUrl, desc.URLBase+"/"+icon.FileName)
	}

	return ret, nil
}

// SetBridgeConfig persists the new bridge config on the Hue bridge.
// Only the name can currently be changed.
// TODO: support setting the static IP of the bridge.
func (b *HueBridge) SetBridgeConfig(ctx context.Context, config *pb.BridgeConfig) error {
	updatedConfig := &hue.ConfigArg{}
	updatedConfig.SetName(config.Name)
	return b.bridge.SetConfig(updatedConfig)
}
// SetBridgeState persists the new bridge state on the Hue bridge.
func (b *HueBridge) SetBridgeState(ctx context.Context, state *pb.BridgeState) error {
	return building.ErrOperationNotSupported
}

// SearchForAvailableDevices is a noop that returns immediately (nothing to search for).
func (b *HueBridge) SearchForAvailableDevices(context.Context) error {
	if err := b.bridge.SearchForNewLights(); err != nil {
		return err
	}
	if err := b.bridge.SearchForNewSensors(); err != nil {
		return err
	}
	return nil
}
// AvailableDevices returns an empty result as all devices are always available; never 'to be added'.
func (b *HueBridge) AvailableDevices(ctx context.Context) ([]*pb.Device, error) {
	lights, err := b.bridge.NewLights()
	if err != nil {
		return nil, err
	}
	sensors, err := b.bridge.NewSensors()
	if err != nil {
		return nil, err
	}

	var devices []*pb.Device

	for _, light := range lights {
		devices = append(devices, &pb.Device{
			Address: lightAddrPrefix + light.ID,
			Config: &pb.DeviceConfig{
				Name: light.Name,
			},
		})
	}
	for _, sensor := range sensors {
		devices = append(devices, &pb.Device{
			Address: sensorAddrPrefix + sensor.ID,
			Config: &pb.DeviceConfig{
				Name: sensor.Name,
			},
		})
	}

	return devices, nil
}
// Devices retrieves the list of lights and sensors from the bridge along with their current states.
func (b *HueBridge) Devices(ctx context.Context) ([]*pb.Device, error) {
	lights, err := b.bridge.Lights()
	if err != nil {
		return nil, err
	}

	sensors, err := b.bridge.Sensors()
	if err != nil {
		return nil, err
	}

	var devices []*pb.Device

	for _, light := range lights {
		d := convertLightToDevice(light)
		devices = append(devices, d)
	}
	for _, sensor := range sensors {
		d := convertSensorToDevice(sensor)
		devices = append(devices, d)
	}

	return devices, nil
}
// Device retrieves the specified device ID.
func (b *HueBridge) Device(ctx context.Context, id string) (*pb.Device, error) {
	devices, err := b.Devices(ctx)
	if err != nil {
		return nil, err
	}

	for _, device := range devices {
		if device.Id == id {
			return device, nil
		}
	}

	return nil, building.ErrDeviceNotFound.Err()
}

// SetDeviceConfig updates the bridge with the new config options for the light or sensor.
func (b *HueBridge) SetDeviceConfig(ctx context.Context, dev *pb.Device, config *pb.DeviceConfig) error {
	if strings.Contains(dev.Address, lightAddrPrefix) {
		id := addrToLight(dev.Address)
		args := convertLightDiffToArgs(dev, config)

		err := b.bridge.SetLight(fmt.Sprintf("%d", id), &args)
		if err != nil {
			return err
		} else if len(args.Errors()) > 0 {
			return ErrHueResponseError
		}
		return nil
	} else if strings.Contains(dev.Address, sensorAddrPrefix) {
		id := addrToSensor(dev.Address)
		args := convertSensorDiffToArgs(dev, config)

		err := b.bridge.SetSensor(fmt.Sprintf("%d", id), &args)
		if err != nil {
			return err
		} else if len(args.Errors()) > 0 {
			return ErrHueResponseError
		}
		return nil
	}

	return ErrHueAddressInvalid
}
// SetDeviceState updates the bridge with the new state options for the light (sensors aren't supported).
func (b *HueBridge) SetDeviceState(ctx context.Context, dev *pb.Device, state *pb.DeviceState) error {
	if strings.Contains(dev.Address, lightAddrPrefix) {
		id := addrToLight(dev.Address)
		args, err := convertLightStateDiffToArgs(dev, state)
		if err != nil {
			return err
		}

		err = b.bridge.SetLightState(fmt.Sprintf("%d", id), &args)
		if err != nil {
			return err
		} else if len(args.Errors()) > 0 {
			return ErrHueResponseError
		}
		return nil
	}

	return ErrHueAddressInvalid
}
// AddDevice is not implemented yet.
func (b *HueBridge) AddDevice(ctx context.Context, id string) error {
	return building.ErrNotImplemented.Err()
}
// DeleteDevice is not implemented yet.
func (b *HueBridge) DeleteDevice(ctx context.Context, id string) error {
	return building.ErrNotImplemented.Err()
}

func convertLightToDevice(l hue.Light) *pb.Device {
	d := &pb.Device{}
	d.Reset()

	d.Id = l.UniqueID
	d.Address = lightAddrPrefix + l.ID
	d.IsActive = true

	d.Manufacturer = l.ManufacturerName
	d.ModelId = l.ModelID

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
			rgb.FromXY(xy, l.ModelID)
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

	d.Id = s.UniqueID
	d.Address = sensorAddrPrefix + s.ID

	d.Manufacturer = s.ManufacturerName
	d.ModelId = s.ModelID

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

func convertLightDiffToArgs(currDevice *pb.Device, newConfig *pb.DeviceConfig) hue.LightArg {
	var args hue.LightArg

	if currDevice.Config == nil && newConfig != nil {
		args.SetName(newConfig.Name)
	} else if currDevice.Config != nil && newConfig == nil {
		args.SetName("")
	} else if currDevice.Config.Name != newConfig.Name {
		args.SetName(newConfig.Name)
	}

	return args
}

func convertSensorDiffToArgs(currDevice *pb.Device, newConfig *pb.DeviceConfig) hue.SensorArg {
	var args hue.SensorArg

	if currDevice.Config == nil && newConfig != nil {
		args.SetName(newConfig.Name)
	} else if currDevice.Config != nil && newConfig == nil {
		args.SetName("")
	} else if currDevice.Config.Name != newConfig.Name {
		args.SetName(newConfig.Name)
	}

	return args
}

func convertLightStateDiffToArgs(currDevice *pb.Device, newState *pb.DeviceState) (args hue.LightStateArg, err error) {
	if currDevice.State == nil && newState != nil {
		if newState.Binary != nil {
			args.SetIsOn(newState.Binary.IsOn)
		}
		if newState.Range != nil {
			if currDevice.Range == nil {
				err = ErrDeviceLacksRangeCapability
				return
			} else if newState.Range.Value > currDevice.Range.Maximum || newState.Range.Value < currDevice.Range.Minimum {
				err = ErrDeviceRangeLimitExceeded
				return
			} else {
				args.SetBrightness(uint8(newState.Range.Value))
			}
		}
		if newState.ColorRgb != nil {
			colour := hue.RGB{
				Red: uint8(newState.ColorRgb.Red),
				Green: uint8(newState.ColorRgb.Green),
				Blue: uint8(newState.ColorRgb.Blue)}
			args.SetRGB(colour, currDevice.ModelId)
		}
	} else if currDevice.State != nil && newState == nil {
		args.SetIsOn(false)
	} else if currDevice.State != nil && newState != nil {
		if currDevice.State.Binary != nil && newState.Binary != nil && currDevice.State.Binary.IsOn != newState.Binary.IsOn {
			args.SetIsOn(newState.Binary.IsOn)
		}
		if newState.Range != nil {
			if currDevice.Range == nil {
				err = ErrDeviceLacksRangeCapability
				return
			} else if newState.Range.Value > currDevice.Range.Maximum || newState.Range.Value < currDevice.Range.Minimum {
				err = ErrDeviceRangeLimitExceeded
				return
			} else if currDevice.State.Range.Value != newState.Range.Value {
				args.SetBrightness(uint8(newState.Range.Value))
			}
		}
		if newState.ColorRgb != nil {
			if currDevice.State.ColorRgb == nil ||
				(currDevice.State.ColorRgb.Red != newState.ColorRgb.Red || currDevice.State.ColorRgb.Green != newState.ColorRgb.Green || currDevice.State.ColorRgb.Blue != newState.ColorRgb.Blue) {
				colour := hue.RGB{
					Red: uint8(newState.ColorRgb.Red),
					Green: uint8(newState.ColorRgb.Green),
					Blue: uint8(newState.ColorRgb.Blue)}
				args.SetRGB(colour, currDevice.ModelId)
			}
		}
	}

	return
}
