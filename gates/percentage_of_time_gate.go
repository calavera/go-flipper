package gates

import (
	"math/rand"
	"time"

	"github.com/calavera/go-flipper/actor"
	"github.com/calavera/go-flipper/feature"
)

var seed = rand.New(rand.NewSource(time.Now().UnixNano()))

// PercentageOfTimeGate is a gate that's open only a percentage of times.
// The percentage can be any value between 0 and 100.
type PercentageOfTimeGate struct {
	value uint32
}

// NewPercentageOfTimeGate initializes a PercentageOfTimeGate gate
// with a given percentage.
func NewPercentageOfTimeGate(percentage int) PercentageOfTimeGate {
	return PercentageOfTimeGate{uint32(percentage)}
}

// Key returns the GateKey for a PercentageOfTimeGate gate.
func (PercentageOfTimeGate) Key() GateKey {
	return PercentageOfTimeGateKey
}

// IsOpen check if the gate is open for an feature and an actor.
// It calculates the likeliness of being open by its percentage.
func (g PercentageOfTimeGate) IsOpen(f feature.Feature, a actor.Actor) bool {
	r := seed.Uint32()
	return r < (g.value / rangeFactor)
}

// IntValue returns the gate's percentage as an int.
// This satisfies the IntGateType interface.
func (g PercentageOfTimeGate) IntValue() int {
	return int(g.value)
}
