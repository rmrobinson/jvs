package main

import (
	"errors"
	"fmt"
	"strings"

	"faltung.ca/jvs/lib/proto-go"
	"github.com/rmrobinson/hue-go"
)

// This represents an instance of a Hue bridge.
// One bridge instance will be created for each profile specified in the configuration at startup.
// The UPnP discovery process will also create one instance for each profile detected not already running.
type Bridge struct {
	isActive bool

	b hue_go.Bridge
}

func (b *Bridge) Id() string {
	return b.b.Id()
}

func (b *Bridge) IsActive() bool {
	return b.isActive
}

func (b *Bridge) Disable() {
	b.isActive = false
}

func (b *Bridge) BridgeData() (bData proto.Bridge, err error) {
	bData.Config = &proto.BridgeConfig{}
	bData.State = &proto.BridgeState{}

	bDesc, err := b.b.Description()

	if err != nil {
		return
	}

	bData.Id = bDesc.Device.SerialNumber

	bData.IsActive = b.isActive

	bData.ModelId = bDesc.Device.ModelNumber
	bData.ModelName = bDesc.Device.ModelName
	bData.ModelDescription = bDesc.Device.ModelDescription
	bData.Manufacturer = bDesc.Device.Manufacturer

	for _, icon := range bDesc.Device.Icons {
		bData.IconUrl = append(bData.IconUrl, bDesc.UrlBase+icon.FileName)
	}

	bData.Config.Name = bDesc.Device.FriendlyName

	bConfig, err := b.b.Config()

	if err != nil {
		return
	}

	bData.ModelId = bConfig.ModelVersion

	bData.State.Zigbee = &proto.BridgeState_Zigbee{}
	bData.State.Zigbee.Channel = bConfig.ZigbeeChannel

	bData.State.IsPaired = true
	bData.State.Version = &proto.BridgeState_Version{}
	bData.State.Version.Api = bConfig.ApiVersion
	bData.State.Version.Sw = bConfig.SwVersion

	bData.Config.Name = bConfig.Name
	bData.Config.Timezone = bConfig.Timezone

	bData.Config.Address = &proto.Address{}
	bData.Config.Address.Ip = &proto.Address_Ip{}
	bData.Config.Address.Ip.Host = bConfig.IpAddress
	bData.Config.Address.Ip.Netmask = bConfig.SubnetMask
	bData.Config.Address.Ip.Gateway = bConfig.GatewayAddress
	bData.Config.Address.Ip.ViaDhcp = bConfig.IsDhcpAcquired

	return
}

func (b *Bridge) Pair(identifier string) (err error) {
	_, err = b.b.Pair(identifier)
	return
}

func (b *Bridge) SetConfig(config *proto.BridgeConfig) (err error) {
	err = errors.New("Not implemented")
	return
}

func (b *Bridge) SetState(state *proto.BridgeState) (err error) {
	err = errors.New("Not implemented")
	return
}

func (b *Bridge) Devices() (devices []proto.Device, err error) {
	lights, err := b.b.Lights()

	if err != nil {
		return
	}

	sensors, err := b.b.Sensors()

	if err != nil {
		return
	}

	for _, light := range lights {
		d := convertLightToDevice(light)
		devices = append(devices, d)
	}
	for _, sensor := range sensors {
		d := convertSensorToDevice(sensor)
		devices = append(devices, d)
	}

	return
}

func (b *Bridge) Device(address string) (device proto.Device, err error) {
	if strings.Contains(address, "light") {
		light, err := b.b.Light(address)

		if err == nil {
			device = convertLightToDevice(light)
		}
	} else if strings.Contains(address, "sensor") {
		sensor, err := b.b.Sensor(address)

		if err == nil {
			device = convertSensorToDevice(sensor)
		}
	} else {
		err = errors.New("Invalid address specified")
	}
	return
}

func (b *Bridge) SetDeviceConfig(device proto.Device, deviceConfig *proto.DeviceConfig) (err error) {
	addr := strings.Split(device.Address, "/")

	if len(addr) != 2 {
		err = errors.New("Invalid address specified")
		return
	}

	if addr[0] == "light" {
		args := convertLightDiffToArgs(device, deviceConfig)

		err = b.b.SetLight(addr[1], &args)

		if err != nil {
			return
		} else if len(args.Errors()) > 0 {
			for _, err := range args.Errors() {
				return errors.New(err.Description)
			}
		} else {
			// TODO: take the values which were changed and write them back to deviceConfig
			return
		}
	} else if addr[0] == "sensor" {
		args := convertSensorDiffToArgs(device, deviceConfig)

		err = b.b.SetSensor(addr[1], &args)

		if err != nil {
			return
		} else if len(args.Errors()) > 0 {
			for _, err := range args.Errors() {
				return errors.New(err.Description)
			}
		} else {
			// TODO: take the values which were changed and write them back to deviceConfig
			return
		}
	} else {
		err = errors.New("Invalid address specified")
	}

	return
}

