package devicemanager

import (
	"log"
	"sync"

	"faltung.ca/jvs/lib/proto-go"
)

// This interface contains functions to be called when the device changes state
type deviceWatcher struct {
	updates chan proto.WatchDevicesResponse
	client  string
}

type deviceWatchers struct {
	sync.Mutex
	watchers map[*deviceWatcher]bool
}

func newDeviceWatchers() deviceWatchers {
	dw := deviceWatchers{
		watchers: make(map[*deviceWatcher]bool),
	}

	return dw
}

func (dw *deviceWatchers) Add(w *deviceWatcher) {
	dw.Lock()
	defer dw.Unlock()

	if dw.watchers[w] {
		return
	}

	dw.watchers[w] = true
}

func (dw *deviceWatchers) Remove(w *deviceWatcher) {
	dw.Lock()
	defer dw.Unlock()

	if !dw.watchers[w] {
		log.Printf("Cannot remove watcher %s as it no longer is active\n", w.client)
		return
	}

	log.Printf("Removed watcher %s\n", w.client)
	dw.watchers[w] = false
}

func (dw *deviceWatchers) Broadcast(update proto.WatchDevicesResponse) {
	dw.Lock()
	defer dw.Unlock()

	for watcher, active := range dw.watchers {
		if !active {
			continue
		}

		log.Printf("Broadcasting to %s\n", watcher.client)

		go func() {
			watcher.updates <- update
		}()
	}
}
