package discovery

type socketClient struct {
	conn          *websocket.Conn
	send          chan []byte
	hub           hub
	writeWaitTime time.Duration
}

func newSocketClient(h hub, conn *websocket.Conn, opts ...socketClientOpts) socketClient {
	client := socketClient{hub: h, conn: conn, send: make(chan []byte, 256)}

	for _, opt := range opts {
		opt(&client)
	}

	return client
}

type socketClientOpts func(socketClient *socketClient)

func WithWriteWaitTime(duration time.Duration) socketClientOpts {
	return func(socketClient *socketClient) {
		socketClient.writeWaitTime = duration
	}
}

func (sc *socketClient) readPipe() {
	defer func() {
		sc.hub.unregister <- sc
		sc.conn.Close()
	}()
	ticker := time.NewTicker(time.Second * 2)

	for {
		select {
		case message, ok := <-sc.send:
			sc.conn.SetWriteDeadline(time.Now().Add(sc.writeWaitTime))
			if !ok {
				// The hub closed the channel.
				sc.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := sc.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			sc.conn.SetWriteDeadline(time.Now().Add(sc.writeWaitTime))
			if err := sc.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
