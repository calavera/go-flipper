package memory

import (
	"fmt"

	"github.com/calavera/go-flipper/driver"
	"github.com/calavera/go-flipper/feature"
	"github.com/calavera/go-flipper/gates"
	"github.com/pkg/errors"
)

const keyFormat = "feature/%s/%s"

// Driver is a store driver that keeps features and gates in memory.
type Driver struct {
	store map[string]interface{}
}

// NewDriver initializes a new memory driver.
func NewDriver() *Driver {
	store := make(map[string]interface{})
	return &Driver{store}
}

// Configure configures the memory driver.
// This driver doesn't have any configuration, so this is a NOOP.
func (a *Driver) Configure(config map[string]interface{}) error {
	return nil
}

// Enable opens a feature for a give gate.
func (a *Driver) Enable(feature feature.Feature, gate gates.Gate) error {
	k := key(feature.Name, gate.Key())

	if g, ok := gate.(gates.IntGateType); ok {
		a.store[k] = g.IntValue()
	} else if _, ok := gate.(gates.BoolGateType); ok {
		a.store[k] = true
	} else if g, ok := gate.(gates.SetGateType); ok {
		var gs gates.Set
		if s, ok := a.store[k]; ok {
			gs, ok = s.(gates.Set)
			if !ok {
				return errors.Errorf("unexpected set value enabling feature: %v", s)
			}
		} else {
			gs = gates.Set{}
		}

		for k, v := range g.SetValue() {
			gs[k] = v
		}

		a.store[k] = gs
	} else {
		return errors.Errorf("unsupported data type: %v", gate.Key())
	}

	return nil
}

// Disable closes a feature for a given gate.
func (a *Driver) Disable(feature feature.Feature, gate gates.Gate) error {
	k := key(feature.Name, gate.Key())

	if g, ok := gate.(gates.IntGateType); ok {
		a.store[k] = g.IntValue()
	} else if _, ok := gate.(gates.BoolGateType); ok {
		delete(a.store, k)
	} else if g, ok := gate.(gates.SetGateType); ok {
		if s, ok := a.store[k]; ok {
			gs, ok := s.(gates.Set)
			if !ok {
				return errors.Errorf("unexpected set value disabling feature: %v", s)
			}
			for k := range g.SetValue() {
				delete(gs, k)
			}
			a.store[k] = gs
		}
	} else {
		return errors.Errorf("unsupported data type: %v", gate.Key())
	}

	return nil
}

// Get returns the enabled gates for a feature given a set of gate keys.
// Gates are skipped if they are not open for a feature.
func (a *Driver) Get(feature feature.Feature, keys []gates.GateKey) ([]gates.Gate, error) {
	var g []gates.Gate

	for _, t := range keys {
		k := key(feature.Name, t)
		v, ok := a.store[k]
		if !ok {
			continue
		}

		switch t {
		case gates.BoolGateKey:
			g = append(g, gates.NewBoolGate(ok))
		case gates.ActorGateKey:
			gs, ok := v.(gates.Set)
			if !ok {
				return nil, errors.Errorf("unexpected set value stored: %v", v)
			}
			g = append(g, gates.NewActorGate(gs))
		case gates.GroupGateKey:
			gs, ok := v.(gates.Set)
			if !ok {
				return nil, errors.Errorf("unexpected set value stored: %v", v)
			}
			g = append(g, gates.NewGroupGate(gs))
		case gates.PercentageOfActorsGateKey:
			gi, ok := v.(int)
			if !ok {
				return nil, errors.Errorf("unexpected int value: %v", v)
			}
			g = append(g, gates.NewPercentageOfActorsGate(gi))
		case gates.PercentageOfTimeGateKey:
			gi, ok := v.(int)
			if !ok {
				return nil, errors.Errorf("unexpected int value: %v", v)
			}
			g = append(g, gates.NewPercentageOfTimeGate(gi))
		default:
			return nil, errors.Errorf("unsupported gate: %v", t)
		}
	}

	return g, nil
}

func key(featureName string, gateKey gates.GateKey) string {
	return fmt.Sprintf(keyFormat, featureName, gateKey)
}

func init() {
	driver.Init("memory", NewDriver())
}
