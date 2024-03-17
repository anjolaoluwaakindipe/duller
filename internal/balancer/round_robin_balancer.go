package balancer

import (
	"sync"

	"github.com/anjolaoluwaakindipe/duller/internal/registry"
	"github.com/anjolaoluwaakindipe/duller/internal/service"
)

type RoundRobin struct {
	reg     registry.Registry
	pathMap map[string]int
	mutex   sync.Mutex
}

func (lb *RoundRobin) AddService(service *service.ServiceInfo) error {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	if err := lb.reg.RegisterService(service); err != nil {
		return err
	}
	return nil
}

func (lb *RoundRobin) RemoveService(service *service.ServiceInfo) error {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	if err := lb.reg.DeregisterService(service.Path, service.ServiceId); err != nil {
		return err
	}
	return nil
}

func (lb *RoundRobin) GetNextService(path string) (*service.ServiceInfo, error) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	services, err := lb.reg.GetServicesByPath(path)
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

func (lb *RoundRobin) GetRegitry() registry.Registry {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	return lb.reg
}

func NewRoundRobinLoadBalancer(reg registry.Registry) LoadBalancer {
	return &RoundRobin{reg: reg, pathMap: make(map[string]int)}
}
