package service

import (
	"context"
	"time"
)

type Registry interface {
	// Takes in a specific RegisterServiceMessage an stores that service. This will update a service
	// if the path and service name already exists.
	RegisterService(msg *ServiceInfo) error
	// GetServiceById takes in an ID  and returns its corresponding service.
	GetServiceById(serviceId string) (*ServiceInfo, error)
	// GetServiceByPath takes in a path returns all services associated with that given path.
	// Note: an error is return if the path does not exist. A path could exist and not have
	// any service tied to it thus the slice return can still have a length of zero even if there
	// is no error
	GetServicesByPath(path string) ([]*ServiceInfo, error)
	// Returns all available services in Registry
	GetServices() []*ServiceInfo
	// RefreshRegistry helps remove dead services. This is meant to be used in a goroutine
	RefreshRegistry(duration time.Duration, ctx context.Context)
	// DeregisterService a service from registry given a path and serviceId
	// returns an error if an invalid path or serviceId is given
	DeregisterService(path string, serviceId string) error
	// UpdateServiceCurrentUse is used to increment the CurrentUse of a service until it
	// reaches its weighted use value
	UpdateServiceCurrentUse(serviceId string)
	// IsServiceWeightFull checks if the CurrentUse of a service has reached its weighted value
	IsServiceWeightFull(serviceId string) (bool, error)
	// ResetCurrentUse resets the current use of the specified service to 0
	ResetCurrentUse(serviceId string)
}
