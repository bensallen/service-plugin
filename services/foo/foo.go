package main

import (
	"os"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/u-root/service-plugin/pkg/service"
	"github.com/u-root/service-plugin/pkg/service/onfail"
	"github.com/u-root/service-plugin/pkg/service/state"
)

// foo is a bar
type foo service.Unit

var logger = hclog.New(&hclog.LoggerOptions{
	Level:      hclog.Trace,
	Output:     os.Stderr,
	JSONFormat: true,
})

//New returns a service.Unit of foo
func New() service.Servicer {

	return &foo{
		Name:            "foo",
		Description:     "foo does all of the bar",
		Type:            service.Simple,
		OnFail:          onfail.Restart,
		Before:          []string{},
		After:           []string{},
		Requires:        []string{},
		Enabled:         true,
		EnvironmentFile: "/etc/service.d/foo.conf",
		State:           state.Unknown,
	}
}

func (f *foo) Unit() service.Unit {
	return service.Unit(*f)
}

func (f *foo) Start() error {
	time.Sleep(5 * time.Second)
	logger.Debug("Hello world")
	f.State = state.Starting

	return nil
}

func (f *foo) Stop() error {
	logger.Debug("Goodbye world")
	f.State = state.Stopping

	return nil
}

func (f *foo) Restart() error {
	f.Stop()
	f.Start()

	return nil
}

func (f *foo) Reload() error {
	f.Restart()

	return nil
}

func (f *foo) Status() state.Value {
	return f.State

}

func main() {

	foo := New()

	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"foo": &service.Wrapper{Impl: foo},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: service.HandshakeConfig,
		Plugins:         pluginMap,
		Logger:          logger,
	})
}
