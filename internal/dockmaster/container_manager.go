package dockmaster

import (
	"fmt"
	"strings"

	"github.com/kunalvirwal/Vortex/internal/docker"
	"github.com/kunalvirwal/Vortex/internal/state"
	"github.com/kunalvirwal/Vortex/types"
	"golang.org/x/sync/errgroup"
)

func ReplaceDiedContainer(cfg *types.ContainerConfig) {
	// replace the container
	newName := docker.GenerateContainerName(cfg)
	docker.FindOrPullImage(cfg.Image)

	newContainer := &types.ContainerConfig{
		Name:         newName,
		Service:      cfg.Service,
		Deployment:   cfg.Deployment,
		Image:        cfg.Image,
		HealthCheck:  cfg.HealthCheck,
		Env:          cfg.Env,
		ServiceUID:   cfg.ServiceUID,
		StartCommand: cfg.StartCommand,
		CrashData: types.CrashData{
			LastCrashTime:          cfg.CrashData.LastCrashTime,
			CrashCount:             cfg.CrashData.CrashCount,
			CurrentBackoffDuration: cfg.CrashData.CurrentBackoffDuration,
			IsInCrashLoop:          cfg.CrashData.IsInCrashLoop,
			CrashHistory:           cfg.CrashData.CrashHistory,
		},
	}

	err := docker.CreateContainer(newContainer)
	if err != nil {
		fmt.Println("Error in creating container: " + err.Error())
		return
	}

	// Update the containerID in VortexService ContainerIDs
	for i := range state.VortexServices.List {
		if state.VortexServices.List[i].Service.Name == cfg.Service && state.VortexServices.List[i].Deployment == cfg.Deployment {
			for j := range state.VortexServices.List[i].ContainerIDs {
				if state.VortexServices.List[i].ContainerIDs[j] == cfg.ID {
					state.VortexServices.List[i].Mu.Lock()
					state.VortexServices.List[i].ContainerIDs[j] = newContainer.ID
					state.VortexServices.List[i].Mu.Unlock()
					break
				}
			}
			break
		}
	}

	// Remove the old container from VortexContainers
	for i := range state.VortexContainers.List {
		if state.VortexContainers.List[i].ID == cfg.ID {
			state.VortexContainers.Mu.Lock()
			state.VortexContainers.List = append(state.VortexContainers.List[:i], state.VortexContainers.List[i+1:]...)
			state.VortexContainers.Mu.Unlock()
			break
		}
	}

	err = docker.DeleteContainer(cfg.ID)
	if err != nil {
		if strings.Contains(err.Error(), "No such container") {
			fmt.Println("Container already deleted : " + cfg.ID)
		} else {
			fmt.Println("Error in deleting container: " + err.Error())
			return
		}
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
			Service:      VService.Service.Name,
			Deployment:   VService.Deployment,
			Image:        VService.Service.Image,
			HealthCheck:  VService.Service.HealthCheck,
			Env:          VService.Service.Env,
			ServiceUID:   uint(i) + 1,
			StartCommand: VService.Service.StartCommand,
		}
		container.Name = docker.GenerateContainerName(container)
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
