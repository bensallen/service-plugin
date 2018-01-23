package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/u-root/service-plugin/pkg/service"
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
		HandshakeConfig: handshakeConfig,
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

	// We should have a Greeter now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
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

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "SERVICER_PLUGIN",
	MagicCookieValue: "SERVICER",
}

// pluginMap is the map of plugins we can dispense.
var pluginMap = map[string]plugin.Plugin{
	"foo": &service.ServicerPlugin{},
}
