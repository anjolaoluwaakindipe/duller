package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/anjolaoluwaakindipe/duller/internal/balancer"
	"github.com/anjolaoluwaakindipe/duller/internal/service"
	"github.com/mitchellh/mapstructure"
)

// RegistrySettings holds configuration for the Registry Server
type RegistrySettings struct {
	REGISTRY_HOST      string
	REGISTRY_PORT      string
	REGISTRY_TYPE      string
	HEARTBEAT_INTERVAL time.Duration
}

// InitRegistryServer initiates a TCP server and accepts connections for the registry
func InitRegistryServer(balancer balancer.LoadBalancer, rs RegistrySettings, registry service.Registry, ctx context.Context, routerCreator func(balancer.LoadBalancer, service.Registry) Router) error {
	router := routerCreator(balancer, registry)
	server := http.Server{
		Addr:    ":" + rs.REGISTRY_PORT,
		Handler: router.SetupRoutes(),
	}

	server.ListenAndServe()

	return nil
}
