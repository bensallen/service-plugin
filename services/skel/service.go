package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"os"

	"text/template"
)

var skelTmpl = []byte(`
package main

import (
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/u-root/service-plugin/pkg/service"
	"github.com/u-root/service-plugin/pkg/service/onfail"
	"github.com/u-root/service-plugin/pkg/service/state"
)

// {{ .Name }} is a service Unit that implimiments service.Servicer
type {{ .Name }} service.Unit

var logger = hclog.New(&hclog.LoggerOptions{
	Level:      hclog.Trace,
	Output:     os.Stderr,
	JSONFormat: true,
})

//New returns a service.Unit of {{ .Name }}
func New() service.Servicer {

	return &{{ .Name }}{
		Name:            "{{ .Name }}",
		Description:     "{{ .Name }} is a service",
		Type:            service.Simple,
		OnFail:          onfail.Restart,
		Before:          []string{},
		After:           []string{},
		Requires:        []string{},
		Enabled:         true,
		EnvironmentFile: "/etc/service.d/{{ .Name }}.conf",
		State:           state.Unknown,
	}
}

func ({{ .Var }} *{{ .Name }}) Unit() service.Unit {
	return service.Unit(*{{ .Var }})
}

func ({{ .Var }} *{{ .Name }}) Start() error {
	logger.Debug("Hello world")
	{{ .Var }}.State = state.Starting
	// Do something
	{{ .Var }}.State = state.Active

	return nil
}

func ({{ .Var }} *{{ .Name }}) Stop() error {
	logger.Debug("Goodbye world")
	{{ .Var }}.State = state.Stopping
	// Do something
	{{ .Var }}.State = state.Stopped

	return nil
}

func ({{ .Var }} *{{ .Name }}) Restart() error {
	{{ .Var }}.Stop()
	{{ .Var }}.Start()

	return nil
}

func ({{ .Var }} *{{ .Name }}) Reload() error {
	{{ .Var }}.Restart()

	return nil
}

func ({{ .Var }} *{{ .Name }}) Status() state.Value {
	return {{ .Var }}.State

}

func main() {

	{{ .Name }} := New()

	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"{{ .Name }}": &service.Wrapper{Impl: {{ .Name }}},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: service.HandshakeConfig,
		Plugins:         pluginMap,
		Logger:          logger,
	})
}
`)

type Service struct {
	Name string
	Var  string
}

func main() {

	var name string

	flag.StringVar(&name, "name", "", "Name of service")
	flag.Parse()

	if name == "" {
		fmt.Print("name argument must be specified\n")
		flag.Usage()
		os.Exit(1)
	}

	service := Service{name, name[0:1]}

	tmpl, err := template.New("serviceSkel").Parse(string(skelTmpl))
	if err != nil {
		fmt.Printf("Template parse failed, %v", err)
		os.Exit(1)
	}

	var b bytes.Buffer

	if err := tmpl.Execute(&b, service); err != nil {
		fmt.Printf("Template failed to execute, %v", err)
		os.Exit(1)
	}
	if fmted, err := format.Source(b.Bytes()); err != nil {
		fmt.Printf("Source failed to format, %v", err)
		os.Exit(1)
	} else {
		fmt.Printf("%s", fmted)
	}

}
