package discovery

import (
	"context"
	"fmt"
)

type Hub interface {
	Run(context.Context)
	Broadcaster() chan []byte
	// RegistrationSignal() chan *SocketClient
	// DeRegistrationSignal() chan *SocketClient
	Register() chan *SocketClient
	Unregister() chan *SocketClient
}

type InMemoryHub struct {
	SocketClients    map[*SocketClient]bool
	broadcaster      chan []byte
	register         chan *SocketClient
	unregister       chan *SocketClient
	registerSignal   chan *SocketClient
	deregisterSignal chan *SocketClient
}

// NewInMemoryHub instantiates an in memory implementation
// of Hub
func NewInMemoryHub() Hub {
	return &InMemoryHub{
		SocketClients: make(map[*SocketClient]bool),
		broadcaster:   make(chan []byte),
		register:      make(chan *SocketClient),
		unregister:    make(chan *SocketClient),
	}
}

// Broadcaster implements Hub.
func (h *InMemoryHub) Broadcaster() chan []byte {
	return h.broadcaster
}

// Register implements Hub.
func (h *InMemoryHub) Register() chan *SocketClient {
	return h.register
}

// Unregister implements Hub.
func (h *InMemoryHub) Unregister() chan *SocketClient {
	return h.unregister
}

func (h *InMemoryHub) removeClient(client *SocketClient) {
	delete(h.SocketClients, client)
	fmt.Println("hello")
	close(client.send)
}

func (h *InMemoryHub) Run(ctx context.Context) {
	for {
		select {
		case client := <-h.register:
			h.SocketClients[client] = true
		case client := <-h.unregister:
			h.removeClient(client)
		case message := <-h.broadcaster:
			for client := range h.SocketClients {
				select {
				case client.send <- message:
				default:
					h.removeClient(client)
				}
			}
		case <-ctx.Done():
			return
		}
	}
}
