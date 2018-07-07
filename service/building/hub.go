package building

import (
	"errors"
	"log"
	"reflect"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/rmrobinson/jvs/service/building/pb"
	"golang.org/x/net/context"
)

var (
	ErrBridgeAlreadyRegistered = errors.New("bridge already registered")
	ErrBridgeNotRegistered     = errors.New("bridge not registered")
	ErrDeviceAlreadyRegistered = errors.New("device already registered")
	ErrDeviceNotRegistered     = errors.New("device not registered")
	ErrNilArgument             = errors.New("nil argument")
)

// Bridge is an interface to a set of capabilities a device bridge must support.
type Bridge interface {
	Bridge(context.Context) (*pb.Bridge, error)
	SetBridgeConfig(context.Context, *pb.BridgeConfig) error
	SetBridgeState(context.Context, *pb.BridgeState) error

	SearchForAvailableDevices(context.Context) error
	AvailableDevices(context.Context) ([]*pb.Device, error)
	Devices(context.Context) ([]*pb.Device, error)
	Device(context.Context, string) (*pb.Device, error)

	SetDeviceConfig(context.Context, string, *pb.DeviceConfig) error
	SetDeviceState(context.Context, string, *pb.DeviceState) error
	AddDevice(context.Context, string) error
	DeleteDevice(context.Context, string) error

	//Pair(string) error
}

// AsyncBridge is an interface to a bridge that is able to detect changes and alert on them.
type AsyncBridge interface {
	Bridge

	SetNotifier(Notifier)
}

// Notifier is an interface used to signal changes outwards when there is a change to the specified bridge or device.
type Notifier interface {
	BridgeUpdated(bridge *pb.Bridge) error

	DeviceAdded(bridgeID string, device *pb.Device) error
	DeviceUpdated(bridgeID string, device *pb.Device) error
	DeviceRemoved(bridgeID string, device *pb.Device) error
}

// Hub contains the required logic to operate on a collection of bridges.
type Hub struct {
	bridgesLock sync.RWMutex
	bridges     map[string]*bridgeInstance

	bw bridgeWatchers
	dw deviceWatchers
}

// NewHub sets up a new bridge manager
func NewHub() *Hub {
	return &Hub{
		bridges: map[string]*bridgeInstance{},
		bw: bridgeWatchers{
			watchers: map[*bridgeWatcher]bool{},
		},
		dw: deviceWatchers{
			watchers: map[*deviceWatcher]bool{},
		},
	}
}

func (h *Hub) Bridges() []*pb.Bridge {
	h.bridgesLock.RLock()
	defer h.bridgesLock.RUnlock()

	var ret []*pb.Bridge
	for _, b := range h.bridges {
		ret = append(ret, b.bridge)
	}
	return ret
}

func (h *Hub) Devices() []*pb.Device {
	h.bridgesLock.RLock()
	defer h.bridgesLock.RUnlock()

	var ret []*pb.Device
	for _, b := range h.bridges {
		for _, d := range b.devices {
			ret = append(ret, d)
		}
	}
	return ret
}

func (h *Hub) SetDeviceConfig(ctx context.Context, id string, config *pb.DeviceConfig) (*pb.Device, error) {
	h.bridgesLock.RLock()
	defer h.bridgesLock.RUnlock()

	var bi *bridgeInstance
	for _, b := range h.bridges {
		if _, ok := b.devices[id]; ok {
			bi = b
			break
		}
	}

	if bi == nil {
		return nil, ErrDeviceNotRegistered
	}

	err := bi.bridgeHandle.SetDeviceConfig(ctx, id, config)

	if err != nil {
		return nil, err
	}

	// Propagate this update.
	// Since we don't have the full device we need to clone our current and update the state.
	d := proto.Clone(bi.devices[id]).(*pb.Device)
	d.Config = config
	h.DeviceUpdated(bi.bridgeID, d)

	return d, nil
}

func (h *Hub) SetDeviceState(ctx context.Context, id string, state *pb.DeviceState) (*pb.Device, error) {
	h.bridgesLock.RLock()
	defer h.bridgesLock.RUnlock()

	var bi *bridgeInstance
	for _, b := range h.bridges {
		if _, ok := b.devices[id]; ok {
			bi = b
			break
		}
	}

	if bi == nil {
		return nil, ErrDeviceNotRegistered
	}

	err := bi.bridgeHandle.SetDeviceState(ctx, id, state)
	if err != nil {
		return nil, err
	}

	// Propagate this update.
	// Since we don't have the full device we need to clone our current and update the state.
	d := proto.Clone(bi.devices[id]).(*pb.Device)
	d.State = state
	h.DeviceUpdated(bi.bridgeID, d)

	return d, nil
}

// AddBridge adds a pre-configured bridge into the collection of managed bridges.
// This will signal outwards that this bridge collection has been updated.
func (h *Hub) AddBridge(b Bridge, refreshRate time.Duration) error {
	bi, err := h.addBridgeInstance(b)
	if err != nil {
		return err
	}

	// We have a floor on the refresh rate for performance considerations.
	if refreshRate < time.Second {
		refreshRate = time.Second
	}

	go bi.monitor(refreshRate)

	return nil
}

