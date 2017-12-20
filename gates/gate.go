package gates

import (
	"github.com/calavera/go-flipper/actor"
	"github.com/calavera/go-flipper/feature"
)

const (
	scalingFactor uint32 = 1000
	rangeFactor   uint32 = 100

	// BoolGateKey is the key for a BoolGate
	BoolGateKey GateKey = "boolean"
	// ActorGateKey is the key for an ActorGate
	ActorGateKey GateKey = "actors"
	// GroupGateKey is the key for a GroupGateKey
	GroupGateKey GateKey = "groups"
	// PercentageOfActorsGateKey is the key for a PercentageOfActorsGate
	PercentageOfActorsGateKey GateKey = "percentage_of_actors"
	// PercentageOfTimeGateKey is the key for a PercentageOfTimeGate
	PercentageOfTimeGateKey GateKey = "percentage_of_time"
)

// Set is a key set.
type Set map[string]string

// GateKey is the type for gate keys.
type GateKey string

// Gate represents a feature constraint.
type Gate interface {
	Key() GateKey
	IsOpen(f feature.Feature, a actor.Actor) bool
}

// BoolGateType represents a gate that uses boolean values.
type BoolGateType interface {
	BoolValue() bool
}

// IntGateType represents a gate that uses int values.
type IntGateType interface {
	IntValue() int
}

// SetGateType represents a gate that uses set values.
type SetGateType interface {
	SetValue() Set
}
