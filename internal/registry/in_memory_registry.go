package registry

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/anjolaoluwaakindipe/duller/internal/service"
	"github.com/anjolaoluwaakindipe/duller/internal/utils"
	"github.com/invopop/validation"
)

// InMemoryRegistry is an in memory implementation of the
// Registry interface.
type InMemoryRegistry struct {
	mutex sync.Mutex
	// PathTable is a store for all the services indexed by their path.
	PathTable map[string][]*service.ServiceInfo
	// servicePathRegex is used to quickly get the available service path that has been stored from
	// a given external path.
	servicPathRegex string
	// Clock is an app interface for time.
	Clock utils.Clock
	// ServiceIdTable is a store for all the services indexed by their id
	ServiceIdTable map[string]*service.ServiceInfo
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
	if len(r.PathTable) == 0 {
		return
	}
	paths := make([]string, 0)
	for k := range r.PathTable {
		reg := regexp.QuoteMeta(k)
		paths = append(paths, reg)
	}
	r.servicPathRegex = "^(" + strings.Join(paths, "|") + ")"
}

func (r *InMemoryRegistry) validateService(msg *service.ServiceInfo) error {
	return validation.ValidateStruct(msg,
		validation.Field(&msg.IP, validation.Required),
		validation.Field(&msg.Port, validation.Required),
		validation.Field(&msg.ServiceId, validation.Required),
		validation.Field(&msg.Path, validation.Required),
	)
}

func (r *InMemoryRegistry) RegisterService(msg *service.ServiceInfo) error {
	var valerr error = r.validateService(msg)

	if valerr != nil {
		return valerr
	}

	utils.MakeUrlPathValid(&msg.Path)
	_, pathExist := r.PathTable[msg.Path]

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if !pathExist {
		r.PathTable[msg.Path] = []*service.ServiceInfo{msg}
		r.ServiceIdTable[msg.ServiceId] = msg
		r.SetServicePathRegex()
		return nil
	}

	service, serviceIdExist := r.ServiceIdTable[msg.ServiceId]

	if !serviceIdExist {
		r.PathTable[msg.Path] = append(r.PathTable[msg.Path], msg)
		r.ServiceIdTable[msg.ServiceId] = msg
		return nil
	}

	service.LastHeartbeat = r.Clock.Now()
	service.IP = msg.IP
	service.Port = msg.Port

	return nil
}

func (r *InMemoryRegistry) UpdateServiceCurrentUse(serviceId string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	service, exists := r.ServiceIdTable[serviceId]

	if !exists {
		return
	}

	if service.CurrentUse < service.WeightedUse {
		service.CurrentUse += 1
	}
	return
}

func (r *InMemoryRegistry) ResetCurrentUse(serviceId string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	service, exists := r.ServiceIdTable[serviceId]

	if !exists {
		return
	}

	service.CurrentUse = 0
}

func (r *InMemoryRegistry) IsServiceWeightFull(serviceId string) (bool, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	service, exists := r.ServiceIdTable[serviceId]

	if !exists {
		return false, fmt.Errorf("service id does not exist in registry")
	}

	if service.CurrentUse == service.WeightedUse {
		return true, nil
	}

	return false, nil
}

func (r *InMemoryRegistry) GetServicesByPath(path string) ([]*service.ServiceInfo, error) {
	servicePath, err := r.GetPathFromRequest(path)
	if err != nil {
		return nil, fmt.Errorf("invalid path")
	}
	r.mutex.Lock()
	services, exist := r.PathTable[servicePath]
	r.mutex.Unlock()

	if !exist {
		return nil, fmt.Errorf("path '%v' does not exist", path)
	}
	return services, nil
}

func (r *InMemoryRegistry) GetServiceById(serviceId string) (*service.ServiceInfo, error) {
	service, exist := r.ServiceIdTable[serviceId]

	if !exist {
		return nil, fmt.Errorf("service with serviceId '%v' does not exist", serviceId)
	}

	return service, nil
}

func (r *InMemoryRegistry) GetServices() []*service.ServiceInfo {
	services := make([]*service.ServiceInfo, 0)

	for _, service := range r.ServiceIdTable {
		services = append(services, service)
	}

	return services
}

func (r *InMemoryRegistry) DeregisterService(path string, serviceId string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	services, pathExist := r.PathTable[path]

	if !pathExist {
		return fmt.Errorf("path '%v' does not exist inside registry", path)
	}

	_, serviceExist := r.ServiceIdTable[serviceId]

	if !serviceExist {
		return fmt.Errorf("service with id '%v' does not exist", serviceId)
	}

	delete(r.ServiceIdTable, serviceId)

	rmInd := -1

	for ind, service := range services {
		if service.ServiceId == serviceId {
			rmInd = ind
			break
		}
	}

	if rmInd != -1 {
		r.PathTable[path][rmInd] = nil
		r.PathTable[path] = append(r.PathTable[path][:rmInd], r.PathTable[path][rmInd+1:]...)
	}

	return nil
}

func (r *InMemoryRegistry) RefreshRegistry(duration time.Duration, ctx context.Context) {
	ticker := time.NewTicker(duration)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.mutex.Lock()
			defer r.mutex.Unlock()
			deadServices := make([]struct {
				Path string
				Id   string
				Ind  int
			}, 0)

			for path, services := range r.PathTable {
				for ind, service := range services {
					if time.Now().After(service.LastHeartbeat.Add(duration).Add(1 * time.Second)) {
						deadServices = append(deadServices, struct {
							Path string
							Id   string
							Ind  int
						}{
							Path: path,
							Id:   service.ServiceId,
							Ind:  ind,
						})
					}
				}
			}

			for _, val := range deadServices {
				delete(r.ServiceIdTable, val.Id)
				r.PathTable[val.Path][val.Ind] = nil
				r.PathTable[val.Path] = append(r.PathTable[val.Path][:val.Ind], r.PathTable[val.Path][val.Ind+1:]...)
			}
		}
	}
}

func InitInMemoryRegistry(clock utils.Clock) Registry {
	return &InMemoryRegistry{PathTable: make(map[string][]*service.ServiceInfo), Clock: clock, ServiceIdTable: make(map[string]*service.ServiceInfo)}
}
