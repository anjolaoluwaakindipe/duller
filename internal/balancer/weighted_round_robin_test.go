package balancer_test

import (
	"testing"

	"github.com/anjolaoluwaakindipe/duller/internal/balancer"
	"github.com/anjolaoluwaakindipe/duller/internal/service"
	"github.com/stretchr/testify/assert"
)

func Test_WeightedRoundRobin_AddService(t *testing.T) {
	t.Run("SHOULD return an error WHEN service added has weightedUse less than 1", func(t *testing.T) {
		registry, _ := stubFactory()
		loadBalancer := balancer.NewWeightedRoundRobinLoadBalancer(registry)

		newService1 := service.ServiceInfo{
			Path:        "/path2",
			WeightedUse: -1,
			ServiceId:   "service2",
			Port:        "9999",
			IP:          "000.00.0",
		}

		newService2 := service.ServiceInfo{
			Path:      "/path2",
			ServiceId: "service2",
			Port:      "9999",
			IP:        "000.00.0",
		}

		err := loadBalancer.AddService(&newService1)
		assert.NotNil(t, err)

		err = loadBalancer.AddService(&newService2)
		assert.NotNil(t, err)
	})
}

func Test_WeightedRoundRobin_GetNextService(t *testing.T) {
	t.Run("SHOULD return a service WHEN given a valid path", func(t *testing.T) {
		registry, _ := stubFactory()
		loadBalancer := balancer.NewWeightedRoundRobinLoadBalancer(registry)

		service, err := loadBalancer.GetNextService("/path1")
		assert.Nil(t, err)
		assert.NotNil(t, service)
	})
	t.Run("SHOULD return nil and an error WHEN an invalid path is given", func(t *testing.T) {
		registry, _ := stubFactory()
		loadBalancer := balancer.NewWeightedRoundRobinLoadBalancer(registry)
		service, err := loadBalancer.GetNextService("/invalid")

		assert.NotNil(t, err)
		assert.Nil(t, service)
	})

	t.Run("SHOULD perform weighted cycle WHEN given a valid path multiple times", func(t *testing.T) {
		registry, stubServices := stubFactory()
		loadBalancer := balancer.NewWeightedRoundRobinLoadBalancer(registry)
		services := make([]*service.ServiceInfo, 0)

		for i := 0; i < 10; i++ {
			service, err := loadBalancer.GetNextService("/path1")
			assert.Nil(t, err)
			services = append(services, service)
		}
		service_count := make(map[string]int)
		for _, service := range services {
			_, exists := service_count[service.ServiceId]
			if !exists {
				service_count[service.ServiceId] = 0
			}
			service_count[service.ServiceId] += 1
		}

		for _, service := range stubServices {
			assert.Contains(t, service_count, service.ServiceId)
			assert.Equal(t, service.WeightedUse, service_count[service.ServiceId])
		}
	})
}
