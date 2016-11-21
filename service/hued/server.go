package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"faltung.ca/jvs/lib/devicemanager-go"
	"faltung.ca/jvs/lib/proto-go"
	"github.com/rmrobinson/hue-go"
)

type server struct {
	locator *hue_go.Locator

	pairings Storage

	manager *devicemanager.Manager
}

func newServer() *server {
	s := &server{
		manager: devicemanager.NewManager(),
		locator: hue_go.NewLocator(),
	}

	return s
}

func (s *server) run(configPath string, port int) error {
	s.pairings.Open(configPath)
	defer s.pairings.Close()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	defer lis.Close()

	if err != nil {
		return err
	}

	newBridges := make(chan hue_go.Bridge)

	go s.locator.Run(newBridges)
	go s.runNewBridgeHandler(newBridges)

	s.manager.Run(lis)

	return nil
}

func (s *server) runNewBridgeHandler(newBridges chan hue_go.Bridge) {
	for {
		hueBridge := <-newBridges

		p, err := s.pairings.Profile(hueBridge.Id())

		if err != nil {
			log.Printf("Unable to get pairing for ID: %s\n", err)
		} else {
			hueBridge.Username = p.Username
		}

		var bridge Bridge
		bridge.b = hueBridge
		bridge.isActive = true

		var dmBridge devicemanager.Bridge
		dmBridge = &bridge

		s.manager.AddBridge(dmBridge)

		go s.runBridgeMonitor(bridge)
	}
}

// Check the status of the bridge very minute
func (s *server) runBridgeMonitor(bridge Bridge) {
	ticker := time.NewTicker(60 * time.Second)
	quit := make(chan struct{})

	go s.runDeviceMonitor(bridge, quit)

	for {
		select {
		case <-ticker.C:
			b, err := s.manager.BridgeData(bridge.Id())

			if err != nil {
				log.Printf("Bridge %s not present; removing...\n", bridge.Id())
				close(quit)

				continue
			}

			var existingConfig proto.BridgeConfig
			var existingState proto.BridgeState

			existingConfig = *b.Config
			existingState = *b.State

			bData, err := bridge.BridgeData()

			if err != nil {
				log.Printf("Bridge %s config load error; removing...\n", bridge.Id())
				close(quit)

				// TODO: ensure this propagates to the sensor and light monitors
				continue
			}

			// Info will not change so doesn't need to be checked.
			if existingConfig.Name != bData.Config.Name ||
				existingConfig.Timezone != bData.Config.Timezone ||
				existingConfig.Address.Ip.Host != bData.Config.Address.Ip.Host ||
				existingConfig.Address.Ip.Netmask != bData.Config.Address.Ip.Netmask ||
				existingConfig.Address.Ip.Gateway != bData.Config.Address.Ip.Gateway ||
				existingConfig.Address.Ip.ViaDhcp != bData.Config.Address.Ip.ViaDhcp ||
				existingState.IsPaired != bData.State.IsPaired ||
				existingState.Version.Api != bData.State.Version.Api ||
				existingState.Version.Sw != bData.State.Version.Sw ||
				existingState.Zigbee.Channel != bData.State.Zigbee.Channel {

				var tmp devicemanager.Bridge
				tmp = &bridge

				s.manager.UpdateBridge(tmp, bData)
			}
		case <-quit:
			log.Fatal("quit channel triggered, no longer monitoring bridge %s\n", bridge.Id())
			ticker.Stop()
			return
		}
	}
}

// Check the status of the devices every second
func (s *server) runDeviceMonitor(bridge Bridge, quit chan struct{}) {
	ticker := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-ticker.C:
			existingDevices, err := s.manager.DevicesByBridgeId(bridge.Id())

			if err != nil {
				log.Printf("Devices for bridge %s not present; removing... (%s)\n", bridge.Id(), err)
				close(quit)

				continue
			}

			currentDevices, err := bridge.Devices()

			if err != nil {
				log.Printf("Bridge %s device load error; removing... (%s)\n", bridge.Id(), err)
				close(quit)

				continue
			}

			var existingDevicesMap map[string]proto.Device
			existingDevicesMap = make(map[string]proto.Device)

			for _, existingDevice := range existingDevices {
				existingDevicesMap[existingDevice.Id] = existingDevice
			}

			var currentDevicesMap map[string]proto.Device
			currentDevicesMap = make(map[string]proto.Device)

			for _, currentDevice := range currentDevices {
				currentDevicesMap[currentDevice.Id] = currentDevice
			}

			var addedDevices []proto.Device
			var removedDevices []proto.Device
			var changedDevices []proto.Device

			// Iterate over the current devices to find any a) ones which have changed, or b) new ones
			for _, currentDevice := range currentDevices {
				if existingDevice, ok := existingDevicesMap[currentDevice.Id]; ok {
					/*
						if currentDevice.Info.ModelId != existingDevice.Info.ModelId ||
							currentDevice.Info.ModelName != existingDevice.Info.ModelName ||
							currentDevice.Info.ModelDescription != existingDevice.Info.ModelDescription ||
							currentDevice.Info.Manufacturer != existingDevice.Info.Manufacturer ||
							currentDevice.Info.IsActive != existingDevice.Info.IsActive ||
							*currentDevice.Config != *existingDevice.Config {
					*/
					if currentDevice.String() != existingDevice.String() {
						// TODO: compare state structs
						changedDevices = append(changedDevices, currentDevice)
					}
				} else {
					addedDevices = append(addedDevices, currentDevice)
				}
			}

			// Iterate over the existing devices to find a) ones which no longer exist
			for _, existingDevice := range existingDevices {
				if _, ok := currentDevicesMap[existingDevice.Id]; !ok {
					removedDevices = append(removedDevices, existingDevice)
				}
			}

			// Send the relevant notifications to the different callers
			for _, device := range addedDevices {
				log.Printf("Adding device to bridge %s (%+v)\n", bridge.Id(), device)
				s.manager.AddDeviceByBridgeId(bridge.Id(), device)
			}

			for _, device := range changedDevices {
				log.Printf("Changed device on bridge %s (%+v)\n", bridge.Id(), device)
				s.manager.UpdateDeviceByBridgeId(bridge.Id(), device)
			}

			for _, device := range removedDevices {
				log.Printf("Removed device on bridge %s (%+v)\n", bridge.Id(), device)
				s.manager.RemoveDeviceByBridgeId(bridge.Id(), device.Id)
			}

		case <-quit:
			ticker.Stop()
			return
		}
	}
}
