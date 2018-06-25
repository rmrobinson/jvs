package device

import (
	"sync"

	"github.com/rmrobinson/jvs/service/device/pb"
)

// Bridge is an interface to a set of capabilities a device bridge must support.
type Bridge interface {

	//SearchForNewDevices() error
	//CreateDevice(*pb.Device) error

	//NewDevices() ([]pb.Device, error)
	//Devices() ([]pb.Device, error)

	//Device(string) (pb.Device, error)
	//SetDeviceConfig(pb.Device, *proto.DeviceConfig) error
	//SetDeviceState(pb.Device, *proto.DeviceState) error
	//DeleteDevice(string) error

	ID() string
	ModelID() string
	ModelName() string
	ModelDescription() string
	Manufacturer() string

	//IsActive() bool
	//Activate() error
	//Deactivate() error

	Name() string
	SetName(string) error

	//Address() pb.Address
	//SetAddress(*pb.Address) error

	//Pair(string) error

	IconURLs() []string
}

// bridgeNotifier is an interface used to signal changes outwards when there is a change to the specified bridge.
type bridgeNotifier interface {
	bridgeUpdated(bridge Bridge)
}

// BridgeManager contains the required logic to operate on a collection of bridges.
type BridgeManager struct {
	bridges []Bridge

	updates  chan *pb.BridgeUpdate
	watchers bWatchers
}

// NewBridgeManager sets up a new bridge manager
func NewBridgeManager() *BridgeManager {
	return &BridgeManager{
		updates: make(chan *pb.BridgeUpdate),
		watchers: bWatchers{
			watchers: map[*bWatcher]bool{},
		},
	}
}

// AddBridge adds a pre-configured bridge into the collection of managed bridges.
// This will signal outwards that this bridge collection has been updated.
func (bm *BridgeManager) AddBridge(b Bridge) {
	bm.bridges = append(bm.bridges, b)

	bm.notify(pb.BridgeUpdate_ADDED, bridgeToProto(b))
}

// bridgeUpdated should be used by a bridge implementation to signal outwards that something about its state has changed.
// This is kept package-scoped to prevent spurious updates coming from external packages.
func (bm *BridgeManager) bridgeUpdated(bridge Bridge) {
	bm.notify(pb.BridgeUpdate_CHANGED, bridgeToProto(bridge))
}

// notify is the internal function that takes a notification and propogates it to all registered watchers.
func (bm *BridgeManager) notify(action pb.BridgeUpdate_Action, bridge *pb.Bridge) {
	bm.watchers.Lock()
	defer bm.watchers.Unlock()

	for watcher, active := range bm.watchers.watchers {
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

// bWatchers is a synchronized map that gives us the ability to safely add and remove watchers while sending updates.
type bWatchers struct {
	sync.Mutex
	watchers map[*bWatcher]bool
}

func (bw *bWatchers) add(watcher *bWatcher) {
	bw.Lock()
	defer bw.Unlock()

	bw.watchers[watcher] = true
}

func (bw *bWatchers) remove(watcher *bWatcher) {
	bw.Lock()
	defer bw.Unlock()

	bw.watchers[watcher] = false
}

// bWatcher is a bridge watcher. A user should try to keep the updates channel empty but failure to read this
// will not block updates from being propagated to other watchers.
type bWatcher struct {
	updates  chan *pb.BridgeUpdate
	peerAddr string
}
