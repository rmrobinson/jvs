package main

import (
	"faltung.ca/jvs/lib/proto-go"
	"google.golang.org/grpc"
)

type conn struct {
	addr string
	conn *grpc.ClientConn

	bridgeClient proto.BridgeManagerClient
	deviceClient proto.DeviceManagerClient

	cancelBridgeWatcher func()
	cancelDeviceWatcher func()
}

func connect(addr string) (c conn, err error) {
	c = conn{
		addr:                addr,
		cancelBridgeWatcher: nil,
		cancelDeviceWatcher: nil,
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	c.conn, err = grpc.Dial(c.addr, opts...)

	if err != nil {
		return
	}

	c.bridgeClient = proto.NewBridgeManagerClient(c.conn)
	c.deviceClient = proto.NewDeviceManagerClient(c.conn)

	return
}

func (c *conn) close() {
	if c.cancelBridgeWatcher != nil {
		c.cancelBridgeWatcher()
		c.cancelBridgeWatcher = nil
	}

	if c.cancelDeviceWatcher != nil {
		c.cancelDeviceWatcher()
		c.cancelDeviceWatcher = nil
	}

	c.conn.Close()
	c.conn = nil
}
