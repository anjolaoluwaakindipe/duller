package discovery_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/anjolaoluwaakindipe/duller/internal/discovery"
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
	registry := discovery.InitInMemoryRegistry(utils.NewClock())
	t.Run("When ctx is canceled before function invocation SHOULD send false through status channel and close goroutine ", func(t *testing.T) {
		fmt.Println("Hello")
		ctx, cancel := context.WithCancel(context.Background())
		status := make(chan bool)

		cancel()
		go discovery.InitRegistryServer(registrySettings, registry, ctx, status)
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
		registry := discovery.InitInMemoryRegistry(utils.NewClock())
		testServer, testClient := net.Pipe()

		go func(con net.Conn) {
			discovery.HandleRequest(testServer, registry)
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
		registry := discovery.InitInMemoryRegistry(utils.NewClock())
		testServer, testClient := net.Pipe()

		go func(net.Conn) {
			discovery.HandleRequest(testServer, registry)
		}(testServer)
		// Send discv
		message := discovery.Message{Type: "registerServiceMsg", Data: discovery.RegisterServiceMessage{ServiceName: "testService", Path: "/test", Address: "http://localhost:2020"}}
		messageAsBytes, err := json.Marshal(message)
		assert.Nil(t, err)
		testClient.Write(messageAsBytes)
	})
}
