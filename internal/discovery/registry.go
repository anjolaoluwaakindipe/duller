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
	// Takes in a specific RegisterServiceMessage an stores that service. This will update a service
	// if the path and service name already exists.
	RegisterService(msg RegisterServiceMessage) error
	// Takes in a path gives and returns a service.
	// Note: that this could utilize load-balancing technique to distribute service access
	GetService(path string) (*ServiceInfo, error)
	// Takes in a paht and returns a list of paths
	GetServices(path string) ([]*ServiceInfo, error)
	// RefreshRegistry helps remove dead services. This is meant to be used in a goroutine
	RefreshRegistry(duration time.Duration, ctx context.Context)

	// Removes service from registry
	DeregisterService(path string, serviceId string) error
}

// InMemoryRegistry is an in memory implementation of the
// Registry interface.
type InMemoryRegistry struct {
	Mutex sync.Mutex
	// Services is a store for all the services based indexed by their path.
	Services map[string]*ServiceInfo
	// servicePathRegex is used to quickly get the available service path that has been stored from
	// a given external path.
	servicPathRegex string
	// Clock is an app interface for time.
	Clock utils.Clock
}

func (r *InMemoryRegistry) GetServicePathRegex() string {
	return r.servicPathRegex
}

func (r *InMemoryRegistry) GetServices(path string) ([]*ServiceInfo, error) {
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
	return []*ServiceInfo{val}, nil
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
		r.Services[msg.Path] = &ServiceInfo{LastHeartbeat: r.Clock.Now(), ServiceId: msg.ServiceName, Path: msg.Path, Address: msg.Address}
	} else {
		service.LastHeartbeat = r.Clock.Now()
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
	return val, nil
}

func (r *InMemoryRegistry) RefreshRegistry(duration time.Duration, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(duration)
			r.Mutex.Lock()
			for path, service := range r.Services {
				fmt.Printf("%v", duration)
				if time.Since(service.LastHeartbeat) > duration {
					delete(r.Services, path)
				}
			}
			r.Mutex.Unlock()
		}
	}
}

func (r *InMemoryRegistry) DeregisterService(path string, id string) error {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	servicePath, err := r.GetPathFromRequest(path)
	if err != nil {
		return fmt.Errorf("invalid path")
	}
	r.Mutex.Lock()
	service, ok := r.Services[servicePath]
	r.Mutex.Unlock()
	if !ok {
		return fmt.Errorf("path '%v' does not exist", path)
	}

	if service.ServiceId != id {
		return fmt.Errorf("Server with id %v does not exist", id)
	}

	delete(r.Services, servicePath)
	return nil
}

func InitInMemoryRegistry(clock utils.Clock) Registry {
	return &InMemoryRegistry{Services: make(map[string]*ServiceInfo), Clock: clock}
}
