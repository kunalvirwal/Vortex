package state

import (
	"github.com/kunalvirwal/Vortex/types"
)

var VortexDeployments []*types.Deployment
var VortexServices []*types.VService
var VortexContainers []*types.ContainerConfig

func GetState() types.State {
	var state types.State
	var h types.HealthCheck
	for _, service := range VortexServices {
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
		state = append(state, types.ServiceState{
			Deployment:   service.Deployment,
			Service:      service.Service,
			ContainerIDs: service.ContainerIDs,
			HealthCheck:  h,
		})
	}
	// fmt.Println(state)
	return state
}
