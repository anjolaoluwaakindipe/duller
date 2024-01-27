package utils

import "time"

// defaults for the gateway & discovery cli
const (
	REGISTRY_HOST            = "localhost"
	REGISTRY_PORT            = "9876"
	REGISTRY_TYPE            = "tcp"
	DISCOVERY_SERVICE_PATH   = "/get-service/"
	DISCOVERY_HEARTBEAT_PATH = "/sendheartbeat"
	GATEWAY_PORT             = "5923"
	GATEWAY_GRACEFULL_WAIT   = 15 * time.Second
	HEARTBEAT_INTERVAL       = 15 * time.Second
)
