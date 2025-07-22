package swarmlet

import (
	"io"
	"sync"
)

type RunContext struct {
	RunID        string
	NodeInputs   map[string]string
	NodeOutputs  map[string]string
	NodeErrors   map[string]error
	StreamWriter io.Writer
	mu           sync.RWMutex
}

func NewRunContext(runID string, w io.Writer) *RunContext {
	return &RunContext{
		RunID:        runID,
		NodeInputs:   make(map[string]string),
		NodeOutputs:  make(map[string]string),
		NodeErrors:   make(map[string]error),
		StreamWriter: w,
	}
}

func (rc *RunContext) AddInput(key string, value string) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.NodeInputs[key] = value
}

func (rc *RunContext) GetInput(key string) (string, bool) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	val, ok := rc.NodeInputs[key]
	return val, ok
}

func (rc *RunContext) AddOutput(key string, value string) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.NodeOutputs[key] = value
}

func (rc *RunContext) GetOutput(key string) (string, bool) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	val, ok := rc.NodeOutputs[key]
	return val, ok
}

func (rc *RunContext) AddError(key string, err error) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.NodeErrors[key] = err
}

func (rc *RunContext) GetError(key string) (error, bool) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	err, ok := rc.NodeErrors[key]
	return err, ok
}
