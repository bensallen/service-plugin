package wrapper

import (
	"net/rpc"

	plugin "github.com/hashicorp/go-plugin"
	"github.com/u-root/service-plugin/pkg/service"
	"github.com/u-root/service-plugin/pkg/service/client"
	"github.com/u-root/service-plugin/pkg/service/server"
)

var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "SERVICER_PLUGIN",
	MagicCookieValue: "SERVICER",
}

type ServicerWrapper struct {
	// Impl Injection
	Impl service.Servicer
}

func (p *ServicerWrapper) Server(*plugin.MuxBroker) (interface{}, error) {
	return &server.ServicerRPCServer{Impl: p.Impl}, nil
}

func (ServicerWrapper) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &client.ServicerRPC{Client: c}, nil
}
