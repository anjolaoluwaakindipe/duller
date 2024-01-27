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
		fmt.Printf("disc usage: %s disc [OPTIONS] argument ...\n", os.Args[0])
		gc.fs.PrintDefaults()
		fmt.Printf("\n\n")
	}
	flag.DurationVar(&gc.gatewayHearbeatInterval, "rheartbeat", utils.HEARTBEAT_INTERVAL, fmt.Sprintf("The interval of heartbeats expected. The default interval is %s", utils.HEARTBEAT_INTERVAL))
	gc.discoveryServicePath = *flag.String("dservice_path", utils.DISCOVERY_SERVICE_PATH, fmt.Sprintf("Path for proxying users. Default path is %v .", utils.DISCOVERY_SERVICE_PATH))
	flag.DurationVar(&gc.gatewayGracefullWait, "gwait", utils.GATEWAY_GRACEFULL_WAIT, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	gc.gatewayPort = *flag.String("gport", utils.GATEWAY_PORT, fmt.Sprintf("The PORT number the gateway should run on. Default value is %v", utils.GATEWAY_PORT))
	gc.discoveryPort = *flag.String("dport", utils.REGISTRY_PORT, fmt.Sprintf("The PORT number the discovery server is running on. Default value is %v", utils.REGISTRY_PORT))
	gc.discoveryHost = *flag.String("dhost", utils.REGISTRY_HOST, fmt.Sprintf("The IP Address/Host of the discovery server. Default value is %v", utils.REGISTRY_HOST))
	return gc.fs.Parse(args)
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
