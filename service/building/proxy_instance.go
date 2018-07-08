package building

import (
	"context"

	"github.com/rmrobinson/jvs/service/building/pb"
	"google.golang.org/grpc"
)

type proxyInstance struct {
	bridgeID string

	remoteBridge pb.BridgeManagerClient
	remoteDevice pb.DeviceManagerClient
	notifier     Notifier
}

func newProxyInstance(conn *grpc.ClientConn, id string) *proxyInstance {
	return &proxyInstance{
		bridgeID:     id,
		remoteBridge: pb.NewBridgeManagerClient(conn),
		remoteDevice: pb.NewDeviceManagerClient(conn),
	}
}

// SetNotifier saves the notifier to use for the proxy instance.
func (pi *proxyInstance) SetNotifier(notifier Notifier) {
	pi.notifier = notifier
}

// Bridge queries the linked peer to retrieve information on the requested bridge.
func (pi *proxyInstance) Bridge(ctx context.Context) (*pb.Bridge, error) {
	resp, err := pi.remoteBridge.GetBridges(ctx, &pb.GetBridgesRequest{})
	if err != nil {
		return nil, err
	}

	for _, bridge := range resp.Bridges {
		if bridge.Id == pi.bridgeID {
			return bridge, nil
		}
	}

	return nil, ErrBridgeNotRegistered
}

// SetBridgeConfig updates the linked peer with the config of the specified bridge.
func (pi *proxyInstance) SetBridgeConfig(ctx context.Context, config *pb.BridgeConfig) error {
	req := &pb.SetBridgeConfigRequest{
		Id:     pi.bridgeID,
		Config: config,
	}

	_, err := pi.remoteBridge.SetBridgeConfig(ctx, req)
	return err
}

// SetBridgeState updates the linked peer with the state of the specified bridge.
func (pi *proxyInstance) SetBridgeState(context.Context, *pb.BridgeState) error {
	return ErrNotImplemented.Err()
}

// SearchForAvailableDevices requests that the linked peer begin searching for new devices.
func (pi *proxyInstance) SearchForAvailableDevices(context.Context) error {
	return ErrNotImplemented.Err()
}

// AvailableDevices queries the linked peer for the list of available but not yet added devices.
func (pi *proxyInstance) AvailableDevices(ctx context.Context) ([]*pb.Device, error) {
	resp, err := pi.remoteDevice.GetAvailableDevices(ctx, &pb.GetDevicesRequest{})
	if err != nil {
		return nil, err
	}

	return resp.Devices, nil
}

// Devices queries the linked peer for its collection of devices.
func (pi *proxyInstance) Devices(ctx context.Context) ([]*pb.Device, error) {
	resp, err := pi.remoteDevice.GetDevices(ctx, &pb.GetDevicesRequest{BridgeId: pi.bridgeID})
	if err != nil {
		return nil, err
	}

	return resp.Devices, nil
}

// Device queries the linked peer for the specified device.
func (pi *proxyInstance) Device(ctx context.Context, id string) (*pb.Device, error) {
	req := &pb.GetDeviceRequest{
		Id: id,
	}
	resp, err := pi.remoteDevice.GetDevice(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Device, nil
}

// SetDeviceConfig updates the linked peer with the config for the specified device.
func (pi *proxyInstance) SetDeviceConfig(ctx context.Context, id string, config *pb.DeviceConfig) error {
	req := &pb.SetDeviceConfigRequest{
		Id:     id,
		Config: config,
	}
	_, err := pi.remoteDevice.SetDeviceConfig(ctx, req)
	return err
}

// SetDeviceState updates the linked peer with the state for the specified device.
func (pi *proxyInstance) SetDeviceState(ctx context.Context, id string, state *pb.DeviceState) error {
	req := &pb.SetDeviceStateRequest{
		Id:    id,
		State: state,
	}
	_, err := pi.remoteDevice.SetDeviceState(ctx, req)
	return err
}

// AddDevice requests that the linked peer add the specified device.
func (pi *proxyInstance) AddDevice(context.Context, string) error {
	return ErrNotImplemented.Err()
}

// DeleteDevice requests that the linked peer remove the specified device.
func (pi *proxyInstance) DeleteDevice(context.Context, string) error {
	return ErrNotImplemented.Err()
}
