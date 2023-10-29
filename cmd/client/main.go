package main

import (
	"flag"
	"time"

	"github.com/anjolaoluwaakindipe/duller/internal/client"
)

// defaults for the test client cli
const (
	DISCOVERY_LOCATION = "localhost:9876"
	CLIENT_PORT        = "3000"
	GATEWAY_PATH       = "/test"
	SERVER_NAME        = "server1"
	HEARTBEAT_INTERVAL = 14 * time.Second
)

func main() {
	var heartBeatInterval time.Duration
	cport := flag.String("cport", CLIENT_PORT, "The PORT number the client server should run on")
	gpath := flag.String("gpath", GATEWAY_PATH, "The path the client server will use in the discovery server. This is used by the gateway when proxying")
	cname := flag.String("cname", SERVER_NAME, "Name of the client server")
	dlocation := flag.String("dlocation", DISCOVERY_LOCATION, "Address of the discovery server")
	flag.DurationVar(&heartBeatInterval, "cinterval", HEARTBEAT_INTERVAL, "Interval at which the client server will be sending out heartbeats")

	flag.Parse()

	clientServerSettings := client.ClientServerSettings{
		ClientPort:        *cport,
		Path:              *gpath,
		ServerName:        *cname,
		HeartBeatInterval: heartBeatInterval,
		RegistryLocation:  *dlocation,
	}
	client.InitServer(clientServerSettings)
}
