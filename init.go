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

func exercise(logger hclog.Logger, servicers map[string]service.Servicer) error {
	for name, svcr := range servicers {
		u := svcr.Unit()
		logger.Debug(fmt.Sprintf("Exercising %s unit: %#v", name, u))

		if err := svcr.Start(); err != nil {
			return fmt.Errorf("%s failed to start", name)
		}

		logger.Debug(fmt.Sprintf("After Start, state of %s is: %v", name, svcr.Status()))

		if err := svcr.Stop(); err != nil {
			return fmt.Errorf("%s failed to stop", name)
		}

		logger.Debug(fmt.Sprintf("After Stop, state of %s is: %v", name, svcr.Status()))

		if err := svcr.Restart(); err != nil {
			return fmt.Errorf("%s failed to restart", name)
		}

		logger.Debug(fmt.Sprintf("After Restart, state of %s is: %v", name, svcr.Status()))

	}
	return nil
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

	var servicers = map[string]service.Servicer{}

	// Launch and connect to each service binary
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

		// Cast to a Servicer
		if svcr, ok := raw.(service.Servicer); ok {
			servicers[name] = svcr
		} else {
			logger.Error(fmt.Sprintf("%s failed to launch", name))
		}
	}

	err = exercise(logger, servicers)
	if err != nil {
		log.Fatal(err)
	}

}
