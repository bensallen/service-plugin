package service

import (
	"github.com/u-root/service-plugin/pkg/service/onfail"
	"github.com/u-root/service-plugin/pkg/service/state"
)

// Servicer is the interface that we're exposing as a plugin.
type Servicer interface {
	Start() error
	Stop() error
	Reload() error
	Restart() error
	Status() state.Value
	Unit() Unit
}

type Unit struct {
	Name            string        // Name of service
	Description     string        // Description of service
	EnvironmentFile string        // File to source for overriding settings
	Type            Type          // Type of service
	OnFail          onfail.Action // Specified action if the service exists unexpectedly
	Enabled         bool          // Service enabled
	State           state.Value   // Current state of service
	Before          []string      // Start this service before specified services
	After           []string      // Start this service after specified services
	Requires        []string      // Start this service if specified services have successfully started
}

type Type int

const (
	Simple  Type = 0
	Forking Type = 1
)
