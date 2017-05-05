package deviceclient_go

import (
	"faltung.ca/jvs/lib/proto-go"
	"google.golang.org/grpc"
)

type Conn struct {
	Addr                string
	Conn                *grpc.ClientConn

	BridgeClient        proto.BridgeManagerClient
	DeviceClient        proto.DeviceManagerClient

	cancelBridgeWatcher func()
	cancelDeviceWatcher func()
}

func Connect(addr string) (c Conn, err error) {
	c = Conn{
		Addr:                addr,
		cancelBridgeWatcher: nil,
		cancelDeviceWatcher: nil,
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	c.Conn, err = grpc.Dial(c.Addr, opts...)

	if err != nil {
		return
	}

	c.BridgeClient = proto.NewBridgeManagerClient(c.Conn)
	c.DeviceClient = proto.NewDeviceManagerClient(c.Conn)

	return
}

func (c *Conn) Close() {
	if c.cancelBridgeWatcher != nil {
		c.cancelBridgeWatcher()
		c.cancelBridgeWatcher = nil
	}

	if c.cancelDeviceWatcher != nil {
		c.cancelDeviceWatcher()
		c.cancelDeviceWatcher = nil
	}

	c.Conn.Close()
	c.Conn = nil
}
