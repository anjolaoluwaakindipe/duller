package service

import "time"

type ServiceInfo struct {
	LastHeartbeat time.Time
	ServiceId     string
	IP            string
	Port          string
	Path          string
}
