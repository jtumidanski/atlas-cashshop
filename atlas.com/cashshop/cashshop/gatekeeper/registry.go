package gatekeeper

import (
	"sync"
)

type registry struct {
	gatekeepers uint32
	lock        sync.RWMutex
}

var r *registry
var once sync.Once

func GetRegistry() *registry {
	once.Do(func() {
		r = &registry{
			gatekeepers: 0,
			lock:        sync.RWMutex{},
		}
	})
	return r
}

func (r *registry) Register(_ string) {
	r.lock.Lock()
	r.gatekeepers++
	r.lock.Unlock()
}

func (r *registry) Unregister(_ string) {
	r.lock.Lock()
	r.gatekeepers--
	r.lock.Unlock()
}

func (r *registry) Count() uint32 {
	return r.gatekeepers
}
