package bridge

import (
	"errors"
	"log"

	"github.com/rmrobinson/jvs/service/device/pb"
	monopamp "github.com/rmrobinson/monoprice-amp-go"
	"github.com/tarm/serial"
)

const (
	portBaudRate = 9600
)

var (
	ErrUnableToSetupMonopAmp = errors.New("unable to set up monoprice amp")
)

// MonopAmpBridge is an implementation of a bridge for the Hue service.
type MonopAmpBridge struct {
	bn  BridgeNotifier
	amp *monopamp.SerialAmplifier
}

// SetupNewMonopAmpBridge takes a bridge config and returns the appropriate monoprice amp bridge if possible.
func SetupNewMonopAmpBridge(config *pb.BridgeConfig, notifier BridgeNotifier) (*MonopAmpBridge, error) {
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

	return NewMonopAmpBridge(notifier, amp), err
}

// NewMonopAmpBridge takes a previously set up MonopAmp handle and exposes it as a MonopAmp bridge.
func NewMonopAmpBridge(notifier BridgeNotifier, amp *monopamp.SerialAmplifier) *MonopAmpBridge {
	ret := &MonopAmpBridge{
		bn:  notifier,
		amp: amp,
	}

	return ret
}

func (b *MonopAmpBridge) ID() string {
	return "2"
}

func (b *MonopAmpBridge) ModelID() string {
	return "10761"
}
func (b *MonopAmpBridge) ModelName() string {
	return "Monoprice Amp"
}
func (b *MonopAmpBridge) ModelDescription() string {
	return "6 Zone Home Audio Multizone Controller"
}
func (b *MonopAmpBridge) Manufacturer() string {
	return "Monoprice"
}
func (b *MonopAmpBridge) IconURLs() []string {
	return []string{}
}
func (b *MonopAmpBridge) Name() string {
	return ""
}
func (b *MonopAmpBridge) SetName(name string) error {
	return nil
}

func (b *MonopAmpBridge) Devices() ([]*pb.Device, error) {
	devices := []*pb.Device{}
	return devices, nil
}
