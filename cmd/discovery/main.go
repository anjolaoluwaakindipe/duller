package main

import (
	"flag"
	"time"

	"github.com/anjolaoluwaakindipe/duller/internal/discovery"
	"github.com/anjolaoluwaakindipe/duller/internal/routing"
)

// cli constants
const (
	REGISTRY_HOST          = "localhost"
	REGISTRY_PORT          = "9876"
	REGISTRY_TYPE          = "tcp"
	GATEWAY_PORT           = "5923"
	GATEWAY_GRACEFULL_WAIT = 15 * time.Second
)

func main() {
	var gtimeout time.Duration

	rport := flag.String("rport", REGISTRY_PORT, "The PORT number the registry should run on")
	rhost := flag.String("rhost", REGISTRY_HOST, "Registry Host")
	gport := flag.String("gport", GATEWAY_PORT, "The PORT number the gateway should run on")
	flag.DurationVar(&gtimeout, "gtimeout", GATEWAY_GRACEFULL_WAIT, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	serviceRegistry := discovery.InitInMemoryRegistry()

	go discovery.InitRegistryServer(discovery.RegistrySettings{
		REGISTRY_HOST: *rhost,
		REGISTRY_PORT: *rport,
		REGISTRY_TYPE: REGISTRY_TYPE,
	}, serviceRegistry)

	gatewayRouter := routing.InitMuxRouter()

	routing.InitGateway(gatewayRouter, routing.GatewaySetting{
		GATEWAY_PORT:           *gport,
		GATEWAY_GRACEFULL_WAIT: gtimeout,
	})
}
