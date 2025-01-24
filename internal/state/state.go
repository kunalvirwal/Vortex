package state

import (
	"sync"
	"time"

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

// State variables
var VortexDeployments = &DeploymentList{}
var VortexServices = &ServiceList{}
var VortexContainers = &ContainerList{}

// ExponentialCrashBackoff variables
var MaxCrashCount uint = 2                                 // Maximum number of crashes allowed before exponential backoff kicks in
var MaxBackoffDuration time.Duration = time.Minute * 2     // Upperlimit for exponential backoff duration
var InitialBackoffDuration time.Duration = time.Second * 1 // Initial backoff duration
var BackoffMultiplier uint = 2                             // Exponential Multiplier for backoff
var CrashHistorySaveDuration time.Duration = 0             // 0 means save crash events forever, can take any time value
var BackoffResetDuration time.Duration = time.Minute * 5   // Reset the backoff duration, if no crashes for this duration after the last crash

func GetState() types.State {

	var state types.State

	VortexServices.Mu.RLock()
	for _, service := range VortexServices.List {
		var h types.HealthCheck
		var s []types.Cntr
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

		VortexContainers.Mu.RLock()
		// for _, id := range service.ContainerIDs {
		// 	fmt.Println(id)
		// }
		for _, id := range service.ContainerIDs {
			for _, container := range VortexContainers.List {
				if container.ID == id {
					s = append(s, types.Cntr{
						ID:               container.ID,
						ServiceUID:       container.ServiceUID,
						CrashLoopBackOff: container.CrashData.IsInCrashLoop,
					})
					break
				}
			}
		}
		VortexContainers.Mu.RUnlock()

		state = append(state, types.ServiceState{
			Deployment:   service.Deployment,
			Service:      *service.Service,
			ContainerIDs: s,
			HealthCheck:  h,
		})
	}
	VortexServices.Mu.RUnlock()
	// fmt.Println(state)
	return state
}
