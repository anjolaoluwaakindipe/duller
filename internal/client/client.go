package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/anjolaoluwaakindipe/duller/internal/discovery"
)


type DiscoveryClient struct {
	ServerName string
	Path       string
	Address    string
}

func (dc DiscoveryClient) SendHeartBeat(interval time.Duration, serverLocation string) {
	for {
		conn, err := net.Dial("tcp", serverLocation)
		if err != nil {
			log.Println("Error connecting to server:", err)
			return
		}
		defer conn.Close()
		message := discovery.Message{
			Type: "registerServiceMsg",
			Data: discovery.RegisterServiceMessage{
				ServerName: dc.ServerName,
				Path:       dc.Path,
				Address:    dc.Address,
			},
		}
		jsonMessage, _ := json.Marshal(message)
		_, err = conn.Write(jsonMessage)
		if err != nil {
			log.Println("Error sending JSON message to server:", err)
			return
		}

		decoder, _ := io.ReadAll(conn)

		fmt.Println("Sent JSON message to server:", string(decoder))
		time.Sleep(interval)
	}
}

type DiscoveryClientBuilder struct {
	discoveryClient DiscoveryClient
}

func (dcb *DiscoveryClientBuilder) SetServerName(serverName string) *DiscoveryClientBuilder {
	dcb.discoveryClient.ServerName = serverName
	return dcb
}

func (dcb *DiscoveryClientBuilder) SetPath(path string) *DiscoveryClientBuilder {
	dcb.discoveryClient.Path = path
	return dcb
}

func (dcb *DiscoveryClientBuilder) SetAddress(address string) *DiscoveryClientBuilder {
	dcb.discoveryClient.Address = address
	return dcb
}

func (dcb DiscoveryClientBuilder) Build() *DiscoveryClient {
	return &dcb.discoveryClient
}

func InitDiscoveryClient() *DiscoveryClientBuilder {
	client := DiscoveryClient{}
	return &DiscoveryClientBuilder{
		discoveryClient: client,
	}
}