func (b *Bridge) SetDeviceState(device proto.Device, deviceState *proto.DeviceState) (err error) {
	addr := strings.Split(device.Address, "/")

	if len(addr) != 2 {
		err = errors.New("Invalid address specified")
		return
	}

	if addr[0] == "light" {
		args, convErr := convertLightStateDiffToArgs(device, deviceState)

		if convErr != nil {
			err = convErr
			return
		}

		err = b.b.SetLightState(addr[1], &args)

		if err != nil {
			return
		} else if len(args.Errors()) > 0 {
			for _, argErr := range args.Errors() {
				return errors.New(argErr.Description)
			}
		} else {
			// TODO: take the values which were changed and write them back to deviceState
			return
		}
	} else {
		// TODO: allow the user to set the state of sensors (only supported by the CLIP Hue sensors currently).
		err = errors.New("Invalid address specified")
	}

	return
}

func (b *Bridge) DeleteDevice(address string) error {
	addr := strings.Split(address, "/")

	if len(addr) != 2 {
		return errors.New("Invalid address specified")
	}

	if addr[0] == "light" {
		return b.b.DeleteLight(addr[1])
	} else if addr[0] == "sensor" {
		return b.b.DeleteSensor(addr[1])
	} else {
		return errors.New("Invalid address specified")
	}
}

func (b *Bridge) SearchForNewDevices() (err error) {
	if err = b.b.SearchForNewLights(); err != nil {
		return
	}
	if err = b.b.SearchForNewSensors(); err != nil {
		return
	}

	return nil
}

func (b *Bridge) CreateDevice(device *proto.Device) (err error) {
	err = errors.New("Not implemented")
	return
}

func (b *Bridge) NewDevices() (devices []proto.Device, err error) {
	lights, err := b.b.NewLights()

	if err != nil {
		return
	}

	sensors, err := b.b.NewSensors()

	if err != nil {
		return
	}

	for _, light := range lights {
		var device proto.Device
		device.Config = &proto.DeviceConfig{}

		device.Address = fmt.Sprintf("light/%d", light.Id)
		device.Config.Name = light.Name

		devices = append(devices, device)
	}

	for _, sensor := range sensors {
		var device proto.Device
		device.Config = &proto.DeviceConfig{}

		device.Address = fmt.Sprintf("sensor/%d", sensor.Id)
		device.Config.Name = sensor.Name

		devices = append(devices, device)
	}

	return
}

func convertLightToDevice(l hue_go.Light) proto.Device {
	var d proto.Device
	d.Reset()

	d.Id = l.UniqueId
	d.Address = "light/" + l.Id
	d.IsActive = true

	d.Manufacturer = l.ManufacturerName
	d.ModelId = l.ModelId

	config := &proto.DeviceConfig{}
	d.Config = config

	config.Name = l.Name

	state := &proto.DeviceState{}
	d.State = state

	state.IsReachable = l.State.Reachable
	state.Version = l.SwVersion

	state.Binary = &proto.DeviceState_BinaryState{IsOn: l.State.On}

	if l.Model == "Dimmable light" ||
		l.Model == "Color light" ||
		l.Model == "Extended color light" {

		// We are hardcoded to only supporting a uint8 worth of brightness values
		d.Range = &proto.Device_RangeDevice{Minimum: 0, Maximum: 254}

		state.Range = &proto.DeviceState_RangeState{Value: int32(l.State.Brightness)}
	}

	if l.Model == "Color light" ||
		l.Model == "Extended color light" ||
		l.Model == "Color temperature light" {
		if l.State.ColorMode == "xy" {
			xy := hue_go.XY{X: l.State.XY[0], Y: l.State.XY[1]}

			rgb := hue_go.RGB{}
			rgb.FromXY(xy, l.ModelId)
			state.ColorRgb = &proto.DeviceState_RGBState{Red: int32(rgb.Red), Blue: int32(rgb.Blue), Green: int32(rgb.Green)}
		} else if l.State.ColorMode == "ct" {
			rgb := hue_go.RGB{}
			rgb.FromCT(l.State.ColorTemperature)
			state.ColorRgb = &proto.DeviceState_RGBState{Red: int32(rgb.Red), Blue: int32(rgb.Blue), Green: int32(rgb.Green)}
		} else if l.State.ColorMode == "hs" {
			hsb := hue_go.HSB{Hue: l.State.Hue, Saturation: l.State.Saturation, Brightness: l.State.Brightness}

			rgb := hue_go.RGB{}
			rgb.FromHSB(hsb)
			state.ColorRgb = &proto.DeviceState_RGBState{Red: int32(rgb.Red), Blue: int32(rgb.Blue), Green: int32(rgb.Green)}
		}

	}

	return d
}

