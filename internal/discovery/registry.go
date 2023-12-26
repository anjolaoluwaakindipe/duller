package discovery

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/anjolaoluwaakindipe/duller/internal/utils"
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
	for k := range r.Services {
		reg := regexp.QuoteMeta(k)
		paths = append(paths, reg)
	}
	r.servicPathRegex = "^(" + strings.Join(paths, "|") + ")"
}

// registring and updating services on the network
// if heartbeat is implemented the timeCreated field is always updated
func (r *InMemoryRegistry) RegisterService(msg RegisterServiceMessage) error {
	if len(msg.Path) == 0 {
		return fmt.Errorf("path field is empty")
	}
	utils.MakeUrlPathValid(&msg.Path)
	service, exist := r.Services[msg.Path]
	r.Mutex.Lock()
	if !exist {
		r.Services[msg.Path] = ServiceInfo{timeCreated: time.Now(), serverName: msg.ServiceName, path: msg.Path, address: msg.Address}
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
	r.Mutex.Lock()
	val, ok := r.Services[servicePath]
	r.Mutex.Unlock()
	if !ok {
		return nil, fmt.Errorf("path '%v' does not exist", path)
	}
	return &val, nil
}

func (r *InMemoryRegistry) RefreshRegistry(duration time.Duration) {
	r.Mutex.Lock()
	for path, service := range r.Services {
		fmt.Printf("%v", duration)
		if time.Since(service.timeCreated) > duration {
			delete(r.Services, path)
		}
	}
	r.Mutex.Unlock()
}

func removeDeadServices(r Registry, duration time.Duration, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(duration)
			r.RefreshRegistry(duration)
		}
	}
}

func InitInMemoryRegistry() Registry {
	return &InMemoryRegistry{Services: make(map[string]ServiceInfo)}
}
