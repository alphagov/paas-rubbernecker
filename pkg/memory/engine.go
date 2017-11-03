package memory

import "fmt"

// SetupEngine should compose the storage.
func SetupEngine() *Engine {
	return &Engine{
		storage: map[string]interface{}{},
	}
}

// Engine module configuration.
type Engine struct {
	storage map[string]interface{}
}

// Get a specific value from a memory store.
func (e *Engine) Get(key string) (interface{}, error) {
	if value, ok := e.storage[key]; ok {
		return value, nil
	}

	return nil, fmt.Errorf("memory: key not found in storage")
}

// Put specific value into a memory store.
func (e *Engine) Put(key string, value interface{}) error {
	e.storage[key] = value

	return nil
}
