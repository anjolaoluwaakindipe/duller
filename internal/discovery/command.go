package discovery

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/anjolaoluwaakindipe/duller/internal/balancer"
	"github.com/anjolaoluwaakindipe/duller/internal/service"
	"github.com/anjolaoluwaakindipe/duller/internal/utils"
)

// DiscCommand is the command subset for a Discovery Server for the.
// It implement the Runner interface
type DiscCommand struct {
	fs                         *flag.FlagSet
	discoveryHost              string
	discoveryPort              string
	discoveryType              string
	discoveryServicePath       string
	discoveryHeartbeatPath     string
	discoveryHeartbeatInterval time.Duration
}

// Name returns the name of the command
func (dc *DiscCommand) Name() string {
	return dc.fs.Name()
}

func (dc *DiscCommand) Init(args []string) error {
	dc.fs.Usage = func() {
		fmt.Printf("disc usage: %s disc [OPTIONS] argument ...\n", os.Args[0])
		dc.fs.PrintDefaults()
		fmt.Printf("\n\n")
	}
	dc.fs.StringVar(&dc.discoveryPort, "dport", utils.REGISTRY_PORT, "The PORT number the discovery should run on")
	dc.fs.StringVar(&dc.discoveryHost, "dhost", utils.REGISTRY_HOST, "The Host for the discovery server")
	dc.fs.StringVar(&dc.discoveryServicePath, "dservice_path", utils.DISCOVERY_SERVICE_PATH, "Path the discoveryService will use to proxy requests to corresponding services")
	dc.fs.DurationVar(&dc.discoveryHeartbeatInterval, "rheartbeat", utils.HEARTBEAT_INTERVAL, "The interval of heartbeats expected")
	return dc.fs.Parse(args)
}

func (dc *DiscCommand) UsageInfo() {
	dc.Init([]string{})
	dc.fs.Usage()
}

func (dc *DiscCommand) Run() error {
	serviceRegistry := service.InitInMemoryRegistry(utils.NewClock())
	loadBalancer := balancer.NewRoundRobinLoadBalancer(serviceRegistry)

	ctx := context.TODO()

	return InitRegistryServer(loadBalancer, RegistrySettings{
		REGISTRY_HOST:      dc.discoveryHost,
		REGISTRY_PORT:      dc.discoveryPort,
		REGISTRY_TYPE:      dc.discoveryType,
		HEARTBEAT_INTERVAL: dc.discoveryHeartbeatInterval,
	}, serviceRegistry, ctx, NewMuxRouter)
}

func NewDiscCommand() *DiscCommand {
	return &DiscCommand{
		fs: flag.NewFlagSet("disc", flag.ContinueOnError),
	}
}
