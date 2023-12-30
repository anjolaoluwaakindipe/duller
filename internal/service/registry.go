package service

import (
	"time"
)

type Registry interface {
	// Takes in a specific RegisterServiceMessage an stores that service. This will update a service
	// if the path and service name already exists.
	RegisterService(msg ServiceInfo) error
	// GetServiceById takes in an ID  and returns its corresponding service.
	GetServiceById(serviceId string) (*ServiceInfo, error)
	// GetServiceByPath takes in a path returns all services associated with that given path.
	// Note: If an empty string is paseed as a serviceId then a service wi
	// that this could utilize load-balancing technique to distribute service access
	GetServicesByPath(path string) ([]*ServiceInfo, error)
	// Returns all available services in Registry
	GetServices() []*ServiceInfo
	// RefreshRegistry helps remove dead services. This is meant to be used in a goroutine
	RefreshRegistry(duration time.Duration)
	// DeregisterService a service from registry given a path and serviceId
	// returns an error if an invalid path or serviceId is given
	DeregisterService(path string, serviceId string) error
}
