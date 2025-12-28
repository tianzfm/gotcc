package register

import (
	"sync"
)

type registryCenter struct {
	mux        sync.RWMutex
	components map[string]component.TCCComponent
}

func newRegistryCenter() *registryCenter {
	return &registryCenter{
		components: make(map[string]component.TCCComponent),
	}
}
