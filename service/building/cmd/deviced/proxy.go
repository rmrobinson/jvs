package main

import (
	"fmt"
	"log"

	"github.com/rmrobinson/jvs/service/building"
	"github.com/rmrobinson/jvs/service/building/pb"
	"google.golang.org/grpc"
)

type proxyImpl struct {
	conn *grpc.ClientConn
	p *building.ProxyHub
}

func (b *proxyImpl) setup(config *pb.BridgeConfig, hub *building.Hub) error {
	if config.Address.Ip == nil {
		return ErrBridgeConfigInvalid
	}

	var err error
	addr := fmt.Sprintf("%s:%d", config.Address.Ip.Host, config.Address.Ip.Port)
	log.Printf("Proxying requests to %s\n", addr)

	// Setup the proxyImpl connection first
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	b.conn, err = grpc.Dial(addr, opts...)
	if err != nil {
		log.Printf("Error initializing proxyImpl connection to %s: %s\n", addr, err.Error())
		return err
	}

	b.p = building.NewProxyBridge(hub, b.conn)
	go b.p.Run()

	return nil
}

// Close cleans up any open resources.
func (b *proxyImpl) Close() error {
	var connErr error
	if b.conn != nil {
		connErr = b.conn.Close()
	}
	return connErr
}
