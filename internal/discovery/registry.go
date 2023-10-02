package discovery

import (
	"fmt"
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


func (r *Registry) RegisterService(msg registerServiceMessage) error {

		_, exist := r.Services[msg.Path]
		if exist {
			return fmt.Errorf("service already exists")
		}
		mutex.Lock()
		r.Services[msg.Path] = ServiceInfo{timeCreated: time.Now(), serverName: msg.ServerName, path: msg.Path, address: msg.Address}
		mutex.Unlock()

		return nil
}
