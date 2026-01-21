package backend

import (
	"fmt"
	"sync"
)

// Constructor is a function that creates a new backend instance.
type Constructor func(cfg Config) (Backend, error)

var (
	registry = make(map[string]Constructor)
	mu       sync.RWMutex
)

// Register adds a backend constructor to the registry.
// This should be called in the init() function of each backend package.
func Register(name string, constructor Constructor) {
	mu.Lock()
	defer mu.Unlock()
	registry[name] = constructor
}

// New creates a new backend instance based on the config type.
func New(cfg Config) (Backend, error) {
	mu.RLock()
	constructor, ok := registry[cfg.Type]
	mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("unknown backend type: %s", cfg.Type)
	}

	return constructor(cfg)
}

// Available returns a list of registered backend types.
func Available() []string {
	mu.RLock()
	defer mu.RUnlock()

	types := make([]string, 0, len(registry))
	for t := range registry {
		types = append(types, t)
	}
	return types
}
