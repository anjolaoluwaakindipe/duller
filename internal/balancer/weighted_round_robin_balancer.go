package balancer

import (
	"sync"

	"github.com/anjolaoluwaakindipe/duller/internal/registry"
	"github.com/anjolaoluwaakindipe/duller/internal/service"
	"github.com/invopop/validation"
)

type WeightedRoundRobin struct {
	reg   registry.Registry
	mutex sync.Mutex
}

func (wrb *WeightedRoundRobin) validateService(service *service.ServiceInfo) error {
	return validation.ValidateStruct(service, validation.Field(service.WeightedUse, validation.Required, validation.Min(1)))
}

func (wrb *WeightedRoundRobin) AddService(service *service.ServiceInfo) error {
	wrb.mutex.Lock()
	defer wrb.mutex.Unlock()

	valErr := wrb.validateService(service)

	if valErr != nil {
		return valErr
	}

	if err := wrb.reg.RegisterService(service); err != nil {
		return err
	}
	return nil
}

func (wrb *WeightedRoundRobin) RemoveService(service *service.ServiceInfo) error {
	wrb.mutex.Lock()
	defer wrb.mutex.Unlock()
	if err := wrb.reg.DeregisterService(service.Path, service.ServiceId); err != nil {
		return err
	}
	return nil
}

func (wrb *WeightedRoundRobin) GetNextService(path string) (*service.ServiceInfo, error) {
	wrb.mutex.Lock()
	defer wrb.mutex.Unlock()
	services, err := wrb.reg.GetServicesByPath(path)
	if err != nil {
		return nil, err
	}

	if len(services) == 0 {
		return nil, nil
	}

	var selectedService *service.ServiceInfo = nil

	for i := range services {
		service := services[i]
		isFull, err := wrb.reg.IsServiceWeightFull(service.ServiceId)
		if err != nil {
			return nil, err
		}
		if !isFull {
			wrb.reg.UpdateServiceCurrentUse(service.ServiceId)
			selectedService = service
			break
		}
	}

	if selectedService == nil {
		for i := range services {
			service := services[i]
			wrb.reg.ResetCurrentUse(service.ServiceId)
		}
		firstService := services[0]
		selectedService = firstService
		wrb.reg.UpdateServiceCurrentUse(firstService.ServiceId)
	}

	return selectedService, nil
}

func NewWeightedRoundRobinLoadBalancer(reg registry.Registry) LoadBalancer {
	return &WeightedRoundRobin{reg: reg}
}
