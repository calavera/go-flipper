package mongodb

import (
	"os"
	"testing"

	mgo "gopkg.in/mgo.v2"

	"github.com/calavera/go-flipper/actor/testhelpers"
	"github.com/calavera/go-flipper/feature"
	"github.com/calavera/go-flipper/gates"
	"github.com/stretchr/testify/require"
)

const testConnectionURL = "FLIPPER_MONGODB_URL"

func TestMongoDB(t *testing.T) {
	url := os.Getenv(testConnectionURL)
	if url == "" {
		t.SkipNow()
	}

	t.Run("configure", func(t *testing.T) {
		d := NewDriver()
		err := d.Configure(map[string]interface{}{
			"url": url,
		})
		require.NoError(t, err)
	})

	session, err := mgo.Dial(url)
	require.NoError(t, err)

	db := session.DB("")
	collection := db.C(defaultCollectionName)
	driver := NewDriverWithCollection(collection)

	t.Run("initialize with collection", func(t *testing.T) {
		err = driver.Configure(map[string]interface{}{
			"url": url,
		})
		require.NoError(t, err)
	})

	t.Run("enable for system", func(t *testing.T) {
		gate := gates.NewBoolGate(true)
		feat := feature.NewFeature("test")

		err := driver.Enable(feat, gate)
		require.NoError(t, err)

		g, err := driver.Get(feat, []gates.GateKey{gates.BoolGateKey})
		require.NoError(t, err)
		require.Len(t, g, 1)

		require.IsType(t, gates.BoolGate{}, g[0])
		b := g[0].(gates.BoolGate)
		require.Equal(t, true, b.BoolValue())
	})

	db.DropDatabase()

	t.Run("enable for actor", func(t *testing.T) {
		actor := testhelpers.Actor{"id"}

		gate := gates.NewActorGate(gates.NewSet(actor.ID))
		feat := feature.NewFeature("test")

		err := driver.Enable(feat, gate)
		require.NoError(t, err)

		g, err := driver.Get(feat, []gates.GateKey{gates.ActorGateKey})
		require.NoError(t, err)
		require.Len(t, g, 1)

		require.IsType(t, gates.ActorGate{}, g[0])
		b := g[0].(gates.ActorGate)

		set := b.SetValue()
		require.Len(t, set, 1)
		require.Contains(t, set, actor.FlipperID())
	})

	db.DropDatabase()

	t.Run("enable for groups", func(t *testing.T) {
		gate := gates.NewGroupGate(gates.NewSet("admins"))
		feat := feature.NewFeature("test")

		err := driver.Enable(feat, gate)
		require.NoError(t, err)

		g, err := driver.Get(feat, []gates.GateKey{gates.GroupGateKey})
		require.NoError(t, err)
		require.Len(t, g, 1)

		require.IsType(t, gates.GroupGate{}, g[0])
		b := g[0].(gates.GroupGate)

		set := b.SetValue()
		require.Len(t, set, 1)
		require.Contains(t, set, "admins")
	})

	db.DropDatabase()
	t.Run("enable for percentage of actors", func(t *testing.T) {
		feat := feature.NewFeature("test")
		gate := gates.NewPercentageOfActorsGate(30)

		err := driver.Enable(feat, gate)
		require.NoError(t, err)

		g, err := driver.Get(feat, []gates.GateKey{gates.PercentageOfActorsGateKey})
		require.NoError(t, err)
		require.Len(t, g, 1)

		require.IsType(t, gates.PercentageOfActorsGate{}, g[0])
		b := g[0].(gates.PercentageOfActorsGate)

		require.Equal(t, 30, b.IntValue())
	})

	db.DropDatabase()

	t.Run("enable for percentage of time", func(t *testing.T) {
		feat := feature.NewFeature("test")
		gate := gates.NewPercentageOfTimeGate(30)

		err := driver.Enable(feat, gate)
		require.NoError(t, err)

		g, err := driver.Get(feat, []gates.GateKey{gates.PercentageOfTimeGateKey})
		require.NoError(t, err)
		require.Len(t, g, 1)

		require.IsType(t, gates.PercentageOfTimeGate{}, g[0])
		b := g[0].(gates.PercentageOfTimeGate)

		require.Equal(t, 30, b.IntValue())
	})
}
