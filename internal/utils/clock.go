package utils

import "time"

type Clock interface {
	Now() time.Time
}

type AppTime struct{}

func (at AppTime) Now() time.Time {
	return time.Now()
}

func NewClock() Clock {
	return AppTime{}
}
