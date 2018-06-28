package device

import (
	"github.com/google/uuid"
	"github.com/rmrobinson/jvs/service/device/pb"
)

// LoopbackBridge offers the standard bridge capabilities in a form that doesn't do anything.
type LoopbackBridge struct {
	// TODO: use a logger
	id string
}

// NewLoopbackBridge creates a new loopback bridge.
func NewLoopbackBridge() *LoopbackBridge {
	return &LoopbackBridge{
		id: uuid.New().String(),
	}
}

func (b *LoopbackBridge) ID() string {
	return b.id
}
func (b *LoopbackBridge) ModelID() string {
	return "LoopbackID"
}
func (b *LoopbackBridge) ModelName() string {
	return "Loopback"
}
func (b *LoopbackBridge) ModelDescription() string {
	return "Loopback bridge"
}
func (b *LoopbackBridge) Manufacturer() string {
	return "Faltung Systems"
}
func (b *LoopbackBridge) IconURLs() []string {
	return []string{}
}
func (b *LoopbackBridge) Name() string {
	return "Loopback"
}
func (b *LoopbackBridge) SetName(name string) error {
	return nil
}

func (b *LoopbackBridge) Devices() ([]pb.Device, error) {
	devices := []pb.Device{}
	return devices, nil
}
