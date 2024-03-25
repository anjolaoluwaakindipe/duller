package discovery

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/anjolaoluwaakindipe/duller/internal/balancer"
	"github.com/anjolaoluwaakindipe/duller/internal/registry"
	"github.com/anjolaoluwaakindipe/duller/internal/utils"
)

// DiscCommand is the command subset for a Discovery Server for the.
// It implement the Runner interface
type DiscCommand struct {
	fs                         *flag.FlagSet
	DiscoveryKey               string
	DiscoveryPort              string
	DiscoveryServicePath       string
	DiscoveryHeartbeatPath     string
	DiscoveryHeartbeatInterval time.Duration
}

// Name returns the name of the command
func (dc *DiscCommand) Name() string {
	return dc.fs.Name()
}

// Init initializes DiscCommand fields with appropriate flags
func (dc *DiscCommand) Init(args ...string) error {
	dc.fs.Usage = func() {
		fmt.Printf("disc usage: %s disc [OPTIONS] argument ...\n", os.Args[0])
		dc.fs.PrintDefaults()
		fmt.Printf("\n\n")
	}
	dc.fs.StringVar(&dc.DiscoveryKey, utils.DISCOVERY_KEY_FLAG, utils.DISCOVERY_KEY, "Secret Key used for secure service registration. If empty or not passed no security will be applied to service registration.")
	dc.fs.StringVar(&dc.DiscoveryPort, utils.DISCOVERY_PORT_FLAG, utils.DISCOVERY_PORT, "The PORT number the discovery should run on")
	dc.fs.StringVar(&dc.DiscoveryServicePath, utils.DISCOVERY_SERVICE_PATH_FLAG, utils.DISCOVERY_SERVICE_PATH, "Path the discoveryService will use to proxy requests to corresponding services")
	dc.fs.DurationVar(&dc.DiscoveryHeartbeatInterval, utils.HEARTBEAT_INTERVAL_FLAG, utils.HEARTBEAT_INTERVAL, "The interval of heartbeats expected")
	return dc.fs.Parse(args)
}

func (dc *DiscCommand) UsageInfo() {
	dc.Init()
	dc.fs.Usage()
}

func (dc *DiscCommand) Run() error {
	serviceRegistry := registry.InitInMemoryRegistry(utils.NewClock())
	loadBalancer := balancer.NewRoundRobinLoadBalancer(serviceRegistry)

	ctx := context.TODO()

	router, err := NewMuxRouter(loadBalancer, serviceRegistry, ctx, WithSecretKey(dc.DiscoveryKey))
	if err != nil {
		return err
	}

	return InitRegistryServer(DiscoveryConfig{
		DISCOVERY_PORT:     dc.DiscoveryPort,
		DISCOVERY_KEY:      dc.DiscoveryKey,
		HEARTBEAT_INTERVAL: dc.DiscoveryHeartbeatInterval,
	}, ctx, router)
}

func NewDiscCommand() *DiscCommand {
	return &DiscCommand{
		fs: flag.NewFlagSet("disc", flag.ContinueOnError),
	}
}
