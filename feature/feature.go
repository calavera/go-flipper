package feature

import (
	"fmt"

	"github.com/calavera/go-flipper/actor"
)

// Feature represents a feature check.
type Feature struct {
	Name string
}

// FeaturedActor is an specialized actor
// that appends the feature name as a prefix to the actor's FlipperID.
type FeaturedActor struct {
	f Feature
	a actor.Actor
}

// FlipperID returns the compound id for a FeaturedActor.
func (f FeaturedActor) FlipperID() string {
	return fmt.Sprintf("%s%s", f.f.Name, f.a.FlipperID())
}

// NewFeature initializes a feature by its name.
func NewFeature(name string) Feature {
	return Feature{
		Name: name,
	}
}

// NewFeaturedActor initializes a FeaturedActor with a feature and an actor.
func NewFeaturedActor(f Feature, a actor.Actor) FeaturedActor {
	return FeaturedActor{
		f: f,
		a: a,
	}
}
