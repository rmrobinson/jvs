package device

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
	bn  bridgeNotifier
	amp *monopamp.Amplifier
}

// SetupNewMonopAmpBridge takes a bridge config and returns the appropriate monoprice amp bridge if possible.
func SetupNewMonopAmpBridge(config *pb.BridgeConfig, notifier bridgeNotifier) (*MonopAmpBridge, error) {
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

	amp := monopamp.NewAmplifier(s)

	return NewMonopAmpBridge(notifier, amp), nil
}

// NewMonopAmpBridge takes a previously set up MonopAmp handle and exposes it as a MonopAmp bridge.
func NewMonopAmpBridge(notifier bridgeNotifier, amp *monopamp.Amplifier) *MonopAmpBridge {
	ret := &MonopAmpBridge{
		bn:  notifier,
		amp: amp,
	}

	return ret
}

func (h *MonopAmpBridge) ID() string {
	return "2"
}

func (h *MonopAmpBridge) ModelID() string {
	return "10761"
}
func (h *MonopAmpBridge) ModelName() string {
	return "Monoprice Amp"
}
func (h *MonopAmpBridge) ModelDescription() string {
	return "6 Zone Home Audio Multizone Controller"
}
func (h *MonopAmpBridge) Manufacturer() string {
	return "Monoprice"
}
func (h *MonopAmpBridge) IconURLs() []string {
	return []string{}
}
func (h *MonopAmpBridge) Name() string {
	return ""
}
func (h *MonopAmpBridge) SetName(name string) error {
	return nil
}
