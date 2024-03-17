package service

import "time"

type ServiceInfo struct {
	LastHeartbeat time.Time `json:"lastHearbeat"`
	ServiceId     string    `json:"serviceId"`
	IP            string    `json:"ip"`
	Port          string    `json:"port"`
	Path          string    `json:"path"`
	IsHealthy     bool      `json:"isHealthy"`
	CurrentUse    int       `json:"-"`
	WeightedUse   int       `json:"weightedUse,omitempty"`
}
