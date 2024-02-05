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

	if len(services) == 0 {
		return nil, nil
	}

	var selectedService *service.ServiceInfo = nil

	for i := range services {
		service := services[i]
		isFull, err := wrb.registry.IsServiceWeightFull(service.ServiceId)
		if err != nil {
			return nil, err
		}
		if !isFull {
			wrb.registry.UpdateServiceCurrentUse(service.ServiceId)
			selectedService = service
			break
		}
	}

	if selectedService == nil {
		for i := range services {
			service := services[i]
			wrb.registry.ResetCurrentUse(service.ServiceId)
		}
		firstService := services[0]
		selectedService = firstService
		wrb.registry.UpdateServiceCurrentUse(firstService.ServiceId)
	}

	return selectedService, nil
}

func NewWeightedRoundRobinLoadBalancer(registry service.Registry) LoadBalancer {
	return &WeightedRoundRobin{registry: registry}
}
