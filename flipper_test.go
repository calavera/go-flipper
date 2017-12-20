package flipper_test

import (
	"fmt"
	"strings"
	"testing"

	flipper "github.com/calavera/go-flipper"
	"github.com/calavera/go-flipper/actor"
	"github.com/calavera/go-flipper/actor/testhelpers"
	"github.com/calavera/go-flipper/gates"
	"github.com/stretchr/testify/require"

	_ "github.com/calavera/go-flipper/driver/memory"
	_ "github.com/calavera/go-flipper/driver/mongodb"
)

// ExampleNewClient_Global shows how to initialize
// the Flipper client with a store driver.
func ExampleNewClient_Global() {
	// Use an import with side-effects to load the store driver you want to use.
	// For example:
	// import _ "github.com/calavera/go-flipper/driver/memory"
	// This will load the driver to be initialized by the client.
	c, err := flipper.NewClient("memory", nil)
	if err != nil {
		panic("error initializing flipper client")
	}

	c.Enable("feature")
	enabled, err := c.IsEnabled("feature")
	if err != nil {
		panic("error loading feature")
	}

	if enabled {
		fmt.Println("feature enabled")
	}
}

// ExampleNewClient_MongoDB shows how to initialize
// the Flipper client with a store driver.
func ExampleNewClient_MongoDB() {
	// Use an import with side-effects to load the store driver you want to use.
	// For example:
	// import _ "github.com/calavera/go-flipper/driver/mongodb"
	// This will load the driver to be initialized by the client.
	config := map[string]interface{}{
		"url":        "127.0.0.1:27017",
		"database":   "testing",
		"collection": "flipper",
	}
	c, err := flipper.NewClient("mongodb", config)
	if err != nil {
		panic("error initializing flipper client")
	}

	c.Enable("feature")
	enabled, err := c.IsEnabled("feature")
	if err != nil {
		panic("error loading feature")
	}

	if enabled {
		fmt.Println("feature enabled")
	}
}

// ExampleNewClient_Actor shows how to enable a feature for a given actor.
func ExampleNewClient_Actor() {
	// Use an import with side-effects to load the store driver you want to use.
	// For example:
	// import _ "github.com/calavera/go-flipper/driver/memory"
	// This will load the driver to be initialized by the client.
	c, err := flipper.NewClient("memory", nil)
	if err != nil {
		panic("error initializing flipper client")
	}

	a := testhelpers.Actor{"id"}

	c.EnableForActors("feature", a)
	enabled, err := c.IsEnabled("feature", a)
	if err != nil {
		panic("error loading feature")
	}

	if enabled {
		fmt.Println("feature enabled")
	}
}

// ExampleNewClient_Group shows how to enable a feature for a group of actors.
func ExampleNewClient_Group() {
	// Use an import with side-effects to load the store driver you want to use.
	// For example:
	// import _ "github.com/calavera/go-flipper/driver/memory"
	// This will load the driver to be initialized by the client.
	c, err := flipper.NewClient("memory", nil)
	if err != nil {
		panic("error initializing flipper client")
	}

	a := testhelpers.Actor{"Admin;id"}

	gates.RegisterGroup("admins", func(a actor.Actor) bool {
		return strings.HasPrefix(a.FlipperID(), "Admin;")
	})

	c.EnableForGroups("feature", "admins")
	enabled, err := c.IsEnabled("feature", a)
	if err != nil {
		panic("error loading feature")
	}

	if enabled {
		fmt.Println("feature enabled")
	}
}

// ExampleNewClient_PercentageOfActors shows how to enable a feature for a percentage of actors.
func ExampleNewClient_PercentageOfActors() {
	// Use an import with side-effects to load the store driver you want to use.
	// For example:
	// import _ "github.com/calavera/go-flipper/driver/memory"
	// This will load the driver to be initialized by the client.
	c, err := flipper.NewClient("memory", nil)
	if err != nil {
		panic("error initializing flipper client")
	}

	a := testhelpers.Actor{"id"}

	c.EnableForPercentageOfActors("feature", 20)
	enabled, err := c.IsEnabled("feature", a)
	if err != nil {
		panic("error loading feature")
	}

	if enabled {
		fmt.Println("feature enabled")
	}
}

// ExampleNewClient_PercentageOfTime shows how to enable a feature for a percentage of checks
func ExampleNewClient_PercentageOfTime() {
	// Use an import with side-effects to load the store driver you want to use.
	// For example:
	// This will load the driver to be initialized by the client.
	c, err := flipper.NewClient("memory", nil)
	if err != nil {
		panic("error initializing flipper client")
	}

	a := testhelpers.Actor{"id"}

	c.EnableForPercentageOfTime("feature", 20)
	enabled, err := c.IsEnabled("feature", a)
	if err != nil {
		panic("error loading feature")
	}

	if enabled {
		fmt.Println("feature enabled")
	}
}

func TestNewClient(t *testing.T) {
	_, err := flipper.NewClient("memory", nil)
	require.NoError(t, err)
}
