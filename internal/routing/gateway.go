package routing

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type GatewaySetting struct {
	GATEWAY_PORT           string
	GATEWAY_GRACEFULL_WAIT time.Duration
}

func InitGateway(router Router, settings GatewaySetting) {
	router.RegisterRoutes()
	server := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%v", settings.GATEWAY_PORT),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router.GetRouter(),
	}

	go func() {
		err := server.ListenAndServe()
		if err == nil {
			log.Printf("Gateway server starting on port %v \n", settings.GATEWAY_PORT)
		}
		if err != nil {
			log.Printf("Gateway Server could not be started: %v \n", err)
		}
	}()

	sig_chan := make(chan os.Signal, 1)

	signal.Notify(sig_chan, os.Interrupt)

	<-sig_chan

	ctx, cancel := context.WithTimeout(context.Background(), settings.GATEWAY_GRACEFULL_WAIT)
	defer cancel()

	server.Shutdown(ctx)

	log.Println("Shutting down gateway server")
	os.Exit(0)
}

