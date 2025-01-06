package master

import (
	"fmt"

	"github.com/kunalvirwal/Vortex/internal/dockmaster"
	"github.com/kunalvirwal/Vortex/internal/state"
	"github.com/kunalvirwal/Vortex/internal/utils"
	"github.com/kunalvirwal/Vortex/types"
)

func ContainerDied(containerID string) {
	container, tracked := IsVortexContainer(containerID)
	if tracked {
		vService, err := utils.GetServiceByName(container.Service)
		if vService.Service.RestartPolicy == "Always" {
			if err != nil {
				fmt.Println("Service not found!")
				return
			}
			dockmaster.ReplaceDiedContainer(container)
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
