package discovery

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type WSConn interface {
	Close() error
	// NextWriter returns a writer for the next message to send. The writer's Close method flushes
	// the complete message to the network.
	// There can be at most one open writer on a connection. NextWriter closes the previous writer
	// if the application has not already done so.
	// All message types (TextMessage, BinaryMessage, CloseMessage, PingMessage and PongMessage) are supported.
	NextWriter(messageType int) (io.WriteCloser, error)
	SetWriteDeadline(t time.Time) error
	// WriteMessage is a helper method for getting a writer using NextWriter,
	// writing the message and closing the writer.
	WriteMessage(messageType int, data []byte) error
}

// SocketClient is a struct object responsible
// for handling client connections involved in a hub
type SocketClient struct {
	Id            uuid.UUID
	conn          WSConn
	send          chan []byte
	hub           Hub
	writeWaitTime time.Duration
}

// NewSocketClient is a constructor function for a new socketClient
func NewSocketClient(h Hub, conn WSConn, opts ...socketClientOpts) SocketClient {
	client := SocketClient{Id: uuid.New(), hub: h, conn: conn, send: make(chan []byte, 256)}

	for _, opt := range opts {
		opt(&client)
	}

	return client
}

func (sc *SocketClient) Send() chan []byte {
	return sc.send
}

// socketClientOpts is an option function for creating
// new socketClients
type socketClientOpts func(socketClient *SocketClient)

// WithWriteWaitTimes is a socketClient option that sets the
// writeWaitTime for a socketClient upon creation
func WithWriteWaitTime(duration time.Duration) socketClientOpts {
	return func(socketClient *SocketClient) {
		socketClient.writeWaitTime = duration
	}
}

// ReadPipe facilitates reading messages broadcasted from the hub and sending
// messages through the websocket connection to the user
func (sc *SocketClient) ReadPipe(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 2)

	defer func() {
		sc.hub.Unregister() <- sc
		ticker.Stop()
		if sc.conn != nil {
			sc.conn.Close()
		}
	}()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Cancel was called")
			return

		case message, ok := <-sc.Send():
			if sc.conn == nil {
				return
			}
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
			fmt.Println("hello")
			for i := 0; i < n; i++ {
				fmt.Println("hello2")
				w.Write([]byte{'\n'})
				w.Write(<-sc.send)
			}

			if err := w.Close(); err != nil {
				slog.Warn(fmt.Sprintf("Could not close text message connection %v", sc.conn))
				return
			}
			fmt.Println("hello3")
			fmt.Println(string(message))

		case <-ticker.C:
			sc.conn.SetWriteDeadline(time.Now().Add(sc.writeWaitTime))
			if err := sc.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				slog.Warn(fmt.Sprintf("Could not ping client with websocket %v", sc.Id))
				return
			}

		}
	}
}
