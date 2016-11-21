package main

import (
	"errors"
	"fmt"

	"faltung.ca/jvs/lib/proto-go"
	"github.com/rmrobinson/bottlerocket-go"
)

type Bridge struct {
	id string

	name string

	isActive bool

	br bottlerocket_go.Bottlerocket

	db storage

	devices []proto.Device
}

func (b *Bridge) Setup(devicePath string, dbPath string) error {
	err := b.db.Open(dbPath)

	if err != nil {
		return err
	}

	bData, err := b.db.Bridges()

	if err != nil {
		return err
	}

	if len(bData) < 1 {
		b.id = "1"
		b.name = "X10 Firecracker"

		var bd BridgeData
		bd.Id = b.id
		bd.Name = b.name

		err = b.db.SetBridge(bd)

		if err != nil {
			return err
		}
	} else {
		b.id = bData[0].Id
		b.name = bData[0].Name
	}

	b.isActive = true

	dData, err := b.db.Devices()

	if err != nil {
		return err
	}

	if len(dData) < 1 {
		b.devices = []proto.Device{}

		for _, house := range "ABCDEFGHIJKLMNOP" {
			for i := 1; i < 17; i++ {
				addr := fmt.Sprintf("%c%d", house, i)

				dd := DeviceData{}
				dd.Id = "X10-" + addr
				dd.Address = addr
				dd.Name = addr
				dd.IsActive = false
				dd.IsOn = false

				err = b.db.SetDevice(dd)

				if err != nil {
					return err
				}

				device := proto.Device{}
				device.Id = dd.Id
				device.Address = dd.Address

				device.ModelDescription = "Generic X10 device"
				device.ModelName = "Generic X10 device"
				device.ModelId = "Unknown"
				device.Manufacturer = "Unknown"

				device.IsActive = dd.IsActive

				device.Config = &proto.DeviceConfig{}
				device.Config.Name = dd.Name
				device.Config.Description = "Generic X10 device"

				device.State = &proto.DeviceState{}
				device.State.IsReachable = true
				device.State.Binary = &proto.DeviceState_BinaryState{
					IsOn: dd.IsOn,
				}

				b.devices = append(b.devices, device)

			}
		}
	} else {
		b.devices = []proto.Device{}

		for _, dd := range dData {
			device := proto.Device{}
			device.Id = dd.Id
			device.Address = dd.Address

			device.ModelDescription = "Generic X10 device"
			device.ModelName = "Generic X10 device"
			device.ModelId = "Unknown"
			device.Manufacturer = "Unknown"

			device.IsActive = dd.IsActive

			device.Config = &proto.DeviceConfig{}
			device.Config.Name = dd.Name
			device.Config.Description = "Generic X10 device"

			device.State = &proto.DeviceState{}
			device.State.IsReachable = true
			device.State.Binary = &proto.DeviceState_BinaryState{
				IsOn: dd.IsOn,
			}

			b.devices = append(b.devices, device)
		}
	}

	return b.br.Open(devicePath)
}

func (b *Bridge) Shutdown() {
	b.br.Close()
}

func (b *Bridge) Id() string {
	return b.id
}

func (b *Bridge) IsActive() bool {
	return b.isActive
}

func (b *Bridge) Disable() {
	b.isActive = false
}

func (b *Bridge) BridgeData() (bData proto.Bridge, _ error) {
	bData.Config = &proto.BridgeConfig{}
	bData.State = &proto.BridgeState{}

	bData.Id = b.Id()
	bData.IsActive = b.IsActive()
	bData.ModelName = "Firecracker"
	bData.Manufacturer = "x10.com"
	bData.ModelId = "CM17A"
	bData.ModelDescription = "Serial-X10 bridge"

	bData.Config.Address = &proto.Address{
		Usb: &proto.Address_Usb{
			Path: b.br.Path(),
		},
	}

	bData.Config.Name = b.name
	bData.Config.Timezone = "UTC"

	bData.State.IsPaired = true
	bData.State.Version = &proto.BridgeState_Version{
		Api: "1.0.0",
		Sw:  "0.05b3",
	}

	return
}

func (b *Bridge) Pair(identifier string) (err error) {
	err = errors.New("Does not support pairing")
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

func (b *Bridge) Devices() (devices []proto.Device, _ error) {
	devices = b.devices
	return
}

func (b *Bridge) Device(address string) (proto.Device, error) {
	for _, device := range b.devices {
		if device.Address == address {
			return device, nil
		}
	}

	return proto.Device{}, errors.New("Invalid address specified")
}

func (b *Bridge) SetDeviceConfig(device proto.Device, deviceConfig *proto.DeviceConfig) error {
	return errors.New("Not implemented")
}

func (b *Bridge) SetDeviceState(device proto.Device, deviceState *proto.DeviceState) error {
	if deviceState.Binary == nil {
		return errors.New("Device state must specify a binary condition")
	}

	for _, d := range b.devices {
		if d.Address != device.Address {
			continue
		}

		var err error

		if deviceState.Binary.IsOn {
			err = b.br.SendCommand(device.Address, "ON")
		} else {
			err = b.br.SendCommand(device.Address, "OFF")
		}

		if err != nil {
			return err
		}

		var dd DeviceData
		dd.Id = device.Id
		dd.Address = device.Address
		dd.Name = device.Config.Name
		dd.IsActive = true
		dd.IsOn = deviceState.Binary.IsOn

		err = b.db.SetDevice(dd)

		return err
	}

	return errors.New("Invalid address specified")
}

func (b *Bridge) DeleteDevice(address string) (err error) {
	err = errors.New("Not implemented")
	return
}

func (b *Bridge) SearchForNewDevices() (err error) {
	err = errors.New("Not implemented")
	return
}

func (b *Bridge) CreateDevice(device *proto.Device) (err error) {
	err = errors.New("Not implemented")
	return
}

func (b *Bridge) NewDevices() (_ []proto.Device, err error) {
	err = errors.New("Not implemented")
	return
}
