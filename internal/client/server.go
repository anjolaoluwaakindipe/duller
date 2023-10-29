package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type ClientServerSettings struct {
	RegistryLocation  string
	ClientPort        string
	HeartBeatInterval time.Duration
	Path              string
	ServerName        string
}

// Setups a test client and sends a heartbeat to service discovery
// server
func InitServer(settings ClientServerSettings) {
	serviceAddress := fmt.Sprintf("http://localhost:%v", settings.ClientPort)
	client := InitDiscoveryClient().
		SetAddress(serviceAddress).
		SetPath(settings.Path).
		SetServerName(settings.ServerName).
		Build()

	go client.SendHeartBeat(settings.HeartBeatInterval, settings.RegistryLocation)

	err := http.ListenAndServe(fmt.Sprintf(":%v", settings.ClientPort), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.Marshal(map[string]interface{}{
			"message": fmt.Sprintf("Hello from server %v, with address %v, and you used path %v from the gateway to get to me", settings.ServerName, serviceAddress, settings.Path),
		})
	}))

	if err != nil {
		log.Fatalf("Server error occured: %v", err)
	}
	fmt.Printf("starting server at port %v", settings.ClientPort)
}
