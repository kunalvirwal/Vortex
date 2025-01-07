package master

import (
	"fmt"
	"time"

	"github.com/docker/docker/api/types/events"
	"github.com/kunalvirwal/Vortex/internal/dockmaster"
	"github.com/kunalvirwal/Vortex/internal/state"
	"github.com/kunalvirwal/Vortex/internal/utils"
	"github.com/kunalvirwal/Vortex/types"
)

func ContainerDied(event events.Message) {
	containerID := event.Actor.ID
	container, tracked := IsVortexContainer(containerID)
	if tracked {
		vService, err := utils.GetServiceByName(container.Service)
		if vService.Service.RestartPolicy == "Always" {
			if err != nil {
				fmt.Println("Service not found!")
				return
			}
			RecordCrash(event, container)
			if !container.CrashData.IsInCrashLoop {
				if container.CrashData.CrashCount <= state.MaxCrashCount {
					dockmaster.ReplaceDiedContainer(container)
				} else {
					s := utils.NewScheduler(1, container.CrashData.CurrentBackoffDuration, func() { dockmaster.ReplaceDiedContainer(container) })
					s.Run()
				}
			}

		}
	}
}

func ContainerStarted(containerID string) {

}

func ContainerStopped(containerID string) {

}

func IsVortexContainer(containerID string) (*types.ContainerConfig, bool) {
	state.VortexContainers.Mu.RLock()
	defer state.VortexContainers.Mu.RUnlock()

	for _, container := range state.VortexContainers.List {
		if container.ID == containerID {
			return container, true
		}
	}
	return nil, false
}

func RecordCrash(event events.Message, container *types.ContainerConfig) {
	crashevent := types.CrashEvent{
		Timestamp: time.Unix(event.Time, 0),
		ExitCode:  event.Actor.Attributes["exitCode"],
		Signal:    event.Actor.Attributes["signal"],
		ExitError: event.Actor.Attributes["error"],
	}

	container.CrashData.Mu.Lock()
	defer container.CrashData.Mu.Unlock()

	container.CrashData.CrashCount++
	container.CrashData.LastCrashTime = time.Unix(event.Time, 0)
	container.CrashData.CrashHistory = append(container.CrashData.CrashHistory, crashevent)

	// if already surpassed max backoff duration, then the container is in crash loop
	if container.CrashData.CurrentBackoffDuration == state.MaxBackoffDuration {
		container.CrashData.IsInCrashLoop = true
		return
	}

	// if the container is not in crash loop, then update the backoff duration
	if container.CrashData.CurrentBackoffDuration == time.Duration(0) && container.CrashData.CrashCount > state.MaxCrashCount {
		container.CrashData.CurrentBackoffDuration = state.InitialBackoffDuration
	} else {
		container.CrashData.CurrentBackoffDuration = container.CrashData.CurrentBackoffDuration * time.Duration(state.BackoffMultiplier)
		if container.CrashData.CurrentBackoffDuration >= state.MaxBackoffDuration {
			container.CrashData.CurrentBackoffDuration = state.MaxBackoffDuration
		}
	}

}