func convertSensorToDevice(s hue_go.Sensor) proto.Device {
	var d proto.Device
	d.Reset()

	d.Id = s.UniqueId
	d.Address = "sensor/" + s.Id

	d.Manufacturer = s.ManufacturerName
	d.ModelId = s.ModelId

	config := &proto.DeviceConfig{}
	d.Config = config

	config.Name = s.Name

	state := &proto.DeviceState{}
	d.State = state

	state.IsReachable = s.Config.Reachable
	state.Version = s.SwVersion

	if s.Type == "ZGPSwitch" {
		button := &proto.DeviceState_ButtonState{}
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
		button := &proto.DeviceState_ButtonState{}

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
		d.Range = &proto.Device_RangeDevice{Minimum: 0, Maximum: s.SensitivityMax}
		state.Range = &proto.DeviceState_RangeState{Value: s.Sensitivity}

		state.Presence = &proto.DeviceState_PresenceState{IsPresent: s.State.Presence}
	} else if s.Type == "ZLLTemperature" {
		state.Temperature = &proto.DeviceState_TemperatureState{TemperatureCelsius: s.State.Temperature / 100}
	}

	return d
}

func convertLightDiffToArgs(currDevice proto.Device, newConfig *proto.DeviceConfig) hue_go.LightArg {
	var args hue_go.LightArg

	if currDevice.Config == nil && newConfig != nil {
		args.SetName(newConfig.Name)
	} else if currDevice.Config != nil && newConfig == nil {
		args.SetName("")
	} else if currDevice.Config.Name != newConfig.Name {
		args.SetName(newConfig.Name)
	}

	return args
}

func convertSensorDiffToArgs(currDevice proto.Device, newConfig *proto.DeviceConfig) hue_go.SensorArg {
	var args hue_go.SensorArg

	if currDevice.Config == nil && newConfig != nil {
		args.SetName(newConfig.Name)
	} else if currDevice.Config != nil && newConfig == nil {
		args.SetName("")
	} else if currDevice.Config.Name != newConfig.Name {
		args.SetName(newConfig.Name)
	}

	return args
}

func convertLightStateDiffToArgs(currDevice proto.Device, newState *proto.DeviceState) (args hue_go.LightStateArg, err error) {
	if currDevice.State == nil && newState != nil {
		if newState.Binary != nil {
			args.SetIsOn(newState.Binary.IsOn)
		}
		if newState.Range != nil {
			if currDevice.Range == nil {
				err = errors.New("Invalid argument supplied: device lacks range capabilities")
				return
			} else if newState.Range.Value > currDevice.Range.Maximum || newState.Range.Value < currDevice.Range.Minimum {
				err = errors.New("Invalid argument supplied: range value outside of allowed values")
				return
			} else {
				args.SetBrightness(uint8(newState.Range.Value))
			}
		}
		if newState.ColorRgb != nil {
			colour := hue_go.RGB{Red: uint8(newState.ColorRgb.Red), Green: uint8(newState.ColorRgb.Green), Blue: uint8(newState.ColorRgb.Blue)}
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
				err = errors.New("Invalid argument supplied: device lacks range capabilities")
				return
			} else if newState.Range.Value > currDevice.Range.Maximum || newState.Range.Value < currDevice.Range.Minimum {
				err = errors.New("Invalid argument supplied: range value outside of allowed values")
				return
			} else if currDevice.State.Range.Value != newState.Range.Value {
				args.SetBrightness(uint8(newState.Range.Value))
			}
		}
		if newState.ColorRgb != nil {
			if currDevice.State.ColorRgb == nil ||
				(currDevice.State.ColorRgb.Red != newState.ColorRgb.Red || currDevice.State.ColorRgb.Green != newState.ColorRgb.Green || currDevice.State.ColorRgb.Blue != newState.ColorRgb.Blue) {
				colour := hue_go.RGB{Red: uint8(newState.ColorRgb.Red), Green: uint8(newState.ColorRgb.Green), Blue: uint8(newState.ColorRgb.Blue)}
				args.SetRGB(colour, currDevice.ModelId)
			}
		}
	}

	return
}
