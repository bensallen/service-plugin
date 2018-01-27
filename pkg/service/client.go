package service

import (
	"net/rpc"
)

// ServicerRPC is an implemtation of Servicer over RPC
type ServicerRPC struct{ Client *rpc.Client }

func (g *ServicerRPC) Start() error {
	var resp string
	return g.Client.Call("Plugin.Start", new(interface{}), &resp)
}

func (s *ServicerRPC) Stop() error {
	var resp string
	return s.Client.Call("Plugin.Stop", new(interface{}), &resp)
}

func (s *ServicerRPC) Reload() error {
	var resp string
	return s.Client.Call("Plugin.Reload", new(interface{}), &resp)
}

func (s *ServicerRPC) Restart() error {
	var resp string
	return s.Client.Call("Plugin.Restart", new(interface{}), &resp)
}

func (s *ServicerRPC) Status() error {
	var resp string
	return s.Client.Call("Plugin.Status", new(interface{}), &resp)
}

func (s *ServicerRPC) Unit() Unit {
	var u Unit
	err := s.Client.Call("Plugin.Unit", new(interface{}), &u)
	if err != nil {
		panic(err)
	}
	return u
}
