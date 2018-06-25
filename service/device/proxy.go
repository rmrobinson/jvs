package device

import (
	"context"
	"errors"
	"log"

	"github.com/rmrobinson/jvs/service/device/pb"
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
	bn bridgeNotifier
	c  pb.BridgeManagerClient

	state *pb.Bridge
}

// SetupNewProxyBridge creates a new proxy bridge using the specified bridge configuration.
func SetupNewProxyBridge(config *pb.BridgeConfig, notifier bridgeNotifier) (*ProxyBridge, error) {
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
func NewProxyBridge(notifier bridgeNotifier, id string, client pb.BridgeManagerClient) *ProxyBridge {
	ret := &ProxyBridge{
		bn: notifier,
		c:  client,
		state: &pb.Bridge{
			Id: id,
		},
	}

	// Start watching for updates. This will populate the initial state;
	// we just ensure that the state ID is properly set for this to take effect.
	go ret.stateMonitor()

	return ret
}

func (p *ProxyBridge) ID() string {
	return p.state.Id
}

func (p *ProxyBridge) ModelID() string {
	return p.state.ModelId
}
func (p *ProxyBridge) ModelName() string {
	return p.state.ModelName
}
func (p *ProxyBridge) ModelDescription() string {
	return p.state.ModelDescription
}
func (p *ProxyBridge) Manufacturer() string {
	return p.state.Manufacturer
}
func (p *ProxyBridge) IconURLs() []string {
	return p.state.IconUrl
}
func (p *ProxyBridge) Name() string {
	return p.state.Config.Name
}
func (p *ProxyBridge) SetName(name string) error {
	return errors.New("not implemented")
}

func (p *ProxyBridge) stateMonitor() {
	stream, err := p.c.WatchBridges(context.Background(), &pb.WatchBridgesRequest{})

	if err != nil {
		return
	}

	log.Printf("Waiting for updates about bridge ID %s\n", p.state.Id)

	for {
		if update, err := stream.Recv(); err == nil {
			// Filter out updates we don't care about
			if update.Bridge.Id != p.state.Id {
				log.Printf("Received update for bridge ID %s, ignoring\n", update.Bridge.Id)
				continue
			}

			log.Printf("Received update %+v\n", update.Bridge)
			p.state = update.Bridge
			p.bn.bridgeUpdated(p)
		} else {
			return
		}
	}
}
