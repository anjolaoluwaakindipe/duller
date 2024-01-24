package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/anjolaoluwaakindipe/duller/internal/balancer"
	"github.com/anjolaoluwaakindipe/duller/internal/discovery"
	routing "github.com/anjolaoluwaakindipe/duller/internal/gateway"
	"github.com/anjolaoluwaakindipe/duller/internal/service"
	"github.com/anjolaoluwaakindipe/duller/internal/utils"
)

// defaults for the gateway & discovery cli
const (
	REGISTRY_HOST          = "localhost"
	REGISTRY_PORT          = "9876"
	REGISTRY_TYPE          = "tcp"
	GATEWAY_PORT           = "5923"
	GATEWAY_GRACEFULL_WAIT = 15 * time.Second
	HEARTBEAT_INTERVAL     = 15 * time.Second
)

func main() {
	var gtimeout time.Duration
	var rheartbeat time.Duration
	rport := flag.String("rport", REGISTRY_PORT, "The PORT number the registry should run on")
	rhost := flag.String("rhost", REGISTRY_HOST, "Registry Host")
	gport := flag.String("gport", GATEWAY_PORT, "The PORT number the gateway should run on")
	flag.DurationVar(&gtimeout, "gtimeout", GATEWAY_GRACEFULL_WAIT, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.DurationVar(&rheartbeat, "rheartbeat", HEARTBEAT_INTERVAL, "The interval of heartbeats expected")
	flag.Parse()

	serviceRegistry := service.InitInMemoryRegistry(utils.NewClock())
	loadBalancer := balancer.NewRoundRobinLoadBalancer(serviceRegistry)

	ctx := context.Background()

	go discovery.InitRegistryServer(loadBalancer, discovery.RegistrySettings{
		REGISTRY_HOST:      *rhost,
		REGISTRY_PORT:      *rport,
		REGISTRY_TYPE:      REGISTRY_TYPE,
		HEARTBEAT_INTERVAL: rheartbeat,
	}, serviceRegistry, ctx, discovery.NewMuxRouter)

	gatewayRouter := routing.InitMuxRouter(fmt.Sprintf("%v:%v", *rhost, *rport))

	routing.InitGateway(gatewayRouter, routing.GatewaySetting{
		GATEWAY_PORT:           *gport,
		GATEWAY_GRACEFULL_WAIT: gtimeout,
	})
}
