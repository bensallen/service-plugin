package service

import (
	hclog "github.com/hashicorp/go-hclog"
	"github.com/u-root/service-plugin/pkg/service/onfail"
	"github.com/u-root/service-plugin/pkg/service/state"
)

// Servicer is the interface that we're exposing as a plugin.
type Servicer interface {
	Start() error
	Stop() error
	Reload() error
	Restart() error
	Status() error
	Unit() Unit
}

type Unit struct {
	Name            string
	Description     string
	EnvironmentFile string
	PIDFile         string
	Type            Type
	OnFail          onfail.Action
	State           state.Value
	Before          []string
	After           []string
	Requires        []string
	Logger          hclog.Logger
}

type Type int

const (
	Simple  Type = 0
	Forking Type = 1
)
