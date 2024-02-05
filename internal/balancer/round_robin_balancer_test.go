package balancer_test

import (
	"testing"

	"github.com/anjolaoluwaakindipe/duller/internal/balancer"
	"github.com/anjolaoluwaakindipe/duller/internal/service"
	"github.com/anjolaoluwaakindipe/duller/internal/utils"
	"github.com/stretchr/testify/assert"
)

func stubFactory() (registry service.Registry, services []*service.ServiceInfo) {
	registry = service.InitInMemoryRegistry(utils.NewClock())

	services = []*service.ServiceInfo{
		{
			Port:        "4000",
			Path:        "/path1",
			ServiceId:   "server1",
			IP:          "localhost",
			WeightedUse: 5,
		},
		{
			Port:        "5000",
			Path:        "/path1",
			ServiceId:   "server2",
			IP:          "localhost",
			WeightedUse: 3,
		},
		{
			Port:        "6000",
			Path:        "/path1",
			ServiceId:   "server3",
			IP:          "localhost",
			WeightedUse: 2,
		},
	}

	for i := 0; i < len(services); i++ {
		registry.RegisterService(*services[i])
	}

	return
}

func Test_RoundRobinLoadBalancer_GetNextService(t *testing.T) {
	t.Run("SHOULD return a service WHEN a valid path is given", func(t *testing.T) {
		registry, _ := stubFactory()
		loadBalancer := balancer.NewRoundRobinLoadBalancer(registry)
		service, err := loadBalancer.GetNextService("/path1")

		assert.Nil(t, err)
		assert.NotNil(t, service)
	})

	t.Run("SHOULD return nil and an error WHEN an invalid path is given", func(t *testing.T) {
		registry, _ := stubFactory()
		loadBalancer := balancer.NewRoundRobinLoadBalancer(registry)
		service, err := loadBalancer.GetNextService("/invalid")

		assert.NotNil(t, err)
		assert.Nil(t, service)
	})

	t.Run("SHOULD cycle through all services WHEN given a valid path multiple times", func(t *testing.T) {
		registry, services := stubFactory()
		loadBalancer := balancer.NewRoundRobinLoadBalancer(registry)
		for i := 0; i < len(services); i++ {

			service, err := loadBalancer.GetNextService("/path1")

			assert.Nil(t, err)
			assert.NotNil(t, service)

			assert.Equal(t, services[i], service)
		}

		firstService, err := loadBalancer.GetNextService("/path1")
		assert.Nil(t, err)
		assert.Equal(t, services[0], firstService)
	})

	t.Run("SHOULD prevent use of dead services WHEN given a valid path", func(t *testing.T) {
		registry, services := stubFactory()
		loadBalancer := balancer.NewRoundRobinLoadBalancer(registry)

		firstService, err := loadBalancer.GetNextService("/path1")
		assert.Nil(t, err)
		assert.NotNil(t, firstService)
		assert.Equal(t, firstService, services[0])

		registry.DeregisterService("/path1", services[1].ServiceId)

		thirdService, err := loadBalancer.GetNextService("/path1")

		assert.Nil(t, err)
		assert.NotNil(t, thirdService)
		assert.Equal(t, thirdService, services[2])

		firstService, err = loadBalancer.GetNextService("/path1")
		assert.Nil(t, err)
		assert.NotNil(t, firstService)
		assert.Equal(t, firstService, services[0])
	})
}
