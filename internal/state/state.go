package state

import (
	"sync"

	"github.com/kunalvirwal/Vortex/types"
)

type DeploymentList struct {
	Mu   sync.RWMutex
	List []*types.Deployment
}

type ServiceList struct {
	Mu   sync.RWMutex
	List []*types.VService
}

type ContainerList struct {
	Mu   sync.RWMutex
	List []*types.ContainerConfig
}

var VortexDeployments = &DeploymentList{}
var VortexServices = &ServiceList{}
var VortexContainers = &ContainerList{}

func GetState() types.State {

	var state types.State
	var h types.HealthCheck
	var s []types.Cntr

	VortexServices.Mu.RLock()
	for _, service := range VortexServices.List {
		if service.Service.HealthCheck == nil {
			h = types.HealthCheck{
				Command:  "",
				Interval: 0,
				Timeout:  0,
				Retries:  0,
			}
		} else {
			h = *service.Service.HealthCheck
		}

		for _, container := range VortexContainers.List {
			if container.Service == service.Service.Name {
				s = append(s, types.Cntr{
					ID:         container.ID,
					ServiceUID: container.ServiceUID,
				})
			}
		}

		state = append(state, types.ServiceState{
			Deployment:   service.Deployment,
			Service:      service.Service,
			ContainerIDs: s,
			HealthCheck:  h,
		})
	}
	VortexServices.Mu.RUnlock()
	// fmt.Println(state)
	return state
}
