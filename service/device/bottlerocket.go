package device

import (
	"context"
	"errors"
	"log"

	br "github.com/rmrobinson/bottlerocket-go"
	"github.com/rmrobinson/jvs/service/device/pb"
)

var (
	// ErrUnableToSetupBottlerocket is returned if the supplied bridge configuration fails to properly initialize bottlerocket.
	ErrUnableToSetupBottlerocket = errors.New("unable to set up bottlerocket")
)

// BottlerocketDB exposes an interface to allow the state of the bottlerocket network to be persisted.
// Because the X10 protocol doesn't support querying we need to maintain state across process restarts.
type BottlerocketPersister interface {
	ID(ctx context.Context) (string, error)
	Name(ctx context.Context) (string, error)
	SetName(ctx context.Context, name string) error

	Devices(ctx context.Context, bridgeID string) ([]*pb.Device, error)
	AvailableDevices(ctx context.Context, bridgeID string) ([]*pb.Device, error)

	//SaveDevice(ctx context.Context, device *pb.Device) error
}

// BottlerocketBridge offers the standard bridge capabilities over the Bottlerocket X10 USB/serial interface.
type BottlerocketBridge struct {
	// TODO: use a logger
	br *br.Bottlerocket

	persister BottlerocketPersister
}

// SetupNewBottlerocketBridge takes a bridge config and returns the appropriate bottlerocket bridge if possible.
func SetupNewBottlerocketBridge(config *pb.BridgeConfig) (*BottlerocketBridge, error) {
	if config.Address.Usb == nil {
		return nil, ErrBridgeConfigInvalid.Err()
	}

	br := &br.Bottlerocket{}
	err := br.Open(config.Address.Usb.Path)

	if err != nil {
		log.Printf("Error initializing bottlerocket: %s\n", err.Error())
		return nil, ErrUnableToSetupBottlerocket
	}


	return NewBottlerocketBridge(br), nil
}

// NewBottlerocketBridge takes a previously set up bottlerocket handle and exposes it as a bottlerocket bridge.
func NewBottlerocketBridge(bridge *br.Bottlerocket) *BottlerocketBridge {
	return &BottlerocketBridge{
		br: bridge,
	}
}

func (b *BottlerocketBridge) ID() string {
	return "1"
}
func (b *BottlerocketBridge) ModelID() string {
	return "CM17A"
}
func (b *BottlerocketBridge) ModelName() string {
	return "Firecracker"
}
func (b *BottlerocketBridge) ModelDescription() string {
	return "Serial-X10 bridge"
}
func (b *BottlerocketBridge) Manufacturer() string {
	return "x10.com"
}
func (b *BottlerocketBridge) IconURLs() []string {
	return []string{}
}
func (b *BottlerocketBridge) Name() string {
	return "X10 Firecracker"
}
func (b *BottlerocketBridge) SetName(name string) error {
	return nil
}

func (b *BottlerocketBridge) Devices() ([]pb.Device, error) {
	return []pb.Device{}, nil
}
