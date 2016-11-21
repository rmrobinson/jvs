package devicemanager

import (
	"net"

	"google.golang.org/grpc"

	"faltung.ca/jvs/lib/proto-go"
)

// This interface contains the functions which are expected to be implemented by any protocol bridge.
// It is possible that the implementation returns an error, that is ok and expected.
type Bridge interface {
	// Retrieve data about the bridge.
	Id() string
	IsActive() bool

	BridgeData() (proto.Bridge, error)

	SetState(*proto.BridgeState) error
	SetConfig(*proto.BridgeConfig) error

	Pair(string) error
	Disable()

	SearchForNewDevices() error
	CreateDevice(*proto.Device) error

	NewDevices() ([]proto.Device, error)
	Devices() ([]proto.Device, error)

	Device(string) (proto.Device, error)
	SetDeviceConfig(proto.Device, *proto.DeviceConfig) error
	SetDeviceState(proto.Device, *proto.DeviceState) error
	DeleteDevice(string) error
}

type Manager struct {
	bridges map[string]bridgeImpl

	bridgeWatchers bridgeWatchers
	deviceWatchers deviceWatchers
}

func NewManager() *Manager {
	m := &Manager{
		bridges:        make(map[string]bridgeImpl),
		bridgeWatchers: newBridgeWatchers(),
		deviceWatchers: newDeviceWatchers(),
	}

	return m
}

func (m *Manager) Run(l net.Listener) {
	// Run the gRPC listener.
	// This is assumed to never return.
	grpcServer := grpc.NewServer()

	proto.RegisterBridgeManagerServer(grpcServer, m)
	proto.RegisterDeviceManagerServer(grpcServer, m)

	grpcServer.Serve(l)
}
