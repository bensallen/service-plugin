package service

import "github.com/u-root/service-plugin/pkg/service/state"

type ServicerRPCServer struct {
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

func (s *ServicerRPCServer) Status(args interface{}, resp *state.Value) error {
	state := s.Impl.Status()
	*resp = state
	return nil
}

func (s *ServicerRPCServer) Unit(args interface{}, resp *Unit) error {
	u := s.Impl.Unit()
	*resp = u
	return nil
}
