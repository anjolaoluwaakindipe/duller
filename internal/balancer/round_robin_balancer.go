package balancer

import (
	"sync"

	"github.com/anjolaoluwaakindipe/duller/internal/service"
)

type RoundRobin struct {
	registry service.Registry
	pathMap  map[string]int
	mutex    sync.Mutex
}

func (lb *RoundRobin) AddService(service *service.ServiceInfo) error {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	if err := lb.registry.RegisterService(*service); err != nil {
		return err
	}
	return nil
}

func (lb *RoundRobin) RemoveService(service *service.ServiceInfo) error {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	if err := lb.registry.DeregisterService(service.Path, service.ServiceId); err != nil {
		return err
	}
	return nil
}

func (lb *RoundRobin) GetNextService(path string) (*service.ServiceInfo, error) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	services, err := lb.registry.GetServicesByPath(path)
	if err != nil {
		return nil, err
	}

	if len(services) == 0 {
		return nil, nil
	}

	serviceLength := len(services)

	index := lb.pathMap[path] % serviceLength

	service := services[index]

	lb.pathMap[path] = (index + 1) % serviceLength

	return service, nil
}

func (lb *RoundRobin) GetRegitry() service.Registry {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	return lb.registry
}

func NewRoundRobinLoadBalancer(registry service.Registry) LoadBalancer {
	return &RoundRobin{registry: registry, pathMap: make(map[string]int)}
}
