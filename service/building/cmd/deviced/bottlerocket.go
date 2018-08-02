package main

import (
	"context"
	"log"

	br "github.com/rmrobinson/bottlerocket-go"
	"github.com/rmrobinson/jvs/service/building"
	"github.com/rmrobinson/jvs/service/building/bridge"
	"github.com/rmrobinson/jvs/service/building/pb"
)

type bottlerocketImpl struct {
	db *building.BridgeDB
	br *br.Bottlerocket

	bridge building.Bridge
}

func (b *bottlerocketImpl) setup(config *pb.BridgeConfig) error {
	if config.Address.Usb == nil {
		return ErrBridgeConfigInvalid
	}

	b.br = &br.Bottlerocket{}
	b.db = &building.BridgeDB{}

	err := b.br.Open(config.Address.Usb.Path)
	if err != nil {
		log.Printf("Error initializing bottlerocketImpl: %s\n", err.Error())
		return ErrUnableToSetupBottlerocket
	}

	err = b.db.Open(config.CachePath)
	if err != nil {
		b.br.Close()
		return err
	}

	brBridge := bridge.NewBottlerocketBridge(b.br, b.db)
	b.bridge = brBridge
	return brBridge.Setup(context.Background())
}

// Close cleans up any open resources.
func (b *bottlerocketImpl) Close() error {
	if b.db != nil {
		b.db.Close()
	}
	if b.br != nil {
		b.br.Close()
	}
	return nil
}
