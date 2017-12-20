package gates

import (
	"github.com/calavera/go-flipper/actor"
	"github.com/calavera/go-flipper/feature"
)

// GroupFunc is a function type to check if a feature
// is enabled for a group or not.
type GroupFunc func(actor.Actor) bool

// GroupGate is a gate that checks if an actor is in a group or not.
// It uses registered GroupFunc functions for those checks.
type GroupGate struct {
	value Set
}

var registry = make(map[string]GroupFunc)

// NewGroupGate initializes a GroupGate with a set of group names.
func NewGroupGate(set Set) GroupGate {
	return GroupGate{set}
}

// Key returns the GateKey for an GroupGate gate.
func (GroupGate) Key() GateKey {
	return GroupGateKey
}

// IsOpen check if the gate is open for an feature and an actor.
// It uses the registered GroupFunc to know if the gate is open or not.
func (g GroupGate) IsOpen(f feature.Feature, a actor.Actor) bool {
	for name := range g.value {
		if f, ok := registry[name]; ok && f(a) {
			return true
		}
	}

	return false
}

// SetValue returns the set of groups for which the gate is open.
// This satisfies the SetGateType interface.
func (g GroupGate) SetValue() Set {
	return g.value
}

// RegisterGroup associates group names with functions.
func RegisterGroup(name string, f GroupFunc) {
	registry[name] = f
}
