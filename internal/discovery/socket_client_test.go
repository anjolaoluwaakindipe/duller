package discovery_test

import (
	"context"
	"testing"
	"time"

	"github.com/anjolaoluwaakindipe/duller/internal/discovery"
	"github.com/anjolaoluwaakindipe/duller/internal/mocks"
	"github.com/stretchr/testify/mock"
)

func Test_SocketClient_ReadPipe(t *testing.T) {
	t.Run("Should write message from send channel", func(t *testing.T) {
		t.Parallel()
		hub := discovery.NewInMemoryHub()
		mockWs := new(mocks.MockWebSocket)
		socket := discovery.NewSocketClient(hub, mockWs)

		ctx, cancel := context.WithCancel(context.Background())

		textMessage := 1
		message := []byte("hello")

		writeCloser := new(mocks.MockWriteCloser)
		writeCloser.On("Write", message).Return(len(message), nil)
		mockWs.On("SetWriteDeadline", mock.Anything).Return(nil)
		mockWs.On("NextWriter", textMessage).Return(writeCloser, nil)

		go func(ctx context.Context) {
			socket.ReadPipe(ctx)
		}(ctx)

		socket.Send() <- message

		time.Sleep(100 * time.Millisecond)
		cancel()

		mockWs.AssertExpectations(t)
		mockWs.AssertExpectations(t)
		mockWs.AssertNumberOfCalls(t, "NextWriter", 1)
		writeCloser.AssertNumberOfCalls(t, "Write", 1)
	})
}
