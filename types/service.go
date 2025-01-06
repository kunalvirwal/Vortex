package types

import "sync"

type Service struct {
	Name          string                 `json:"name" yaml:"name" valid:"required~Field Name is required,matches(^.*\\S+.*$)~Name cannot contain only spaces"`
	Image         string                 `json:"image" yaml:"image" valid:"required~Field Image Name is required,matches(^.*\\S+.*$)~Image cannot contain only spaces"`
	Replicas      int                    `json:"replicas" yaml:"replicas" valid:"required~Field Replicas is required,numeric"`
	Env           map[string]interface{} `json:"env" yaml:"env"`
	HealthCheck   *HealthCheck           `json:"healthCheck" yaml:"healthCheck"`
	RestartPolicy string                 `json:"restartPolicy" yaml:"restartPolicy"`
}

type VService struct {
	Service      Service
	Deployment   string
	ContainerIDs []string
	Mu           sync.RWMutex
}
