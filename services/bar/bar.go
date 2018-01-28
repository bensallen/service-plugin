package main

import (
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/u-root/service-plugin/pkg/service"
	"github.com/u-root/service-plugin/pkg/service/onfail"
	"github.com/u-root/service-plugin/pkg/service/state"
)

// foo is a bar
type bar service.Unit

//New returns a service.Unit of foo
func New() service.Servicer {

	return &bar{
		Name:            "bar",
		Description:     "bar does all of the foo",
		Type:            service.Simple,
		OnFail:          onfail.Restart,
		Before:          []string{},
		After:           []string{},
		Requires:        []string{},
		Enabled:         true,
		EnvironmentFile: "/etc/service.d/bar.conf",
		State:           state.Unknown,
	}
}

func (b *bar) Unit() service.Unit {
	return service.Unit(*b)
}

func (b *bar) Start() error {
	//f.Logger.Debug("Hello world")
	b.State = state.Starting

	return nil
}

func (b *bar) Stop() error {
	//f.Logger.Debug("Goodbye world")
	b.State = state.Stopping

	return nil
}

func (b *bar) Restart() error {
	b.Stop()
	b.Start()

	return nil
}

func (b *bar) Reload() error {
	b.Restart()

	return nil
}

func (b *bar) Status() state.Value {
	return b.State
}

func main() {

	var logger = hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	bar := New()

	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"bar": &service.Wrapper{Impl: bar},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: service.HandshakeConfig,
		Plugins:         pluginMap,
		Logger:          logger,
	})
}
