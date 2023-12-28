package discovery

import "time"

type ServiceInfo struct {
	LastHeartbeat time.Time
	ServiceId     string
	Address       string
	Path          string
}
