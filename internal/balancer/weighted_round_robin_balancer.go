package balancer

import (
	"sync"

	"github.com/anjolaoluwaakindipe/duller/internal/service"
)

type WeightedRoundRobin struct {
	registry service.Registry
	mutex    sync.Mutex
}

func (wrb *WeightedRoundRobin) AddService(service *service.ServiceInfo) error {
	wrb.mutex.Lock()
	defer wrb.mutex.Unlock()
	if err := wrb.registry.RegisterService(*service); err != nil {
		return err
	}
	return nil
}

func (wrb *WeightedRoundRobin) RemoveService(service *service.ServiceInfo) error {
	wrb.mutex.Lock()
	defer wrb.mutex.Unlock()
	if err := wrb.registry.DeregisterService(service.Path, service.ServiceId); err != nil {
		return err
	}
	return nil
}

func (wrb *WeightedRoundRobin) GetNextService(path string) (*service.ServiceInfo, error) {
	wrb.mutex.Lock()
	defer wrb.mutex.Unlock()
	services, err := wrb.registry.GetServicesByPath(path)
	if err != nil {
		return nil, err
	}

	for i := range services {
		service := services[i]
	}
}
