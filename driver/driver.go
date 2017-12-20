package driver

import (
	"github.com/calavera/go-flipper/feature"
	"github.com/calavera/go-flipper/gates"
)

var registry = make(map[string]Driver)

// Driver defines how flipper gets information
// from a source. This source can be a database,
// an http endpoint or anything that implements
// this interface.
type Driver interface {
	Configure(config map[string]interface{}) error
	Enable(feature feature.Feature, gate gates.Gate) error
	Disable(feature feature.Feature, gate gates.Gate) error
	Get(feature feature.Feature, keys []gates.GateKey) ([]gates.Gate, error)
}

// Init stores an driver by name to be used
// by a client. This allows drivers to self
// register themselves on initialization
func Init(name string, a Driver) {
	registry[name] = a
}

// Get retrieves a driver from the registry by its name.
// It returns nil if the driver doesn't exist.
func Get(name string) Driver {
	return registry[name]
}
