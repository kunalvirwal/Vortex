package types

import "time"

type HealthCheck struct {
	Command  string
	Interval time.Duration
	Timeout  time.Duration
	Retries  int
}

type ContainerConfig struct {
	Image       string
	Name        string
	ID          string
	Service     string
	Env         map[string]interface{}
	HealthCheck HealthCheck
}
