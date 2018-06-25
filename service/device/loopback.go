package device

import (
	"github.com/google/uuid"
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

func (h *LoopbackBridge) ID() string {
	return h.id
}
func (h *LoopbackBridge) ModelID() string {
	return "LoopbackID"
}
func (h *LoopbackBridge) ModelName() string {
	return "Loopback"
}
func (h *LoopbackBridge) ModelDescription() string {
	return "Loopback bridge"
}
func (h *LoopbackBridge) Manufacturer() string {
	return "Faltung Systems"
}
func (h *LoopbackBridge) IconURLs() []string {
	return []string{}
}
func (h *LoopbackBridge) Name() string {
	return "Loopback"
}
func (h *LoopbackBridge) SetName(name string) error {
	return nil
}
