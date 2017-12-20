package client

import (
	"strings"
	"testing"

	"github.com/calavera/go-flipper/actor"
	"github.com/calavera/go-flipper/actor/testhelpers"
	"github.com/calavera/go-flipper/driver/memory"
	"github.com/calavera/go-flipper/gates"
	"github.com/stretchr/testify/require"
)

func TestClient_IsOpen(t *testing.T) {
	client := NewClient(memory.NewDriver())
	a := testhelpers.Actor{"58474832756cfb0015870214"}

	t.Run("global flag", func(t *testing.T) {
		enabled, err := client.IsEnabled("test", a)
		require.NoError(t, err)
		require.False(t, enabled)

		err = client.Enable("test")
		require.NoError(t, err)

		enabled, err = client.IsEnabled("test", a)
		require.NoError(t, err)
		require.True(t, enabled)

		enabled, err = client.IsEnabled("test")
		require.NoError(t, err)
		require.True(t, enabled)

		err = client.Disable("test")
		require.NoError(t, err)

		enabled, err = client.IsEnabled("test", a)
		require.NoError(t, err)
		require.False(t, enabled)

		enabled, err = client.IsEnabled("test")
		require.NoError(t, err)
		require.False(t, enabled)
	})

	t.Run("actor flag", func(t *testing.T) {
		enabled, err := client.IsEnabled("test", a)
		require.NoError(t, err)
		require.False(t, enabled)

		err = client.EnableForActors("test", a)
		require.NoError(t, err)

		enabled, err = client.IsEnabled("test", a)
		require.NoError(t, err)
		require.True(t, enabled)

		err = client.DisableForActors("test", a)
		require.NoError(t, err)

		enabled, err = client.IsEnabled("test", a)
		require.NoError(t, err)
		require.False(t, enabled)
	})

	t.Run("percentage of actors flag", func(t *testing.T) {
		enabled, err := client.IsEnabled("test", a)
		require.NoError(t, err)
		require.False(t, enabled)

		err = client.EnableForPercentageOfActors("test", 30)
		require.NoError(t, err)

		enabled, err = client.IsEnabled("test", a)
		require.NoError(t, err)
		require.True(t, enabled)

		err = client.DisableForPercentageOfActors("test")
		require.NoError(t, err)

		enabled, err = client.IsEnabled("test", a)
		require.NoError(t, err)
		require.False(t, enabled)
	})

	t.Run("group of actors flag", func(t *testing.T) {
		gates.RegisterGroup("admins", func(a actor.Actor) bool {
			return strings.HasPrefix(a.FlipperID(), "Admin;")
		})

		enabled, err := client.IsEnabled("test", a)
		require.NoError(t, err)
		require.False(t, enabled)

		err = client.EnableForGroups("test", "admins")
		require.NoError(t, err)

		a = testhelpers.Actor{"Admin;" + a.FlipperID()}

		enabled, err = client.IsEnabled("test", a)
		require.NoError(t, err)
		require.True(t, enabled)

		err = client.DisableForGroups("test", "admins")
		require.NoError(t, err)

		enabled, err = client.IsEnabled("test", a)
		require.NoError(t, err)
		require.False(t, enabled)
	})
}
