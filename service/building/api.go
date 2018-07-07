package building

import (
	"errors"
	"log"

	"github.com/rmrobinson/jvs/service/building/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

var (
	ErrNotImplemented        = status.New(codes.Unimplemented, "not implemented")
	ErrBridgeTypeUndefined   = status.New(codes.InvalidArgument, "bridge type must be specified to create")
	ErrBridgeTypeUnsupported = status.New(codes.InvalidArgument, "bridge type cannot be created manually")
	ErrBridgeConfigInvalid   = status.New(codes.InvalidArgument, "bridge config is not supported for requested type")
	ErrDeviceNotFound        = status.New(codes.NotFound, "device not found")
)

type API struct {
	hub *Hub
}

func NewAPI(hub *Hub) *API {
	return &API{
		hub: hub,
	}
}

func (a *API) GetBridges(ctx context.Context, req *pb.GetBridgesRequest) (*pb.GetBridgesResponse, error) {
	resp := &pb.GetBridgesResponse{}
	for _, bridge := range a.hub.Bridges() {
		resp.Bridges = append(resp.Bridges, bridge)
	}

	return resp, nil
}

func (a *API) UpdateBridge(context.Context, *pb.UpdateBridgeRequest) (*pb.UpdateBridgeResponse, error) {
	return nil, errors.New("not implemented")
}

func (a *API) WatchBridges(req *pb.WatchBridgesRequest, stream pb.BridgeManager_WatchBridgesServer) error {
	peer, isOk := peer.FromContext(stream.Context())

	addr := "unknown"
	if isOk {
		addr = peer.Addr.String()
	}

	log.Printf("WatchBridges request from %s\n", addr)

	watcher := &bridgeWatcher{
		updates:  make(chan *pb.BridgeUpdate),
		peerAddr: addr,
	}

	a.hub.bw.add(watcher)
	defer a.hub.bw.remove(watcher)

	// Send all of the currently active bridges to start.
	for _, impl := range a.hub.bridges {
		update := &pb.BridgeUpdate{
			Action: pb.BridgeUpdate_ADDED,
			Bridge: impl.bridge,
		}

		log.Printf("Sending update %+v to %s\n", update, addr)
		if err := stream.Send(update); err != nil {
			return err
		}
	}
	// TODO: the above is subject to a race condition where the add is processed after we've added the watcher
	// but before we get the range of bridges so we duplicate data.
	// This shouldn't cause issues on the client (they should be tolerant to this) but let's fix this anyways.

	// Now we wait for updates
	for {
		update := <-watcher.updates

		log.Printf("Sending update %+v to %s\n", update, addr)
		if err := stream.Send(update); err != nil {
			return err
		}
	}

	return nil
}

func (a *API) GetDevices(context.Context, *pb.GetDevicesRequest) (*pb.GetDevicesResponse, error) {
	resp := &pb.GetDevicesResponse{}
	for _, device := range a.hub.Devices() {
		resp.Devices = append(resp.Devices, device)
	}

	return resp, nil
}
func (a *API) GetDevice(ctx context.Context, req *pb.GetDeviceRequest) (*pb.GetDeviceResponse, error) {
	for _, device := range a.hub.Devices() {
		if device.Id == req.Id {
			return &pb.GetDeviceResponse{
				Device: device,
			}, nil
		}
	}

	return nil, ErrDeviceNotFound.Err()
}
func (a *API) SetDeviceConfig(ctx context.Context, req *pb.SetDeviceConfigRequest) (*pb.SetDeviceConfigResponse, error) {
	resp, err := a.hub.SetDeviceConfig(ctx, req.Id, req.Config)

	return &pb.SetDeviceConfigResponse{
		Device: resp,
	}, err
}
func (a *API) SetDeviceState(ctx context.Context, req *pb.SetDeviceStateRequest) (*pb.SetDeviceStateResponse, error) {
	resp, err := a.hub.SetDeviceState(ctx, req.Id, req.State)

	return &pb.SetDeviceStateResponse{
		Device: resp,
	}, err
}
func (a *API) WatchDevices(req *pb.WatchDevicesRequest, stream pb.DeviceManager_WatchDevicesServer) error {
	peer, isOk := peer.FromContext(stream.Context())

	addr := "unknown"
	if isOk {
		addr = peer.Addr.String()
	}

	log.Printf("WatchDevices request from %s\n", addr)

	watcher := &deviceWatcher{
		updates:  make(chan *pb.DeviceUpdate),
		peerAddr: addr,
	}

	a.hub.dw.add(watcher)
	defer a.hub.dw.remove(watcher)

	// Send all of the currently active bridges to start.
	for _, impl := range a.hub.Devices() {
		update := &pb.DeviceUpdate{
			Action: pb.DeviceUpdate_ADDED,
			Device: impl,
		}

		log.Printf("Sending update %+v to %s\n", update, addr)
		if err := stream.Send(update); err != nil {
			return err
		}
	}
	// TODO: the above is subject to a race condition where the add is processed after we've added the watcher
	// but before we get the range of devices so we duplicate data.
	// This shouldn't cause issues on the client (they should be tolerant to this) but let's fix this anyways.

	// Now we wait for updates
	for {
		update := <-watcher.updates

		log.Printf("Sending update %+v to %s\n", update, addr)
		if err := stream.Send(update); err != nil {
			return err
		}
	}

	return nil
}
