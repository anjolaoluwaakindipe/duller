package discovery

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"
)

type Registry interface {
	RegisterService(msg RegisterServiceMessage) error
	GetService(path string) (*ServiceInfo, error)
	RefreshRegistry(duration time.Duration)
}
type InMemoryRegistry struct {
	Mutex           sync.Mutex
	Services        map[string]ServiceInfo
	servicPathRegex string
}

// helps edit the path to a more
func (r *InMemoryRegistry) MakePathValid(str *string) {
	val := *str
	if val[0] != '/' {
		*str = "/" + *str
	}
	if val[len(val)-1] == '/' {
		*str = (*str)[:len(*str)-1]
	}
}

func (r *InMemoryRegistry) GetServicePathRegex() string {
	return r.servicPathRegex
}

func (r *InMemoryRegistry) GetPathFromRequest(str string) (string, error) {
	reg, err := regexp.Compile(r.servicPathRegex)

	if err != nil {
		return "", err
	}

	serviceNameIndexes := reg.FindAllStringIndex(str, -1)

	if serviceNameIndexes == nil {
		return "", fmt.Errorf("no services available")
	}

	serviceNameIndex := serviceNameIndexes[len(serviceNameIndexes)-1]
	return str[serviceNameIndex[0]:serviceNameIndex[1]], nil

}

func (r *InMemoryRegistry) SetServicePathRegex() {
	if len(r.Services) == 0 {
		return
	}
	paths := make([]string, 0)
	r.Mutex.Lock()
	for k := range r.Services {
		reg := regexp.QuoteMeta(k)
		paths = append(paths, reg)
	}
	r.servicPathRegex = "^(" + strings.Join(paths, "|") + ")"
	r.Mutex.Unlock()
}

// registring and updating services on the network
// if heartbeat is implemented the timeCreated field is always updated
func (r *InMemoryRegistry) RegisterService(msg RegisterServiceMessage) error {
	if len(msg.Path) == 0 {
		return fmt.Errorf("path field is empty")
	}
	r.MakePathValid(&msg.Path)
	service, exist := r.Services[msg.Path]
	r.Mutex.Lock()
	if !exist {
		r.Services[msg.Path] = ServiceInfo{timeCreated: time.Now(), serverName: msg.ServerName, path: msg.Path, address: msg.Address}
	} else {
		service.timeCreated = time.Now()
	}
	r.SetServicePathRegex()
	r.Mutex.Unlock()
	return nil
}

// Get a specific service based on the path string as input
func (r *InMemoryRegistry) GetService(path string) (*ServiceInfo, error) {
	servicePath, err := r.GetPathFromRequest(path)

	if err != nil {
		return nil, fmt.Errorf("invalid path")
	}

	val, ok := r.Services[servicePath]

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
