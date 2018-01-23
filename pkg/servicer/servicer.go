package servicer

import (
	"net/rpc"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/u-root/service-plugin/pkg/servicer/onfail"
	"github.com/u-root/service-plugin/pkg/servicer/state"
)

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

// Servicer is the interface that we're exposing as a plugin.
type Servicer interface {
	Start() error
	Stop() error
	Reload() error
	Restart() error
	Status() error
	Unit() Unit
}

// Here is an implementation that talks over RPC
type ServicerRPC struct{ client *rpc.Client }

func (g *ServicerRPC) Start() error {
	var resp string
	return g.client.Call("Plugin.Start", new(interface{}), &resp)
}

func (g *ServicerRPC) Stop() error {
	var resp string
	return g.client.Call("Plugin.Stop", new(interface{}), &resp)
}

func (g *ServicerRPC) Reload() error {
	var resp string
	return g.client.Call("Plugin.Reload", new(interface{}), &resp)
}

func (g *ServicerRPC) Restart() error {
	var resp string
	return g.client.Call("Plugin.Restart", new(interface{}), &resp)
}

func (g *ServicerRPC) Status() error {
	var resp string
	return g.client.Call("Plugin.Status", new(interface{}), &resp)
}

func (g *ServicerRPC) Unit() Unit {
	var u Unit
	err := g.client.Call("Plugin.Unit", new(interface{}), &u)
	if err != nil {
		panic(err)
	}
	return u
}

// Here is the RPC server that GreeterRPC talks to, conforming to
// the requirements of net/rpc
type ServicerRPCServer struct {
	// This is the real implementation
	Impl Servicer
}

func (s *ServicerRPCServer) Start(args interface{}, resp *string) error {
	return s.Impl.Start()
}

func (s *ServicerRPCServer) Stop(args interface{}, resp *string) error {
	return s.Impl.Stop()
}

func (s *ServicerRPCServer) Reload(args interface{}, resp *string) error {
	return s.Impl.Reload()
}

func (s *ServicerRPCServer) Restart(args interface{}, resp *string) error {
	return s.Impl.Restart()
}

func (s *ServicerRPCServer) Status(args interface{}, resp *string) error {
	return s.Impl.Status()
}

func (s *ServicerRPCServer) Unit(args interface{}, resp *Unit) error {
	u := s.Impl.Unit()
	resp = &u
	return nil
}

type ServicerPlugin struct {
	// Impl Injection
	Impl Servicer
}

func (p *ServicerPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &ServicerRPCServer{Impl: p.Impl}, nil
}

func (ServicerPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &ServicerRPC{client: c}, nil
}
