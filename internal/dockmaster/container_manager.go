package dockmaster

import (
	"fmt"

	"github.com/kunalvirwal/Vortex/internal/docker"
	"github.com/kunalvirwal/Vortex/internal/state"
	"github.com/kunalvirwal/Vortex/internal/utils"
	"github.com/kunalvirwal/Vortex/types"
	"golang.org/x/sync/errgroup"
)

func ReplaceDiedContainer(cfg *types.ContainerConfig) {
	// replace the container
	newName := utils.GenerateContainerName(cfg)
	docker.FindOrPullImage(cfg.Image)
	newContainer := &types.ContainerConfig{
		Name:        newName,
		Service:     cfg.Service,
		Deployment:  cfg.Deployment,
		Image:       cfg.Image,
		HealthCheck: cfg.HealthCheck,
		Env:         cfg.Env,
	}
	err := docker.CreateContainer(newContainer)
	if err != nil {
		fmt.Println("Error in creating container: " + err.Error())
		return
	}

	for i := range state.VortexServices {
		if state.VortexServices[i].Service.Name == cfg.Service && state.VortexServices[i].Deployment == cfg.Deployment {
			for j := range state.VortexServices[i].ContainerIDs {
				if state.VortexServices[i].ContainerIDs[j] == cfg.ID {
					state.VortexServices[i].Mu.Lock()
					state.VortexServices[i].ContainerIDs[j] = newContainer.ID
					state.VortexServices[i].Mu.Unlock()
					break
				}
			}
			break
		}
	}
	for i := range state.VortexContainers {
		if state.VortexContainers[i].ID == cfg.ID {
			state.VortexContainers[i] = newContainer
			break
		}
	}

	err = docker.DeleteContainer(cfg.ID)
	if err != nil {
		fmt.Println("Error in deleting container " + cfg.ID + " : " + err.Error())
		return
	}

	err = docker.StartContainer(newContainer.ID)
	if err != nil {
		fmt.Println("Error in starting container: " + err.Error())
		return
	}

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
