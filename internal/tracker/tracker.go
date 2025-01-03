package tracker

import (
	"github.com/kunalvirwal/Vortex/internal/master"
	"github.com/kunalvirwal/Vortex/types"
)

var VortexContainers *[]types.ContainerConfig

func ContainerDied(containerID string) {
	container, tracked := IsVortexContainer(containerID)
	if tracked {
		master.ReplaceDiedContainer(container)
	}
}

func ContainerStarted(containerID string) {

}

func ContainerStopped(containerID string) {

}

func IsVortexContainer(containerID string) (*types.ContainerConfig, bool) {
	for _, container := range *VortexContainers {
		if container.ID == containerID {
			return &container, true
		}
	}
	return nil, false
}
