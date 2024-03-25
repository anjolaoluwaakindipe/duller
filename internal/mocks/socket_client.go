package mocks

import (
	"io"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockWebSocket struct {
	mock.Mock
}

func (s *MockWebSocket) Close() error {
	args := s.Called()
	return args.Error(0)
}

func (s *MockWebSocket) NextWriter(messageType int) (io.WriteCloser, error) {
	args := s.Called(messageType)
	return args.Get(0).(io.WriteCloser), args.Error(1)
}

func (s *MockWebSocket) SetWriteDeadline(t time.Time) error {
	args := s.Called(t)
	return args.Error(0)
}

func (s *MockWebSocket) WriteMessage(messageType int, data []byte) error {
	args := s.Called(messageType, data)
	return args.Error(0)
}

type MockWriteCloser struct {
	mock.Mock
}

func (mwc *MockWriteCloser) Close() error {
	args := mwc.Called()
	return args.Error(0)
}

func (mwc *MockWriteCloser) Write(p []byte) (n int, err error) {
	args := mwc.Called(p)
	return args.Int(0), args.Error(1)
}
