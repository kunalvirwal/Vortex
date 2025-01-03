package master

import (
	"github.com/kunalvirwal/Vortex/internal/utils"
	"github.com/kunalvirwal/Vortex/types"
)

func ReplaceDiedContainer(cfg *types.ContainerConfig) {
	// replace the container
	cfg.Name = utils.GenerateContainerName(cfg)
}
