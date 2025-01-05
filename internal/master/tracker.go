package master

import (
	"github.com/kunalvirwal/Vortex/internal/dockmaster"
	"github.com/kunalvirwal/Vortex/internal/state"
	"github.com/kunalvirwal/Vortex/types"
)

func ContainerDied(containerID string) {
	container, tracked := IsVortexContainer(containerID)
	if tracked {
		dockmaster.ReplaceDiedContainer(container)
	}
}

func ContainerStarted(containerID string) {

}

func ContainerStopped(containerID string) {

}

func IsVortexContainer(containerID string) (*types.ContainerConfig, bool) {
	for _, container := range state.VortexContainers {
		if container.ID == containerID {
			return container, true
		}
	}
	return nil, false
}