func (h *Hub) AddAsyncBridge(b AsyncBridge) error {
	if _, err := h.addBridgeInstance(b); err != nil {
		return err
	}

	// We pass the hub back to the async bridge so it is able to notify the hub when things change.
	b.SetNotifier(h)
	return nil
}

func (h *Hub) addBridgeInstance(b Bridge) (*bridgeInstance, error) {
	startB, err := b.Bridge(context.Background())
	if err != nil {
		return nil, err
	}

	bi := newBridgeInstance(b, startB, h)

	h.bridgesLock.Lock()
	if _, ok := h.bridges[startB.Id]; ok {
		return nil, ErrBridgeAlreadyRegistered
	}
	h.bridges[startB.Id] = bi
	h.sendBridgeUpdate(pb.BridgeUpdate_ADDED, bi.bridge)
	h.bridgesLock.Unlock()

	bi.refresh()

	return bi, nil
}

func (h *Hub) RemoveBridge(id string) error {
	h.bridgesLock.Lock()
	h.bridgesLock.Unlock()

	if _, ok := h.bridges[id]; !ok {
		return ErrBridgeNotRegistered
	}

	bi := h.bridges[id]
	bi.cancelRefresh <- true

	h.sendBridgeUpdate(pb.BridgeUpdate_REMOVED, bi.bridge)
	delete(h.bridges, id)

	return nil
}

func (h *Hub) BridgeUpdated(bridge *pb.Bridge) error {
	if bridge == nil {
		return ErrNilArgument
	}

	h.bridgesLock.Lock()
	defer h.bridgesLock.Unlock()

	if _, ok := h.bridges[bridge.Id]; ok {
		nb := proto.Clone(bridge).(*pb.Bridge)
		h.sendBridgeUpdate(pb.BridgeUpdate_CHANGED, nb)
		h.bridges[bridge.Id].bridge = nb
	} else {
		return ErrBridgeNotRegistered
	}

	return nil
}

func (h *Hub) DeviceAdded(bridgeID string, device *pb.Device) error {
	h.bridgesLock.RLock()
	defer h.bridgesLock.RUnlock()

	if bridge, ok := h.bridges[bridgeID]; ok {
		bridge.lock.Lock()
		defer bridge.lock.Unlock()

		if _, ok := bridge.devices[device.Id]; ok {
			return ErrDeviceAlreadyRegistered
		}

		nd := proto.Clone(device).(*pb.Device)
		h.sendDeviceUpdate(pb.DeviceUpdate_ADDED, nd)
		h.bridges[bridgeID].devices[device.Id] = nd
	} else {
		return ErrBridgeNotRegistered
	}

	return nil
}

func (h *Hub) DeviceUpdated(bridgeID string, device *pb.Device) error {
	h.bridgesLock.RLock()
	defer h.bridgesLock.RUnlock()

	if bridge, ok := h.bridges[bridgeID]; ok {
		bridge.lock.Lock()
		defer bridge.lock.Unlock()

		if currDevice, ok := bridge.devices[device.Id]; ok {
			if !reflect.DeepEqual(currDevice, device) {
				nd := proto.Clone(device).(*pb.Device)
				h.sendDeviceUpdate(pb.DeviceUpdate_CHANGED, nd)
				h.bridges[bridgeID].devices[device.Id] = nd
			}
		} else {
			return ErrDeviceNotRegistered
		}
	} else {
		return ErrBridgeNotRegistered
	}

	return nil
}

func (h *Hub) DeviceRemoved(bridgeID string, device *pb.Device) error {
	h.bridgesLock.RLock()
	defer h.bridgesLock.RUnlock()

	if bridge, ok := h.bridges[bridgeID]; ok {
		bridge.lock.Lock()
		defer bridge.lock.Unlock()

		if _, ok := bridge.devices[device.Id]; !ok {
			return ErrDeviceNotRegistered
		}

		h.sendDeviceUpdate(pb.DeviceUpdate_REMOVED, device)
		delete(h.bridges[bridgeID].devices, device.Id)
	} else {
		return ErrBridgeNotRegistered
	}

	return nil
}

// sendDeviceUpdate is the internal function that takes a notification and propagates it to all registered watchers.
func (h *Hub) sendDeviceUpdate(action pb.DeviceUpdate_Action, device *pb.Device) {
	h.dw.Lock()
	defer h.dw.Unlock()

	log.Printf("Device changed: %+v\n", device)

	for watcher, active := range h.dw.watchers {
		if !active {
			continue
		}

		// We perform this in a separate goroutine in case the watcher has not yet finished processing
		// a previously-received message (if, for example, the remote side is timing out).
		go func() {
			watcher.updates <- &pb.DeviceUpdate{
				Action: action,
				Device: device,
			}
		}()
	}
}

// sendBridgeUpdate is the internal function that takes a notification and propagates it to all registered watchers.
func (h *Hub) sendBridgeUpdate(action pb.BridgeUpdate_Action, bridge *pb.Bridge) {
	h.bw.Lock()
	defer h.bw.Unlock()

	log.Printf("Bridge changed: %+v\n", bridge)

	for watcher, active := range h.bw.watchers {
		if !active {
			continue
		}

		// We perform this in a separate goroutine in case the watcher has not yet finished processing
		// a previously-received message (if, for example, the remote side is timing out).
		go func() {
			watcher.updates <- &pb.BridgeUpdate{
				Action: action,
				Bridge: bridge,
			}
		}()
	}
}
