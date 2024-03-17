package discovery

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// socketClient is a struct object responsible
// for handling client connections involved in a hub
type socketClient struct {
	Id            uuid.UUID
	conn          *websocket.Conn
	send          chan []byte
	hub           *hub
	writeWaitTime time.Duration
}

// newSocketClient is a constructor function for a new socketClient
func newSocketClient(h *hub, conn *websocket.Conn, opts ...socketClientOpts) socketClient {
	client := socketClient{Id: uuid.New(), hub: h, conn: conn, send: make(chan []byte, 256)}

	for _, opt := range opts {
		opt(&client)
	}

	return client
}

// socketClientOpts is an option function for creating
// new socketClients
type socketClientOpts func(socketClient *socketClient)

// WithWriteWaitTimes is a socketClient option that sets the
// writeWaitTime for a socketClient upon creation
func WithWriteWaitTime(duration time.Duration) socketClientOpts {
	return func(socketClient *socketClient) {
		socketClient.writeWaitTime = duration
	}
}

// readPipe facilitates reading messages broadcasted from the hub and sending
// messages through the websocket connection to the user
func (sc *socketClient) readPipe() {
	ticker := time.NewTicker(time.Second * 2)

	defer func() {
		sc.hub.unregister <- sc
		ticker.Stop()
		sc.conn.Close()
	}()

	for {
		select {
		case message, ok := <-sc.send:
			sc.conn.SetWriteDeadline(time.Now().Add(sc.writeWaitTime))
			if !ok {
				sc.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := sc.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				slog.Warn(fmt.Sprintf("Could not establish a text message connection %v", sc.conn))
				return
			}
			w.Write(message)

			n := len(sc.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-sc.send)
			}

			if err := w.Close(); err != nil {
				slog.Warn(fmt.Sprintf("Could not close text message connection %v", sc.conn))
				return
			}

		case <-ticker.C:
			sc.conn.SetWriteDeadline(time.Now().Add(sc.writeWaitTime))
			if err := sc.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				slog.Warn(fmt.Sprintf("Could not ping client with websocket %v", sc.Id))
				return
			}
		}
	}
}
