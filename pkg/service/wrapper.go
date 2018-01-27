package service

import (
	"net/rpc"

	plugin "github.com/hashicorp/go-plugin"
)

var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "SERVICER_PLUGIN",
	MagicCookieValue: "SERVICER",
}

type Wrapper struct {
	// Impl Injection
	Impl Servicer
}

func (w *Wrapper) Server(*plugin.MuxBroker) (interface{}, error) {
	return &ServicerRPCServer{Impl: w.Impl}, nil
}

func (Wrapper) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &ServicerRPC{Client: c}, nil
}
