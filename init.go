package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/u-root/service-plugin/pkg/service"
)

func discover(path string, logger hclog.Logger) (map[string]*plugin.ClientConfig, error) {
	services, err := plugin.Discover("*", path)
	if err != nil {
		return nil, err
	}

	configs := map[string]*plugin.ClientConfig{}

	for _, bin := range services {

		name := filepath.Base(bin)
		pluginMap := map[string]plugin.Plugin{
			name: &service.Wrapper{},
		}

		configs[name] = &plugin.ClientConfig{
			HandshakeConfig: service.HandshakeConfig,
			Plugins:         pluginMap,
			Cmd:             exec.Command(bin),
			Logger:          logger,
		}
	}

	return configs, nil
}

func main() {
	// Create an hclog.Logger
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "init",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})

	configs, err := discover("./services/bin", logger)

	if err != nil {
		fmt.Printf("Service binary discover failed: %v\n", err)
	}

	for name, config := range configs {
		client := plugin.NewClient(config)
		defer client.Kill()

		// Connect via RPC
		rpcClient, err := client.Client()
		if err != nil {
			log.Fatal(err)
		}

		// Request the plugin
		raw, err := rpcClient.Dispense(name)
		if err != nil {
			log.Fatal(err)
		}

		// Cast back to a Servicer
		svcr := raw.(service.Servicer)

		err = svcr.Start()
		if err != nil {
			logger.Error(fmt.Sprintf("%s failed to start", name))
		}

		err = svcr.Stop()
		if err != nil {
			logger.Error(fmt.Sprintf("%s failed to stop", name))
			logger.Error("Foo failed to stop")
		}

		err = svcr.Restart()
		if err != nil {
			logger.Error(fmt.Sprintf("%s failed to restart", name))
		}

		err = svcr.Status()
		if err != nil {
			logger.Error(fmt.Sprintf("%s failed to get status", name))
		}

		u := svcr.Unit()

		logger.Debug(fmt.Sprintf("%s unit: %#v", name, u))

	}
	// We're a host! Start by launching the plugin process.

}
