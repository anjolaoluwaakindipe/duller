package duller

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/anjolaoluwaakindipe/duller/internal/discovery"
)

type DiscoveryClient struct {
	heartbeatInterval time.Duration
	discoveryPort     string
	discoveryIP       string
	serviceId         string
	path              string
	ip                string
	port              string
	heartbeatPath     string
}

// SendHearbeat sends a hearbeat message to a servcie discovery server
// using the the properties contained in the DiscoveryClient. Logs out any
// error encoutered while sending a heartbeat.
//
// Note: This function in intended to be used within a goroutine.
func (dc DiscoveryClient) SendHeartBeat(ctx context.Context) {
	ticker := time.NewTicker(dc.heartbeatInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			message := discovery.HeartBeatMessage{
				ServiceId: dc.serviceId,
				Path:      dc.path,
				IP:        dc.ip,
				Port:      dc.port,
			}
			jsonMessage, err := json.Marshal(message)
			if err != nil {
				log.Println("Error occured when parsing heartbeat to json: ", err)
			}

			response, err := http.Post("http://"+dc.discoveryIP+":"+dc.discoveryPort+"/"+dc.heartbeatPath, "application/json", bytes.NewBuffer(jsonMessage))
			if err != nil {
				log.Println("Error occured when sending heartbeat: ", err)
			}

			if (response.StatusCode != http.StatusOK) && (response.StatusCode != http.StatusCreated) {
				body, _ := io.ReadAll(response.Body)
				log.Printf("Error response from discovery server with status code %v: \n %v \n", response.StatusCode, string(body))
			}
		}
	}
}

// WithHeartbeatInterval sets the HeartbeatInterval for the DiscoveryClient.
func WithHeartbeatInterval(interval time.Duration) DiscoveryClientOptions {
	return func(dc *DiscoveryClient) {
		dc.heartbeatInterval = interval
	}
}

// WithHeartbeatPath sets the HeartbeatPath for the DiscoveryClient
func WithHeartbeatPath(path string) DiscoveryClientOptions {
	return func(dc *DiscoveryClient) {
		dc.heartbeatPath = path
	}
}

// WithDiscoveryPort sets the DiscoveryPort for the DiscoveryClient
func WithDiscoveryPort(port string) DiscoveryClientOptions {
	return func(dc *DiscoveryClient) {
		dc.discoveryPort = port
	}
}

// WithDiscoveryIP sets the DiscoveryIP for the DiscoveryClient
func WithDiscoveryIP(ip string) DiscoveryClientOptions {
	return func(dc *DiscoveryClient) {
		dc.discoveryIP = ip
	}
}

// DiscoveryClientOptions is an option fucntion type for any DiscoveryClient
type DiscoveryClientOptions = func(dc *DiscoveryClient)

// NewDiscoveryClient creates a new service discovery client instance for sending heartbeats to a service
// discovery server
//
// Note: Default values for the client include
//
//	discoveryIP: "localhost"
//	discoveryPort: "9876",
//	heartbeatPath: "/heartbeat"
//	heartbeatInterval: 15 * time.Second
func NewDiscoveryClient(serviceId string, path string, ip string, port string, opts ...DiscoveryClientOptions) (DiscoveryClient, error) {
	dc := DiscoveryClient{serviceId: serviceId, path: path, ip: ip, port: port, discoveryIP: "localhost", discoveryPort: "9876", heartbeatPath: "/heartbeat", heartbeatInterval: 15 * time.Second}
	for _, opt := range opts {
		opt(&dc)
	}

	return dc, nil
}
