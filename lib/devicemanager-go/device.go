package devicemanager

import (
	"errors"

	"log"

	"faltung.ca/jvs/lib/proto-go"
	"golang.org/x/net/context"
	"google.golang.org/grpc/peer"
)

func (m *Manager) DevicesByBridgeId(bridgeId string) (devices []proto.Device, err error) {
	if bImpl, ok := m.bridges[bridgeId]; ok {
		devices = bImpl.devices
		return
	} else {
		err = errors.New("Bridge not found")
		return
	}
}

func (m *Manager) AddDeviceByBridgeId(bridgeId string, device proto.Device) (err error) {
	if bImpl, ok := m.bridges[bridgeId]; ok {
		bImpl.devices = append(bImpl.devices, device)

		update := proto.WatchDevicesResponse{
			Action: proto.WatchDevicesResponse_ADDED,
			Device: &device,
		}

		// Send a notification that the devices collection has changed
		m.deviceWatchers.Broadcast(update)

		return
	} else {
		err = errors.New("Bridge not found")
		return
	}
}

func (m *Manager) UpdateDeviceByBridgeId(bridgeId string, updatedDevice proto.Device) (err error) {
	if bImpl, ok := m.bridges[bridgeId]; ok {
		for idx, device := range bImpl.devices {
			if updatedDevice.Id != device.Id {
				continue
			}

			bImpl.devices[idx] = updatedDevice

			update := proto.WatchDevicesResponse{
				Action: proto.WatchDevicesResponse_CHANGED,
				Device: &updatedDevice,
			}

			m.deviceWatchers.Broadcast(update)
		}
		return
	} else {
		err = errors.New("Bridge not found")
		return
	}
}

func (m *Manager) RemoveDeviceByBridgeId(bridgeId string, deviceId string) (err error) {
	if bImpl, ok := m.bridges[bridgeId]; ok {
		for idx, device := range bImpl.devices {
			if deviceId != device.Id {
				continue
			}

			bImpl.devices = append(bImpl.devices[:idx], bImpl.devices[idx+1:]...)

			var device proto.Device
			device.Id = deviceId

			update := proto.WatchDevicesResponse{
				Action: proto.WatchDevicesResponse_CHANGED,
				Device: &device,
			}

			m.deviceWatchers.Broadcast(update)
		}
		return
	} else {
		err = errors.New("Bridge not found")
		return
	}
}

// Below are the gRPC function implementations

func (m *Manager) GetDevices(ctx context.Context, req *proto.GetDevicesRequest) (*proto.GetDevicesResponse, error) {
	resp := &proto.GetDevicesResponse{}

	for _, bImpl := range m.bridges {
		if !bImpl.b.IsActive() {
			continue
		}

		for idx := range bImpl.devices {
			resp.Devices = append(resp.Devices, &bImpl.devices[idx])
		}
	}

	return resp, nil
}

func (m *Manager) GetDevice(ctx context.Context, req *proto.GetDeviceRequest) (*proto.GetDeviceResponse, error) {
	for _, bImpl := range m.bridges {
		if !bImpl.b.IsActive() {
			continue
		}

		for _, device := range bImpl.devices {
			if device.Address != req.Address {
				continue
			}

			resp := &proto.GetDeviceResponse{}

			resp.Device = &device

			return resp, nil
		}
	}

	return nil, errors.New("Device is not available")
}

func (m *Manager) WatchDevices(req *proto.WatchDevicesRequest, stream proto.DeviceManager_WatchDevicesServer) error {
	peer, isOk := peer.FromContext(stream.Context())

	addr := "unknown"
	if isOk {
		addr = peer.Addr.String()
	}

	watcher := &deviceWatcher{
		updates: make(chan proto.WatchDevicesResponse),
		client:  addr,
	}

	m.deviceWatchers.Add(watcher)
	defer func() {
		log.Printf("Calling deviceWatchers.remove for %s\n", watcher.client)
		m.deviceWatchers.Remove(watcher)
	}()

	for {
		update := <-watcher.updates
		if err := stream.Send(&update); err != nil {
			log.Printf("Error sending to addr %s: %s\n", addr, err)
			return err
		}
	}

	return nil
}

func (m *Manager) SetDeviceState(ctx context.Context, req *proto.SetDeviceStateRequest) (*proto.SetDeviceStateResponse, error) {
	for _, bImpl := range m.bridges {
		if !bImpl.b.IsActive() {
			continue
		}

		for dIdx, device := range bImpl.devices {
			if device.Address != req.Address {
				continue
			}

			// TODO: check version

			resp := &proto.SetDeviceStateResponse{}

			err := bImpl.b.SetDeviceState(device, req.State)

			if err != nil {
				resp.Error = err.Error()
			} else {
				bImpl.devices[dIdx].State = req.State

				resp.Device = &bImpl.devices[dIdx]
			}

			return resp, nil
		}
	}

	return nil, errors.New("Device is not available")
}

func (m *Manager) SetDeviceConfig(ctx context.Context, req *proto.SetDeviceConfigRequest) (*proto.SetDeviceConfigResponse, error) {
	for _, bImpl := range m.bridges {
		if !bImpl.b.IsActive() {
			continue
		}

		for dIdx, device := range bImpl.devices {
			if device.Address != req.Address {
				continue
			}

			// TODO: check version

			resp := &proto.SetDeviceConfigResponse{}

			err := bImpl.b.SetDeviceConfig(device, req.Config)

			if err != nil {
				resp.Error = err.Error()
			} else {
				bImpl.devices[dIdx].Config = req.Config

				resp.Device = &bImpl.devices[dIdx]
			}

			return resp, nil
		}
	}

	return nil, errors.New("Device is not available")
}
