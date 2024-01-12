package discovery

import (
	"context"
	"log"
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
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	server.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
	return nil
}
