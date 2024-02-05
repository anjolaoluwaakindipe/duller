package balancer

import "github.com/anjolaoluwaakindipe/duller/internal/service"

type LoadBalancer interface {
	// GetNextService uses implemented load balancing algorithm
	// to return a serviceInfo given a speific path.
	GetNextService(path string) (*service.ServiceInfo, error)
}
