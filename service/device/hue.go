package device

import (
	"context"
	"time"

	hue "github.com/rmrobinson/hue-go"
)

const (
	stateRefreshInterval = 60 * time.Second
)

// HuePersister is an interface to persisting bridge profiles.
type HuePersister interface {
	Profile(ctx context.Context, bridgeID string) (string, error)
	SaveProfile(ctx context.Context, bridgeID string, username string) error
}

// HueBridge is an implementation of a bridge for the Hue service.
type HueBridge struct {
	bn bridgeNotifier
	b  *hue.Bridge

	state               hue.BridgeDescription
	lastStateRefreshErr error
	stopStateRefresh    chan struct{}
}

// NewHueBridge takes a previously set up Hue handle and exposes it as a Hue bridge.
func NewHueBridge(notifier bridgeNotifier, bridge *hue.Bridge) *HueBridge {
	ret := &HueBridge{
		bn:               notifier,
		b:                bridge,
		stopStateRefresh: make(chan struct{}),
	}

	go ret.stateRefresher()

	return ret
}

func (h *HueBridge) Close() {
	close(h.stopStateRefresh)
}

func (h *HueBridge) ID() string {
	return h.state.Device.SerialNumber
}

func (h *HueBridge) ModelID() string {
	return h.state.Device.ModelNumber
}
func (h *HueBridge) ModelName() string {
	return h.state.Device.ModelName
}
func (h *HueBridge) ModelDescription() string {
	return h.state.Device.ModelDescription
}
func (h *HueBridge) Manufacturer() string {
	return h.state.Device.Manufacturer
}
func (h *HueBridge) IconURLs() []string {
	ret := []string{}

	for _, icon := range h.state.Device.Icons {
		ret = append(ret, h.state.UrlBase+"/"+icon.FileName)
	}

	return ret
}
func (h *HueBridge) Name() string {
	return h.state.Device.FriendlyName
}
func (h *HueBridge) SetName(name string) error {
	c := &hue.ConfigArg{}
	c.SetName(name)
	return h.b.SetConfig(c)
}

func (h *HueBridge) stateRefresher() {
	ticker := time.NewTicker(stateRefreshInterval)

	for {
		select {
		case <-ticker.C:
			state, err := h.b.Description()

			if err != nil {
				h.lastStateRefreshErr = err
				continue
			}

			if h.state.UrlBase != state.UrlBase {
				h.state = state
				h.bn.bridgeUpdated(h)
			}
		case <-h.stopStateRefresh:
			ticker.Stop()
			return
		}
	}
}
