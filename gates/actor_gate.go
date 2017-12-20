package gates

import (
	"github.com/calavera/go-flipper/actor"
	"github.com/calavera/go-flipper/feature"
)

// ActorGate is a gate to check if a feature is open for an actor.
type ActorGate struct {
	value Set
}

// NewActorGate initializes an ActorGate with a set of ids.
func NewActorGate(set Set) ActorGate {
	return ActorGate{set}
}

// Key returns the GateKey for an ActorGate gate.
func (ActorGate) Key() GateKey {
	return ActorGateKey
}

// IsOpen check if the gate is open for an feature and an actor.
// It uses the actor set to know whether the gate is open or not.
func (g ActorGate) IsOpen(f feature.Feature, a actor.Actor) bool {
	_, ok := g.value[a.FlipperID()]
	return ok
}

// SetValue returns the set of actors for which the gate is open.
// This satisfies the SetGateType interface.
func (g ActorGate) SetValue() Set {
	return g.value
}
