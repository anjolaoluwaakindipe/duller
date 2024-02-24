package discovery

import "fmt"

type hub struct {
	socketClients map[*socketClient]bool
	broadcaster   chan []byte
	register      chan *socketClient
	unregister    chan *socketClient
}

func newHub() hub {
	return hub{socketClients: make(map[*socketClient]bool), broadcaster: make(chan []byte), register: make(chan *socketClient), unregister: make(chan *socketClient)}
}

func (h *hub) removeClient(client *socketClient) {
	delete(h.socketClients, client)
	close(client.send)
}

func (h *hub) run() {
	for {
		select {
		case client := <-h.register:
			h.socketClients[client] = true
			fmt.Println(len(h.socketClients))
		case client := <-h.unregister:
			h.removeClient(client)
			fmt.Println("client removed")
		case message := <-h.broadcaster:
			fmt.Println(string(message))
			fmt.Println(len(h.socketClients))
			for client := range h.socketClients {
				select {
				case client.send <- message:
					fmt.Println("sending to clients")
				default:
					h.removeClient(client)
				}
			}
		}
	}
}
