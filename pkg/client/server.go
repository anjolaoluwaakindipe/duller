package duller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type ClientServerSettings struct {
	RegistryLocation  string
	DiscoveryIP       string
	DiscoverPort      string
	ClientPort        string
	HeartBeatInterval time.Duration
	Path              string
	ServerName        string
}

// Setups a test client and sends a heartbeat to service discovery
// server
func InitClientServer(settings ClientServerSettings) {
	serviceAddress := fmt.Sprintf("http://localhost:%v", settings.ClientPort)
	client, _ := NewDiscoveryClient(settings.ServerName, settings.Path, settings.DiscoveryIP, settings.DiscoverPort, WithHeartbeatInterval(settings.HeartBeatInterval))

	go client.SendHeartBeat(nil)

	err := http.ListenAndServe(fmt.Sprintf(":%v", settings.ClientPort), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		msg, _ := json.Marshal(map[string]interface{}{
			"message": fmt.Sprintf("Hello from server %v, with address %v, and you used path %v from the gateway to get to me", settings.ServerName, serviceAddress, settings.Path),
		})
		w.WriteHeader(http.StatusOK)
		w.Write(msg)
	}))
	if err != nil {
		log.Fatalf("Server error occured: %v", err)
	}
	fmt.Printf("starting server at port %v", settings.ClientPort)
}
