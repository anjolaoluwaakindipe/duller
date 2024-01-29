package gateway

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/anjolaoluwaakindipe/duller/internal/utils"
)

type GateCommand struct {
	fs                      *flag.FlagSet
	gatewayHearbeatInterval time.Duration
	gatewayPort             string
	gatewayGracefullWait    time.Duration
	discoveryServicePath    string
	discoveryHost           string
	discoveryPort           string
}

func (gc *GateCommand) Name() string {
	return gc.fs.Name()
}

func (gc *GateCommand) Init(args []string) error {
	gc.fs.Usage = func() {
		fmt.Printf("gate usage: %s gate [OPTIONS] argument ...\n", os.Args[0])
		gc.fs.PrintDefaults()
		fmt.Printf("\n\n")
	}
	gc.fs.DurationVar(&gc.gatewayHearbeatInterval, "rheartbeat", utils.HEARTBEAT_INTERVAL, "The interval of heartbeats expected.")
	gc.discoveryServicePath = *gc.fs.String("dservice_path", utils.DISCOVERY_SERVICE_PATH, "Path for proxying users.")
	gc.fs.DurationVar(&gc.gatewayGracefullWait, "gwait", utils.GATEWAY_GRACEFULL_WAIT, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	gc.gatewayPort = *gc.fs.String("gport", utils.GATEWAY_PORT, "The PORT number the gateway should run on.")
	gc.discoveryPort = *gc.fs.String("dport", utils.REGISTRY_PORT, "The PORT number the discovery server is running on.")
	gc.discoveryHost = *gc.fs.String("dhost", utils.REGISTRY_HOST, "The IP Address/Host of the discovery server.")
	return gc.fs.Parse(args)
}

func (gc *GateCommand) UsageInfo() {
	gc.Init([]string{})
	gc.fs.Usage()
}

func (gc *GateCommand) Run() error {
	gatewayRouter := InitMuxRouter(
		WithDiscoveryHost(gc.discoveryHost),
		WithDiscoveryPort(gc.discoveryPort),
		WithDiscoveryPath(gc.discoveryServicePath),
	)

	InitGateway(gatewayRouter, GatewaySetting{
		GATEWAY_PORT:           gc.gatewayPort,
		GATEWAY_GRACEFULL_WAIT: gc.gatewayGracefullWait,
	})

	return nil
}

func NewGateCommand() *GateCommand {
	return &GateCommand{
		fs: flag.NewFlagSet("gate", flag.ContinueOnError),
	}
}
