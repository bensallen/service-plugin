package main

import (
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/u-root/service-plugin/pkg/service"
	"github.com/u-root/service-plugin/pkg/service/onfail"
	"github.com/u-root/service-plugin/pkg/service/state"
)

// baz is a service Unit that implimiments service.Servicer
type baz service.Unit

var logger = hclog.New(&hclog.LoggerOptions{
	Level:      hclog.Trace,
	Output:     os.Stderr,
	JSONFormat: true,
})

//New returns a service.Unit of baz
func New() service.Servicer {

	return &baz{
		Name:            "baz",
		Description:     "baz is a service",
		Type:            service.Simple,
		OnFail:          onfail.Restart,
		Before:          []string{},
		After:           []string{},
		Requires:        []string{"foo"},
		Enabled:         true,
		EnvironmentFile: "/etc/service.d/baz.conf",
		State:           state.Unknown,
	}
}

func (b *baz) Unit() service.Unit {
	return service.Unit(*b)
}

func (b *baz) Start() error {
	logger.Debug("Hello world")
	b.State = state.Starting
	// Do something
	b.State = state.Active

	return nil
}

func (b *baz) Stop() error {
	logger.Debug("Goodbye world")
	b.State = state.Stopping
	// Do something
	b.State = state.Stopped

	return nil
}

func (b *baz) Restart() error {
	b.Stop()
	b.Start()

	return nil
}

func (b *baz) Reload() error {
	b.Restart()

	return nil
}

func (b *baz) Status() state.Value {
	return b.State

}

func main() {

	baz := New()

	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"baz": &service.Wrapper{Impl: baz},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: service.HandshakeConfig,
		Plugins:         pluginMap,
		Logger:          logger,
	})
}
