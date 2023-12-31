package balancer

import "github.com/anjolaoluwaakindipe/duller/internal/service"

type LoadBalancer interface {
	GetNextService(path string) (*service.ServiceInfo, error)
}
