package gates

import (
	"hash/crc32"

	"github.com/calavera/go-flipper/actor"
	"github.com/calavera/go-flipper/feature"
)

// PercentageOfActorsGate is a gate that's open only for a percentage of the actors checked.
// The percentage can be any value between 0 and 100.
type PercentageOfActorsGate struct {
	value uint32
}

// NewPercentageOfActorsGate initializes a PercentageOfActorsGate gate
// with a given percentage.
func NewPercentageOfActorsGate(percentage int) PercentageOfActorsGate {
	return PercentageOfActorsGate{uint32(percentage)}
}

// Key returns the GateKey for a PercentageOfActorsGate gate.
func (PercentageOfActorsGate) Key() GateKey {
	return PercentageOfActorsGateKey
}

// IsOpen check if the gate is open for an feature and an actor.
// It calculates the likeliness of being open by its percentage.
func (g PercentageOfActorsGate) IsOpen(f feature.Feature, a actor.Actor) bool {
	return g.checksum(feature.NewFeaturedActor(f, a)) < (g.value * scalingFactor)
}

// IntValue returns the gate's percentage as an int.
// This satisfies the IntGateType interface.
func (g PercentageOfActorsGate) IntValue() int {
	return int(g.value)
}

func (g PercentageOfActorsGate) checksum(a actor.Actor) uint32 {
	p := crc32.ChecksumIEEE([]byte(a.FlipperID()))
	return p % (scalingFactor * rangeFactor)
}
