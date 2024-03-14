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
	discoveryKey               string
	discoveryPort              string
	discoveryServicePath       string
	discoveryHeartbeatPath     string
	discoveryHeartbeatInterval time.Duration
}

// Name returns the name of the command
func (dc *DiscCommand) Name() string {
	return dc.fs.Name()
}

// Init initializes DiscCommand fields with appropriate flags
func (dc *DiscCommand) Init(args []string) error {
	dc.fs.Usage = func() {
		fmt.Printf("disc usage: %s disc [OPTIONS] argument ...\n", os.Args[0])
		dc.fs.PrintDefaults()
		fmt.Printf("\n\n")
	}
	dc.fs.StringVar(&dc.discoveryKey, "dkey", utils.DISCOVERY_KEY, "Secret Key used for secure service registration. If empty or not passed no security will be applied to service registration.")
	dc.fs.StringVar(&dc.discoveryPort, "dport", utils.REGISTRY_PORT, "The PORT number the discovery should run on")
	dc.fs.StringVar(&dc.discoveryHost, "dhost", utils.REGISTRY_HOST, "The Host for the discovery server")
	dc.fs.StringVar(&dc.discoveryServicePath, "dservice_path", utils.DISCOVERY_SERVICE_PATH, "Path the discoveryService will use to proxy requests to corresponding services")
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

	router, err := NewMuxRouter(loadBalancer, serviceRegistry, WithSecretKey(dc.discoveryKey))
	if err != nil {
		return err
	}

	return InitRegistryServer(DiscoveryConfig{
		DISCOVERY_HOST:     dc.discoveryHost,
		DISCOVERY_PORT:     dc.discoveryPort,
		DISCOVERY_KEY:      dc.discoveryKey,
		HEARTBEAT_INTERVAL: dc.discoveryHeartbeatInterval,
	}, ctx, router)
}

func NewDiscCommand() *DiscCommand {
	return &DiscCommand{
		fs: flag.NewFlagSet("disc", flag.ContinueOnError),
	}
}
