package discovery

import (
	"fmt"
	"sync"
	"time"
)

type Registry interface {
	RegisterService(msg registerServiceMessage) error
	GetService(path string) (*ServiceInfo, error)
	RefreshRegistry(duration time.Duration)
}
type InMemoryRegistry struct {
	Mutex    sync.Mutex
	Services map[string]ServiceInfo
}

// helps edit the path to a more 
func (r *InMemoryRegistry) makePathValid(str *string) {
	val := *str
	if val[0] != '/' {
		*str = "/" + *str
	}
	if val[len(val)-1] == '/' {
		*str = (*str)[:len(*str)-1]
	}
}

// registring and updating services on the network
// if heartbeat is implemented the timeCreated field is always updated
func (r *InMemoryRegistry) RegisterService(msg registerServiceMessage) error {
	if len(msg.Path) == 0 {
		return fmt.Errorf("path field is empty")
	}
	r.makePathValid(&msg.Path)
	service, exist := r.Services[msg.Path]
	r.Mutex.Lock()
	if !exist {
		r.Services[msg.Path] = ServiceInfo{timeCreated: time.Now(), serverName: msg.ServerName, path: msg.Path, address: msg.Address}
		r.Mutex.Unlock()
		return nil
	}
	service.timeCreated = time.Now()
	r.Mutex.Unlock()
	return nil
}

// Get a specific service based on the path string as input
func (r *InMemoryRegistry) GetService(path string) (*ServiceInfo, error) {
	val, ok := r.Services[path]

	if !ok {
		return nil, fmt.Errorf("path '%v' does not exist", path)
	}
	return &val, nil
}

func (r *InMemoryRegistry) RefreshRegistry(duration time.Duration) {
	r.Mutex.Lock()
	for path, service := range r.Services {
		if time.Since(service.timeCreated) > duration+1 {
			delete(r.Services, path)
		}
	}
	r.Mutex.Unlock()
}

func removeDeadServices(r Registry, duration time.Duration) {
	for {
		time.Sleep(duration)
		r.RefreshRegistry(duration)
	}
}

func InitInMemoryRegistry() Registry {
	return &InMemoryRegistry{Services: make(map[string]ServiceInfo)}
}
