package discovery_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/anjolaoluwaakindipe/duller/internal/discovery"
	"github.com/stretchr/testify/assert"
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

	go func(client *discovery.SocketClient, wg *sync.WaitGroup, ctx context.Context) {
		wg.Add(1)
		client.ReadPipe(ctx)
		wg.Done()
	}(&client1, &wg1, ctx2)

	go func(ctx context.Context, wg *sync.WaitGroup) {
		wg.Add(1)
		hub.Run(ctx)
		wg.Done()
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

// func Test_Hub_Broadcaster_With_WebSocket_Connection(t *testing.T) {
// 	t.Parallel()
// 	hub := discovery.NewInMemoryHub()
// 	client := discovery.NewSocketClient(hub, nil)
//
// 	ctx, cancel := context.WithCancel(context.Background())
//
// 	go func(ctx context.Context, client *discovery.SocketClient) {
// 		client.ReadPipe(ctx)
// 	}(ctx, &client)
//
// 	hub.Register() <- &client
// 	hub.Broadcaster() <- []byte("hello")
// }
