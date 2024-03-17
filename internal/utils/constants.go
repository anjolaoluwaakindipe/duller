package utils

import "time"

// defaults for the gateway & discovery cli flags
const (
	DISCOVERY_HOST           = "localhost"
	DISCOVERY_PORT           = "9876"
	DISCOVERY_TYPE           = "tcp"
	DISCOVERY_SERVICE_PATH   = "/get-service/"
	DISCOVERY_HEARTBEAT_PATH = "/sendheartbeat"
	GATEWAY_PORT             = "5923"
	GATEWAY_GRACEFULL_WAIT   = 15 * time.Second
	HEARTBEAT_INTERVAL       = 15 * time.Second
	DISCOVERY_KEY            = ""
)

// flag names for the gateway and cli commands
const (
	DISCOVERY_HOST_FLAG           = "dhost"
	DISCOVERY_PORT_FLAG           = "dport"
	DISCOVERY_TYPE_FLAG           = "dtype"
	DISCOVERY_SERVICE_PATH_FLAG   = "dservice_path"
	DISCOVERY_HEARTBEAT_PATH_FLAG = "dheartbeat_path"
	GATEWAY_PORT_FLAG             = "gport"
	GATEWAY_GRACEFULL_WAIT_FLAG   = "gwait"
	HEARTBEAT_INTERVAL_FLAG       = "dheartbeat"
	DISCOVERY_KEY_FLAG            = "dkey"
)
