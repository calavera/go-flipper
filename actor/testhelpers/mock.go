package testhelpers

// Actor is an struct that uses a string as the ID.
// This ID is used a FlipperID.
type Actor struct {
	ID string
}

// FlipperID returns the actor ID.
// It satisfies the actor.Actor interface.
func (a Actor) FlipperID() string {
	return a.ID
}
