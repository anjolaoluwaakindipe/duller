package discovery

import "time"


type ServiceInfo struct {
	timeCreated time.Time
	serverName  string
	address     string
	path        string
}