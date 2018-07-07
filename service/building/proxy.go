package building

import (
	"context"
	"errors"
	"log"

	"github.com/rmrobinson/jvs/service/building/pb"
	"google.golang.org/grpc"
	"fmt"
)

var (
	// ErrUnableToSetupProxy is returned if there was an error setting up the proxy bridge
	ErrUnableToSetupProxy = errors.New("unable to set up proxy bridge")
)

// ProxyBridge is a bridge implementation that proxies requests to a specified bridge service.
type ProxyBridge struct {
	// TODO: use a logger
	bn     Notifier
	remote pb.BridgeManagerClient

	state *pb.Bridge
}

// SetupNewProxyBridge creates a new proxy bridge using the specified bridge configuration.
func SetupNewProxyBridge(config *pb.BridgeConfig, notifier Notifier) (*ProxyBridge, error) {
	if config.Address.Ip == nil {
		return nil, ErrBridgeConfigInvalid.Err()
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	connStr := fmt.Sprintf("%s:%d", config.Address.Ip.Host, config.Address.Ip.Port)
	conn, err := grpc.Dial(connStr, opts...)

	if err != nil {
		log.Printf("Error initializing proxy connection to %s: %s\n", connStr, err.Error())
		return nil, ErrUnableToSetupProxy
	}

	log.Printf("Connected to %s\n", connStr)
	client := pb.NewBridgeManagerClient(conn)

	return NewProxyBridge(notifier, config.Id, client), nil
}

// NewProxyBridge creates a bridge implementation from a supplied bridge client.
func NewProxyBridge(notifier Notifier, id string, client pb.BridgeManagerClient) *ProxyBridge {
	ret := &ProxyBridge{
		bn:     notifier,
		remote: client,
		state: &pb.Bridge{
			Id: id,
		},
	}

	// Start watching for updates. This will populate the initial state;
	// we just ensure that the state ID is properly set for this to take effect.
	go ret.stateMonitor()

	return ret
}

func (b *ProxyBridge) ID() string {
	return b.state.Id
}

func (b *ProxyBridge) ModelID() string {
	return b.state.ModelId
}
func (b *ProxyBridge) ModelName() string {
	return b.state.ModelName
}
func (b *ProxyBridge) ModelDescription() string {
	return b.state.ModelDescription
}
func (b *ProxyBridge) Manufacturer() string {
	return b.state.Manufacturer
}
func (b *ProxyBridge) IconURLs() []string {
	return b.state.IconUrl
}
func (b *ProxyBridge) Name() string {
	return b.state.Config.Name
}
func (b *ProxyBridge) SetName(name string) error {
	return errors.New("not implemented")
}

func (b *ProxyBridge) Devices() ([]*pb.Device, error) {
	devices := []*pb.Device{}
	return devices, nil
}


func (b *ProxyBridge) stateMonitor() {
	stream, err := b.remote.WatchBridges(context.Background(), &pb.WatchBridgesRequest{})

	if err != nil {
		return
	}

	log.Printf("Waiting for updates about bridge ID %s\n", b.state.Id)

	for {
		if update, err := stream.Recv(); err == nil {
			// Filter out updates we don't care about
			if update.Bridge.Id != b.state.Id {
				log.Printf("Received update for bridge ID %s, ignoring\n", update.Bridge.Id)
				continue
			}

			log.Printf("Received update %+v\n", update.Bridge)
			b.state = update.Bridge
			b.bn.BridgeUpdated(b.state)
		} else {
			return
		}
	}
}
