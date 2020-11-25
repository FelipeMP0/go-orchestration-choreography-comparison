package datastore

import (
	"github.com/google/wire"
)

// WireSet for datastores.
var WireSet = wire.NewSet(
	NewServiceStateDatastore,
)

// Datastore interface.
type Datastore interface {
	CollectionName() string
}
