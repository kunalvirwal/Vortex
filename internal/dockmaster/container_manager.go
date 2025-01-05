package dockmaster

import (
	"fmt"

	"github.com/kunalvirwal/Vortex/internal/docker"
	"github.com/kunalvirwal/Vortex/internal/utils"
	"github.com/kunalvirwal/Vortex/types"
	"golang.org/x/sync/errgroup"
)

func ReplaceDiedContainer(cfg *types.ContainerConfig) {
	// replace the container
	cfg.Name = utils.GenerateContainerName(cfg)
}

func Deploy(VService *types.VService) error {
	fmt.Println("Deploying Service")
	docker.FindOrPullImage(VService.Service.Image)
	var eg errgroup.Group
	for i := 0; i < VService.Service.Replicas; i++ {
		container := &types.ContainerConfig{
			Service:     VService.Service.Name,
			Deployment:  VService.Deployment,
			Image:       VService.Service.Image,
			HealthCheck: VService.Service.HealthCheck,
			Env:         VService.Service.Env,
		}
		container.Name = utils.GenerateContainerName(container)
		eg.Go(func() error {
			err := docker.CreateContainer(container)
			if err != nil {
				return err
			}
			return docker.StartContainer(container.ID)
		})

	}
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

func Modify(VService *types.VService, UpdatedService *types.VService) error {
	fmt.Println("Modify service") // also modify state.VortexServices
	return nil
}
