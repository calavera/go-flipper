package flipper

import (
	"github.com/calavera/go-flipper/driver"
	"github.com/calavera/go-flipper/client"
	"github.com/pkg/errors"
)

// NewClient initializes a Client with a store driver.
// The Driver must be registered before creating a new client.
// The configuration is mapped to the driver requirements before the client is initialized.
func NewClient(driverName string, config map[string]interface{}) (*client.Client, error) {
	a := driver.Get(driverName)
	if a == nil {
		return nil, errors.Errorf("Flipper driver not registered with name: %s", driverName)
	}

	if err := a.Configure(config); err != nil {
		return nil, errors.Wrapf(err, "Configuration error for Flipper driver %s", driverName)
	}

	return client.NewClient(a), nil
}
