package service

import "time"

type ServiceInfo struct {
	LastHeartbeat time.Time `json:"lastHearbeat"`
	ServiceId     string    `json:"serviceId"`
	IP            string    `json:"ip"`
	Port          string    `jsonn:"port"`
	Path          string    `jsonn:"path"`
	IsHealthy     bool      `jsonn:"isHealthy"`
	CurrentUse    int       `jsonn:"-"`
	WeightedUse   int       `jsonn:"weightedUse,omitempty"`
}
