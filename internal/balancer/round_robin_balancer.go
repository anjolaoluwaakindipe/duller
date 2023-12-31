package balancer

import (
	"sync"

	"github.com/anjolaoluwaakindipe/duller/internal/service"
)

type RoundRobinLoadBalancer struct {
	registry service.Registry
	pathMap  map[string]int
	mutex    sync.Mutex
}

func (lb *RoundRobinLoadBalancer) AddService(service *service.ServiceInfo) error {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	if err := lb.registry.RegisterService(*service); err != nil {
		return err
	}
	return nil
}

func (lb *RoundRobinLoadBalancer) RemoveService(service *service.ServiceInfo) error {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	if err := lb.registry.DeregisterService(service.Path, service.ServiceId); err != nil {
		return err
	}
	return nil
}

func (lb *RoundRobinLoadBalancer) GetNextService(path string) (*service.ServiceInfo, error) {
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

func (lb *RoundRobinLoadBalancer) GetRegitry() service.Registry {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	return lb.registry
}

func NewRoundRobinLoadBalancer(registry service.Registry) LoadBalancer {
	return &RoundRobinLoadBalancer{registry: registry, pathMap: make(map[string]int)}
}
