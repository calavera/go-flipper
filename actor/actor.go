package actor

// Actor represents an entity for which a feature can be enabled for.
type Actor interface {
	FlipperID() string
}
