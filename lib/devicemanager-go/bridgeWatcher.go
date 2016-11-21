package devicemanager

import (
	"log"
	"sync"

	"faltung.ca/jvs/lib/proto-go"
)

// This interface contains functions to be called when the bridge changes state
type bridgeWatcher struct {
	updates chan proto.WatchBridgesResponse
	client  string
}

type bridgeWatchers struct {
	sync.Mutex
	watchers map[*bridgeWatcher]bool
}

func newBridgeWatchers() bridgeWatchers {
	bw := bridgeWatchers{
		watchers: make(map[*bridgeWatcher]bool),
	}

	return bw
}

func (bw *bridgeWatchers) Add(w *bridgeWatcher) {
	bw.Lock()
	defer bw.Unlock()

	if bw.watchers[w] {
		return
	}

	bw.watchers[w] = true
}

func (bw *bridgeWatchers) Remove(w *bridgeWatcher) {
	bw.Lock()
	defer bw.Unlock()

	if !bw.watchers[w] {
		return
	}

	bw.watchers[w] = false
}

func (bw *bridgeWatchers) Broadcast(update proto.WatchBridgesResponse) {
	bw.Lock()
	defer bw.Unlock()

	for watcher, active := range bw.watchers {
		if !active {
			continue
		}

		log.Printf("Broadcasting to %s\n", watcher.client)

		// TODO: clean up this handling of things; we'll be spwaning a ton of unnecessary go threads
		go func() {
			watcher.updates <- update
		}()
	}
}
