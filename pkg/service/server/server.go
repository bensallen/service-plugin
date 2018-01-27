package server

import (
	"fmt"

	"github.com/u-root/service-plugin/pkg/service"
)

type ServicerRPCServer struct {
	Impl service.Servicer
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

func (s *ServicerRPCServer) Unit(args interface{}, resp *service.Unit) error {
	u := s.Impl.Unit()
	fmt.Printf("%#v\n", u)
	resp = &u
	return nil
}
