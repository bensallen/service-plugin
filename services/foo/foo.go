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
type foo service.Unit

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
		EnvironmentFile: "",
		PIDFile:         "",
		State:           state.Unknown,
	}
}

func (f *foo) Unit() service.Unit {
	return service.Unit(*f)
}

func (f *foo) Start() error {
	//f.Logger.Debug("Hello world")
	f.State = state.Starting

	return nil
}

func (f *foo) Stop() error {
	//f.Logger.Debug("Goodbye world")
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

func (f *foo) Status() error {
	//f.Logger.Debug(fmt.Sprintf("Foo is %v", f.State))
	return nil

}

func main() {

	var logger = hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

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
