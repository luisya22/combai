package swarmlet

import (
	"fmt"
	"sync"
)

type Memory interface {
	Get(key string) (any, error)
	Set(key string, value any) error
	Append(keys string, value any) error
}

type DummyMemory struct {
	store map[string]any
	mu    sync.RWMutex
}

func NewDummyMemory() *DummyMemory {
	return &DummyMemory{
		store: make(map[string]any),
	}
}

func (d *DummyMemory) Get(key string) (any, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	val, ok := d.store[key]
	if !ok {
		return nil, fmt.Errorf("key '%s' not found in memory", key)
	}
	return val, nil
}
func (d *DummyMemory) Set(key string, value any) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.store == nil {
		d.store = make(map[string]any)
	}
	d.store[key] = value
	return nil
}
func (d *DummyMemory) Append(key string, value any) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.store == nil {
		d.store = make(map[string]any)
	}
	if existing, ok := d.store[key].(string); ok {
		d.store[key] = existing + "\n" + fmt.Sprintf("%v", value)
	} else {
		d.store[key] = value
	}
	return nil
}
