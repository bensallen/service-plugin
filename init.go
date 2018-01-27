package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/u-root/service-plugin/pkg/service"
	"github.com/u-root/service-plugin/pkg/service/wrapper"
)

func main() {
	// Create an hclog.Logger
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "init",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})

	// We're a host! Start by launching the plugin process.
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: wrapper.HandshakeConfig,
		Plugins:         pluginMap,
		Cmd:             exec.Command("./plugin/foo/foo"),
		Logger:          logger,
	})
	defer client.Kill()

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		log.Fatal(err)
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("foo")
	if err != nil {
		log.Fatal(err)
	}

	// Cast back to a Servicer
	foo := raw.(service.Servicer)

	err = foo.Start()
	if err != nil {
		logger.Error("Foo failed to start")
	}

	err = foo.Stop()
	if err != nil {
		logger.Error("Foo failed to stop")
	}

	err = foo.Restart()
	if err != nil {
		logger.Error("Foo failed to restart")
	}

	err = foo.Status()
	if err != nil {
		logger.Error("Foo failed to get status")
	}

	u := foo.Unit()
	fmt.Printf("Foo Unit: %#v\n", u)

}

// pluginMap is the map of plugins we can dispense.
var pluginMap = map[string]plugin.Plugin{
	"foo": &wrapper.ServicerWrapper{},
}
