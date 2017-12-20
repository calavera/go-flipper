package gates

import (
	"testing"

	"github.com/calavera/go-flipper/actor/testhelpers"
	"github.com/calavera/go-flipper/feature"
	"github.com/stretchr/testify/assert"
)

func TestPercentageOfActorsGate(t *testing.T) {
	f := feature.NewFeature("test")
	a := testhelpers.Actor{"58474832756cfb0015870214"}
	p := NewPercentageOfActorsGate(30)

	assert.Equal(t, uint32(28792), p.checksum(a))
	assert.True(t, p.IsOpen(f, a))
}
