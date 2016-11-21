package devicemanager

import (
	"errors"
	"log"

	"faltung.ca/jvs/lib/proto-go"
	"golang.org/x/net/context"
	"google.golang.org/grpc/peer"
)

type bridgeImpl struct {
	b     Bridge
	bData proto.Bridge

	devices []proto.Device
}

func (m *Manager) BridgeData(id string) (bData proto.Bridge, err error) {
	if bridge, ok := m.bridges[id]; ok {
		bData = bridge.bData
		return
	} else {
		err = errors.New("Bridge not found")
		return
	}
}

func (m *Manager) AddBridge(b Bridge) (err error) {
	if _, ok := m.bridges[b.Id()]; ok {
		return errors.New("Bridge implementation already present")
	}

	var bImpl bridgeImpl

	bImpl.b = b

	bImpl.bData, err = b.BridgeData()

	if err != nil {
		return
	}

	bImpl.devices, err = b.Devices()

	if err != nil {
		return
	}

	m.bridges[b.Id()] = bImpl

	update := proto.WatchBridgesResponse{
		Action: proto.WatchBridgesResponse_ADDED,
		Bridge: &bImpl.bData,
	}

	// Send a notification that the bridges collection has changed
	m.bridgeWatchers.Broadcast(update)

	// Send a notification for each device on the bridge
	for _, device := range bImpl.devices {
		log.Printf("Adding device to bridge %s (%s)\n", b.Id(), device.String())

		update := proto.WatchDevicesResponse{
			Action: proto.WatchDevicesResponse_ADDED,
			Device: &device,
		}

		m.deviceWatchers.Broadcast(update)
	}

	log.Printf("Added bridge %s\n", b.Id())

	return nil
}

func (m *Manager) UpdateBridge(b Bridge, bData proto.Bridge) {
	if bImpl, ok := m.bridges[b.Id()]; ok {
		bImpl.bData = bData

		update := proto.WatchBridgesResponse{
			Action: proto.WatchBridgesResponse_CHANGED,
			Bridge: &bData,
		}

		// Send a notification that the bridges collection has changed
		m.bridgeWatchers.Broadcast(update)

		log.Printf("Updated bridge %s\n", b.Id())
	}
}

func (m *Manager) RemoveBridge(id string) {
	if bImpl, ok := m.bridges[id]; ok {
		bImpl.b.Disable()

		update := proto.WatchBridgesResponse{
			Action: proto.WatchBridgesResponse_REMOVED,
			Bridge: &proto.Bridge{Id: id},
		}

		// Send a notification that the bridges collection has changed
		m.bridgeWatchers.Broadcast(update)

		log.Printf("Removed bridge %s\n", id)
	}
}

// Below are the gRPC function implementations

func (m *Manager) GetBridges(ctx context.Context, req *proto.GetBridgesRequest) (*proto.GetBridgesResponse, error) {
	resp := &proto.GetBridgesResponse{}

	for id, bImpl := range m.bridges {
		if !bImpl.b.IsActive() {
			continue
		}

		bData := bImpl.bData
		bData.Id = id

		resp.Bridges = append(resp.Bridges, &bData)
	}

	return resp, nil
}

func (m *Manager) WatchBridges(req *proto.WatchBridgesRequest, stream proto.BridgeManager_WatchBridgesServer) error {
	peer, isOk := peer.FromContext(stream.Context())

	addr := "unknown"
	if isOk {
		addr = peer.Addr.String()
	}

	watcher := &bridgeWatcher{
		updates: make(chan proto.WatchBridgesResponse),
		client:  addr,
	}

	m.bridgeWatchers.Add(watcher)
	defer m.bridgeWatchers.Remove(watcher)

	for {
		update := <-watcher.updates
		if err := stream.Send(&update); err != nil {
			return err
		}
	}

	return nil
}
