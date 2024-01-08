package discovery_test

import (
	"context"
	"encoding/json"
	"net"
	"testing"
	"time"

	"github.com/anjolaoluwaakindipe/duller/internal/balancer"
	"github.com/anjolaoluwaakindipe/duller/internal/discovery"
	"github.com/anjolaoluwaakindipe/duller/internal/service"
	"github.com/anjolaoluwaakindipe/duller/internal/utils"
	"github.com/stretchr/testify/assert"
)

func Test_InitRegistryServer(t *testing.T) {
	registrySettings := discovery.RegistrySettings{
		REGISTRY_HOST:      "localhost",
		REGISTRY_PORT:      "2000",
		HEARTBEAT_INTERVAL: 1 * time.Second,
		REGISTRY_TYPE:      "tcp",
	}
	t.Run("When ctx is canceled before function invocation SHOULD send false through status channel and close goroutine ", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		status := make(chan bool)

		registry := service.InitInMemoryRegistry(utils.NewClock())
		loadBalancer := balancer.NewRoundRobinLoadBalancer(registry)
		cancel()
		go discovery.InitRegistryServer(loadBalancer, registrySettings, registry, ctx, status)
		var finalResult bool

		for c := range status {
			if !c {
				finalResult = c
			}
		}

		assert.Equal(t, false, finalResult)
	})
}

func Test_HandleRequest(t *testing.T) {
	t.Run("When a valid registerServiceMsg message is sent through a connection SHOULD send an okay message through the same connection", func(t *testing.T) {
		registry := service.InitInMemoryRegistry(utils.NewClock())
		testServer, testClient := net.Pipe()
		loadBalancer := balancer.NewRoundRobinLoadBalancer(registry)

		go func(con net.Conn) {
			discovery.HandleRequest(testServer, registry, loadBalancer)
			con.Close()
		}(testServer)

		message := discovery.Message{Type: "registerServiceMsg", Data: discovery.RegisterServiceMessage{ServiceName: "testService", Path: "/test", Address: "http://localhost:2020"}}
		messageAsBytes, err := json.Marshal(message)

		assert.Nil(t, err)

		testClient.Write(messageAsBytes)

		var response []byte
		_, err = testClient.Read(response)

		assert.Nil(t, err)

		result := &discovery.RegistryResponse{}
		json.Unmarshal(response, result)

		assert.NotNil(t, result)
		testClient.Close()
	})

	t.Run("When a valid ", func(t *testing.T) {
		registry := service.InitInMemoryRegistry(utils.NewClock())
		loadBalancer := balancer.NewRoundRobinLoadBalancer(registry)
		testServer, testClient := net.Pipe()

		go func(net.Conn) {
			discovery.HandleRequest(testServer, registry, loadBalancer)
		}(testServer)
		// Send discv
		message := discovery.Message{Type: "registerServiceMsg", Data: discovery.RegisterServiceMessage{ServiceName: "testService", Path: "/test", Address: "http://localhost:2020"}}
		messageAsBytes, err := json.Marshal(message)
		assert.Nil(t, err)
		testClient.Write(messageAsBytes)
	})
}
