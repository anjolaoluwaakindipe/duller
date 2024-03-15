package discovery

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// DiscoveryConfig holds configuration for the Registry Server
type DiscoveryConfig struct {
	DISCOVERY_HOST     string
	DISCOVERY_PORT     string
	DISCOVERY_KEY      string
	HEARTBEAT_INTERVAL time.Duration
}

// InitRegistryServer initiates a TCP server and accepts connections for the registry
func InitRegistryServer(dc DiscoveryConfig, ctx context.Context, router Router) error {
	server := http.Server{
		Addr:         ":" + dc.DISCOVERY_PORT,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router.SetupRoutes(),
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		slog.Info(fmt.Sprintf("Starting Service Discovery Server on port %v... \n", dc.DISCOVERY_PORT))
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
