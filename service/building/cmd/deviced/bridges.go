package main

import (
	"errors"
	"log"

	br "github.com/rmrobinson/bottlerocket-go"
	"github.com/rmrobinson/jvs/service/building"
	"github.com/rmrobinson/jvs/service/building/bridge"
	"github.com/rmrobinson/jvs/service/building/pb"
	monopamp "github.com/rmrobinson/monoprice-amp-go"
	"github.com/tarm/serial"
)

const (
	portBaudRate = 9600
)

var (
	// ErrUnableToSetupBottlerocket is returned if the supplied bridge configuration fails to properly initialize bottlerocket.
	ErrUnableToSetupBottlerocket = errors.New("unable to set up bottlerocket")
	ErrUnableToSetupMonopAmp     = errors.New("unable to set up monoprice amp")
)

func setupBottlerocket(config *pb.BridgeConfig) (building.Bridge, error) {
	if config.Address.Usb == nil {
		return nil, ErrBridgeConfigInvalid.Err()
	}

	br := &br.Bottlerocket{}
	err := br.Open(config.Address.Usb.Path)

	if err != nil {
		log.Printf("Error initializing bottlerocket: %s\n", err.Error())
		return nil, ErrUnableToSetupBottlerocket
	}

	return bridge.NewBottlerocketBridge(br), nil
}

func setupMonopriceAmp(config *pb.BridgeConfig) (building.Bridge, error) {
	if config.Address.Usb == nil {
		return nil, ErrBridgeConfigInvalid.Err()
	}

	c := &serial.Config{
		Name: config.Address.Usb.Path,
		Baud: portBaudRate,
	}
	s, err := serial.OpenPort(c)

	if err != nil {
		log.Printf("Error initializing serial port: %s\n", err.Error())
		return nil, ErrUnableToSetupMonopAmp
	}

	amp, err := monopamp.NewSerialAmplifier(s)

	return bridge.NewMonopAmpBridge(amp), err
}
