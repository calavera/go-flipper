package client

import (
	"errors"

	"github.com/calavera/go-flipper/actor"
	"github.com/calavera/go-flipper/driver"
	"github.com/calavera/go-flipper/feature"
	"github.com/calavera/go-flipper/gates"
)

var (
	globalChecks = []gates.GateKey{
		gates.BoolGateKey,
		gates.PercentageOfTimeGateKey,
	}

	actorChecks = []gates.GateKey{
		gates.BoolGateKey,
		gates.ActorGateKey,
		gates.GroupGateKey,
		gates.PercentageOfActorsGateKey,
		gates.PercentageOfTimeGateKey,
	}
)

// Client is used to access
// the feature flags.
type Client struct {
	driver driver.Driver
}

// NewClient initializes a client with a store driver.
// It assumes that the driver is properly configured.
// See flipper.NewClient as a shortcut to initialize
// a client.
func NewClient(a driver.Driver) *Client {
	return &Client{a}
}

// IsEnabled checks if a feature is enabled.
// It uses only global checks when there are not actors in the list.
// This check is accumulative, it only returns true if the feature is enabled
// for every actor. It returns false if the feature is disabled for any of the actors.
func (c *Client) IsEnabled(featureName string, actors ...actor.Actor) (bool, error) {
	if len(actors) > 0 {
		return c.isEnabledForActors(featureName, actors...)
	}

	return c.isEnabledGlobally(featureName)
}

// Enable enables a feature globally, for every actor.
func (c *Client) Enable(featureName string) error {
	gate := gates.NewBoolGate(true)
	return c.driver.Enable(feature.NewFeature(featureName), gate)
}

// Disable disables a feature globally.
// Actors might still have the feature enabled if other gates
// are open.
func (c *Client) Disable(featureName string) error {
	gate := gates.NewBoolGate(false)
	return c.driver.Disable(feature.NewFeature(featureName), gate)
}

// EnableForActors enables a featue for a list of actors.
func (c *Client) EnableForActors(featureName string, actors ...actor.Actor) error {
	if len(actors) == 0 {
		return errors.New("there are no actors to enable the feature for")
	}
	set := gates.Set{}
	for _, a := range actors {
		set[a.FlipperID()] = a.FlipperID()
	}
	gate := gates.NewActorGate(set)
	return c.driver.Enable(feature.NewFeature(featureName), gate)
}

// DisableForActors disables a featue for a list of actors.
func (c *Client) DisableForActors(featureName string, actors ...actor.Actor) error {
	if len(actors) == 0 {
		return errors.New("there are no actors to disable the feature for")
	}
	set := gates.Set{}
	for _, a := range actors {
		set[a.FlipperID()] = a.FlipperID()
	}
	gate := gates.NewActorGate(set)
	return c.driver.Disable(feature.NewFeature(featureName), gate)
}

// EnableForGroups enables a featue for a list of groups.
func (c *Client) EnableForGroups(featureName string, groups ...string) error {
	if len(groups) == 0 {
		return errors.New("there are no groups to enable the feature for")
	}
	set := gates.Set{}
	for _, n := range groups {
		set[n] = n
	}
	gate := gates.NewGroupGate(set)
	return c.driver.Enable(feature.NewFeature(featureName), gate)
}

// DisableForGroups disables a feature for a list of groups.
func (c *Client) DisableForGroups(featureName string, groups ...string) error {
	if len(groups) == 0 {
		return errors.New("there are no groups to disable the feature for")
	}
	set := gates.Set{}
	for _, n := range groups {
		set[n] = n
	}
	gate := gates.NewGroupGate(set)
	return c.driver.Disable(feature.NewFeature(featureName), gate)
}

// EnableForPercentageOfActors enables a feature for a percentage of the actors checked.
func (c *Client) EnableForPercentageOfActors(featureName string, percentage int) error {
	gate := gates.NewPercentageOfActorsGate(percentage)
	return c.driver.Enable(feature.NewFeature(featureName), gate)
}

// DisableForPercentageOfActors disables a feature for a percentage of the actors checked.
func (c *Client) DisableForPercentageOfActors(featureName string) error {
	gate := gates.NewPercentageOfActorsGate(0)
	return c.driver.Disable(feature.NewFeature(featureName), gate)
}

// EnableForPercentageOfTime enables a feature for a percentage of the checks.
func (c *Client) EnableForPercentageOfTime(featureName string, percentage int) error {
	gate := gates.NewPercentageOfTimeGate(percentage)
	return c.driver.Enable(feature.NewFeature(featureName), gate)
}

// DisableForPercentageOfTime disables a feature for a percentage of the checks.
func (c *Client) DisableForPercentageOfTime(featureName string) error {
	gate := gates.NewPercentageOfTimeGate(0)
	return c.driver.Disable(feature.NewFeature(featureName), gate)
}

func (c *Client) isEnabledGlobally(featureName string) (bool, error) {
	feat := feature.NewFeature(featureName)
	checks, err := c.driver.Get(feat, globalChecks)
	if err != nil {
		return false, err
	}

	if checks == nil || len(checks) == 0 {
		return false, nil
	}

	var open bool
	for _, g := range checks {
		if g.IsOpen(feat, nil) {
			open = true
			break
		}
	}

	return open, nil
}

func (c *Client) isEnabledForActors(featureName string, actors ...actor.Actor) (bool, error) {
	feat := feature.NewFeature(featureName)
	checks, err := c.driver.Get(feat, actorChecks)
	if err != nil {
		return false, err
	}

	if checks == nil || len(checks) == 0 {
		return false, nil
	}

	for _, a := range actors {
		open := false

		for _, g := range checks {
			if g.IsOpen(feat, a) {
				open = true
				break
			}
		}

		if !open {
			return false, nil
		}
	}

	return true, nil
}
