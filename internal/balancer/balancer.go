package balancer

import "github.com/anjolaoluwaakindipe/duller/internal/service"

type LoadBalancer interface {
	AddService(service *service.ServiceInfo)
	RemoveService(service *service.ServiceInfo)
	GetNextService() *service.ServiceInfo
}
