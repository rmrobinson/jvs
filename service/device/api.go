package device

import (
	"log"

	"github.com/rmrobinson/jvs/service/device/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

var (
	ErrBridgeTypeUndefined   = status.New(codes.InvalidArgument, "bridge type must be specified to create")
	ErrBridgeTypeUnsupported = status.New(codes.InvalidArgument, "bridge type cannot be created manually")
	ErrBridgeConfigInvalid   = status.New(codes.InvalidArgument, "bridge config is not supported for requested type")
)

type API struct {
	bm *BridgeManager
}

func (a *API) GetBridgeNotifier() bridgeNotifier {
	return a.bm
}

func NewAPI(bm *BridgeManager) *API {
	return &API{
		bm: bm,
	}
}

func (a *API) AddBridge(ctx context.Context, req *pb.AddBridgeRequest) (*pb.AddBridgeResponse, error) {
	var b Bridge
	var err error

	switch req.Type {
	case pb.BridgeType_Loopback:
		b = NewLoopbackBridge()
	case pb.BridgeType_Proxy:
		b, err = SetupNewProxyBridge(req.Config, a.bm)
	case pb.BridgeType_Bottlerocket:
		b, err = SetupNewBottlerocketBridge(req.Config)
	case pb.BridgeType_MonopriceAmp:
		b, err = SetupNewMonopAmpBridge(req.Config, a.bm)
	case pb.BridgeType_Hue:
		return nil, ErrBridgeTypeUnsupported.Err()
	case pb.BridgeType_Generic:
		return nil, ErrBridgeTypeUndefined.Err()
	default:
		return nil, ErrBridgeTypeUndefined.Err()
	}

	if err != nil {
		return nil, err
	}
	a.bm.AddBridge(b)

	return &pb.AddBridgeResponse{
		Bridge: bridgeToProto(b),
	}, nil
}

func (a *API) GetBridges(ctx context.Context, req *pb.GetBridgesRequest) (*pb.GetBridgesResponse, error) {
	resp := &pb.GetBridgesResponse{}
	for _, bridge := range a.bm.bridges {
		resp.Bridges = append(resp.Bridges, bridgeToProto(bridge))
	}

	return resp, nil
}

func (a *API) WatchBridges(req *pb.WatchBridgesRequest, stream pb.BridgeManager_WatchBridgesServer) error {
	peer, isOk := peer.FromContext(stream.Context())

	addr := "unknown"
	if isOk {
		addr = peer.Addr.String()
	}

	log.Printf("WatchBridges request from %s\n", addr)

	watcher := &bWatcher{
		updates:  make(chan *pb.BridgeUpdate),
		peerAddr: addr,
	}

	a.bm.watchers.add(watcher)
	defer a.bm.watchers.remove(watcher)

	// Send all of the currently active bridges to start.
	for _, bridge := range a.bm.bridges {
		update := &pb.BridgeUpdate{
			Action: pb.BridgeUpdate_ADDED,
			Bridge: bridgeToProto(bridge),
		}

		log.Printf("Sending update %+v to %s\n", update, addr)
		if err := stream.Send(update); err != nil {
			return err
		}
	}
	// TODO: the above is subject to a gross race condition where the add is processed after we've added the watcher
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

func bridgeToProto(b Bridge) *pb.Bridge {
	resp := &pb.Bridge{
		Id:               b.ID(),
		IsActive:         true,
		ModelId:          b.ModelID(),
		ModelName:        b.ModelName(),
		ModelDescription: b.ModelDescription(),
		Manufacturer:     b.Manufacturer(),
	}

	for _, url := range b.IconURLs() {
		resp.IconUrl = append(resp.IconUrl, url)
	}

	return resp
}
