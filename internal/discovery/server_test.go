package discovery_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/anjolaoluwaakindipe/duller/internal/discovery"
	"github.com/stretchr/testify/assert"
)

func TestInitRegistryServer(t *testing.T) {
	registrySettings := discovery.RegistrySettings{
		REGISTRY_HOST:      "localhost",
		REGISTRY_PORT:      "2000",
		HEARTBEAT_INTERVAL: 1 * time.Second,
		REGISTRY_TYPE:      "tcp",
	}
	registry := discovery.InitInMemoryRegistry()
	t.Run("When ctx is canceled before function invocation SHOULD send false through status channel and close goroutine ", func(t *testing.T) {
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

func TestHandleRequest(t *testing.T) {

	registry := discovery.InitInMemoryRegistry()
	t.Run("When a valid registerServiceMsg message is sent through a connection SHOULD send an okay message through the same connection", func(t *testing.T) {
		testServer, testClient := net.Pipe()

		go func(net.Conn) {
			discovery.HandleRequest(testServer, registry)
			testServer.Close()
		}(testServer)
		message := discovery.Message{Type: "registerServiceMsg", Data: discovery.RegisterServiceMessage{ServiceName: "testService", Path: "/test", Address: "http://localhost:2020"}}
		messageAsBytes, err := json.Marshal(message)
		assert.Nil(t, err)
		testClient.Write(messageAsBytes)
		var response []byte
		_, err = testClient.Read(response)
		assert.Nil(t, err)
		var result discovery.RegistryResponse
		json.Unmarshal(response, &result)
		assert.NotNil(t, result)
		fmt.Printf("%v", result)
	})

	t.Run("When a valid ", func(t *testing.T) {
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
