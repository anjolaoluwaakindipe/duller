package discovery_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/anjolaoluwaakindipe/duller/internal/discovery"
	"github.com/anjolaoluwaakindipe/duller/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Hub_Register(t *testing.T) {
	hub := discovery.NewInMemoryHub()
	ctx, cancel := context.WithCancel(context.Background())
	client := discovery.NewSocketClient(hub, nil)
	var wg sync.WaitGroup
	go func(ctx context.Context, wg *sync.WaitGroup) {
		wg.Add(1)
		hub.Run(ctx)
		wg.Done()
	}(ctx, &wg)

	hub.Register() <- &client
	cancel()
	wg.Wait()

	inMemoryHub, ok := hub.(*discovery.InMemoryHub)
	assert.True(t, ok)
	assert.Equal(t, 1, len(inMemoryHub.SocketClients))
}

func Test_Hub_UnRegister(t *testing.T) {
	hub := discovery.NewInMemoryHub()
	ctx1, cancel := context.WithCancel(context.Background())
	client1 := discovery.NewSocketClient(hub, nil)
	client2 := discovery.NewSocketClient(hub, nil)

	var wg sync.WaitGroup
	go func(ctx context.Context, wg *sync.WaitGroup) {
		wg.Add(1)
		hub.Run(ctx)
		wg.Done()
	}(ctx1, &wg)

	hub.Register() <- &client1
	hub.Register() <- &client2
	cancel()
	wg.Wait()

	inMemoryHub, ok := hub.(*discovery.InMemoryHub)
	assert.True(t, ok)
	assert.Equal(t, 2, len(inMemoryHub.SocketClients))

	ctx2, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context, wg *sync.WaitGroup) {
		wg.Add(1)
		hub.Run(ctx)
		wg.Done()
	}(ctx2, &wg)

	hub.Unregister() <- &client1
	hub.Unregister() <- &client2
	cancel()
	wg.Wait()

	assert.Equal(t, 0, len(inMemoryHub.SocketClients))
}

func Test_Hub_Broadcaster(t *testing.T) {
	t.Parallel()
	hub := discovery.NewInMemoryHub()
	ctx1, cancel1 := context.WithCancel(context.Background())
	client1 := discovery.NewSocketClient(hub, nil)

	var wg1 sync.WaitGroup
	var wg2 sync.WaitGroup

	go func(ctx context.Context, wg *sync.WaitGroup) {
		wg.Add(1)
		hub.Run(ctx)
		wg.Done()
	}(ctx1, &wg1)

	hub.Register() <- &client1
	cancel1()
	wg1.Wait()

	ctx1, cancel1 = context.WithCancel(context.Background())
	ctx2 := context.Background()
	wg1.Add(1)
	wg2.Add(1)

	go func(client *discovery.SocketClient, wg *sync.WaitGroup, ctx context.Context) {
		defer wg.Done()
		client.ReadPipe(ctx)
	}(&client1, &wg1, ctx2)

	go func(ctx context.Context, wg *sync.WaitGroup) {
		defer wg.Done()
		hub.Run(ctx)
	}(ctx1, &wg2)

	hub.Broadcaster() <- []byte("hello")
	wg1.Wait()
	time.Sleep(10 * time.Millisecond)
	cancel1()
	wg2.Wait()

	inMemoryHub, ok := hub.(*discovery.InMemoryHub)
	assert.True(t, ok)
	assert.Equal(t, 0, len(inMemoryHub.SocketClients))
}

func Test_Hub_Broadcaster_With_WebSocket_Connection(t *testing.T) {
	t.Parallel()
	hub := discovery.NewInMemoryHub()
	textMessage := 1

	message := []byte("message")
	mockWs := new(mocks.MockWebSocket)
	writeCloser := new(mocks.MockWriteCloser)
	writeCloser.On("Write", message).Return(len(message), nil)
	writeCloser.On("Close").Return(nil)
	mockWs.On("SetWriteDeadline", mock.Anything).Return(nil)
	mockWs.On("NextWriter", textMessage).Return(writeCloser, nil)
	mockWs.On("Close").Return(nil)
	socket := discovery.NewSocketClient(hub, mockWs)
	ctx1, cancel1 := context.WithCancel(context.Background())
	ctx2, cancel2 := context.WithCancel(context.Background())

	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func(ctx context.Context, hub discovery.Hub, wg *sync.WaitGroup) {
		defer wg.Done()
		hub.Run(ctx)
		fmt.Println("hub was cancelled")
	}(ctx2, hub, wg)

	go func(ctx context.Context, client *discovery.SocketClient, wg *sync.WaitGroup) {
		defer wg.Done()
		socket.ReadPipe(ctx)
		fmt.Println("socket client was cancelled")
	}(ctx1, &socket, wg)

	hub.Register() <- &socket
	time.Sleep(100 * time.Millisecond)
	hub.Broadcaster() <- message
	fmt.Println("hello1")

	time.Sleep(100 * time.Millisecond)
	cancel1()
	time.Sleep(100 * time.Millisecond)
	cancel2()

	wg.Wait()
	mockWs.AssertExpectations(t)
	mockWs.AssertExpectations(t)
	mockWs.AssertNumberOfCalls(t, "NextWriter", 1)
	writeCloser.AssertNumberOfCalls(t, "Write", 1)
}
