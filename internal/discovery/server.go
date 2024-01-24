package discovery

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/anjolaoluwaakindipe/duller/internal/balancer"
	"github.com/anjolaoluwaakindipe/duller/internal/service"
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
		Addr:         ":" + rs.REGISTRY_PORT,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router.SetupRoutes(),
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		slog.Info(fmt.Sprintf("Starting Service Discovery Server on port %v... \n", rs.REGISTRY_PORT))
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			slog.Error("Error starting Service Discovery Server: %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	server.Shutdown(ctx)
	slog.Info("Shutting down Service Discovery Server...")
	return nil
}
