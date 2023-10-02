package discovery

import (
	"sync"
	"time"
)

var mutex sync.Mutex

type Registry struct {
	Mutex    sync.Mutex
	Services map[string]ServiceInfo
}

type ServiceInfo struct {
	timeCreated time.Time
	serverName  string
	address     string
	path        string
}



