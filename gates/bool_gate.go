package gates

import (
	"github.com/calavera/go-flipper/actor"
	"github.com/calavera/go-flipper/feature"
)

// BoolGate is a gate that's open when its value is true.
type BoolGate struct {
	value bool
}

// NewBoolGate initializes a new BoolGate with a value.
func NewBoolGate(value bool) BoolGate {
	return BoolGate{value}
}

// Key returns the GateKey for a BoolGate gate.
func (BoolGate) Key() GateKey {
	return BoolGateKey
}

// IsOpen check if the gate is open for an feature and an actor.
// It returns the value of the gate.
func (g BoolGate) IsOpen(f feature.Feature, a actor.Actor) bool {
	return g.value
}

// BoolValue returns the boolean value
// This satisfies the BoolGateType interface.
func (g BoolGate) BoolValue() bool {
	return g.value
}
