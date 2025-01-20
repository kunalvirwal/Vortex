package dockmaster

import (
	"fmt"
	"math"
	"strings"

	"github.com/kunalvirwal/Vortex/internal/docker"
	"github.com/kunalvirwal/Vortex/internal/state"
	"github.com/kunalvirwal/Vortex/internal/utils"
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

	// Remove the old container from VortexContainers
	fmt.Println("Removing old container from VortexContainers")
	state.VortexContainers.Mu.Lock()
	for i := range state.VortexContainers.List {
		if state.VortexContainers.List[i].ID == cfg.ID {
			state.VortexContainers.List = append(state.VortexContainers.List[:i], state.VortexContainers.List[i+1:]...)
			break
		}
	}
	state.VortexContainers.Mu.Unlock()

	fmt.Println("Creating new container")
	// Create new container
	err := docker.CreateContainer(newContainer)
	if err != nil {
		fmt.Println("Error in creating container: " + err.Error())
		return
	}

	fmt.Println("Updating vortex service.ContainerIDS")
	// Update the containerID in VortexService ContainerIDs
	for i := range state.VortexServices.List {
		if state.VortexServices.List[i].Service.Name == cfg.Service && state.VortexServices.List[i].Deployment == cfg.Deployment {
			for j := range state.VortexServices.List[i].ContainerIDs {
				if state.VortexServices.List[i].ContainerIDs[j] == cfg.ID {

					state.VortexServices.List[i].Mu.Lock()
					state.VortexServices.List[i].ContainerIDs = append(state.VortexServices.List[i].ContainerIDs[:j], state.VortexServices.List[i].ContainerIDs[j+1:]...)
					state.VortexServices.List[i].Mu.Unlock()
					break
				}
			}
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

	fmt.Println("Starting the new container")
	err = docker.StartContainer(newContainer.ID)
	if err != nil {
		fmt.Println("Error in starting container: " + err.Error())
	} else {
		fmt.Println("Creating new container and starting its backoff reset scheduler for :" + newContainer.ID)
		s := utils.NewScheduler(1, state.BackoffResetDuration, func() {
			utils.ResetCrashbackOff(newContainer.ID)
		})
		s.RunAsync()
		utils.BackOffResetSchedulers.Mu.Lock()
		utils.BackOffResetSchedulers.Schedule[newContainer.ID] = s
		utils.BackOffResetSchedulers.Mu.Unlock()
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

	VService.Mu.Lock()
	defer VService.Mu.Unlock()
	fmt.Println("got inside modify")
	var serviceUpdated, containerUpdated bool

	if VService.Service.Name != UpdatedService.Service.Name {
		return fmt.Errorf("service can not be modified")
	}

	if VService.Service.Replicas != UpdatedService.Service.Replicas {
		serviceUpdated = true
	}

	if VService.Service.RestartPolicy != UpdatedService.Service.RestartPolicy {
		VService.Service.RestartPolicy = UpdatedService.Service.RestartPolicy
		serviceUpdated = true
	}

	if VService.Service.Image != UpdatedService.Service.Image || compareEnvs(VService.Service.Env, UpdatedService.Service.Env) || compareHealthChecks(VService.Service.HealthCheck, UpdatedService.Service.HealthCheck) || VService.Service.StartCommand != UpdatedService.Service.StartCommand {
		VService.Service.Image = UpdatedService.Service.Image
		VService.Service.Env = UpdatedService.Service.Env
		VService.Service.HealthCheck = UpdatedService.Service.HealthCheck
		VService.Service.StartCommand = UpdatedService.Service.StartCommand
		containerUpdated = true
		// kill and restart
	}

	if containerUpdated {
		// kill and restart whole service as changes are immutable without restarting
		fmt.Println("Container configs updated")
		for i := range VService.ContainerIDs {
			err := docker.DeleteContainer(VService.ContainerIDs[i])
			if err != nil {
				fmt.Println("Error in deleting container: " + err.Error())
			}
			utils.RemoveContainerConfigsByService(state.VortexContainers, VService)
		}
		Deploy(VService)
	} else if serviceUpdated {
		fmt.Println("Service configs updated", VService.Service.Replicas, UpdatedService.Service.Replicas)
		diff := VService.Service.Replicas - UpdatedService.Service.Replicas
		if diff > 0 {
			// scale down
			fmt.Println(diff)
			ScaleDown(VService, int(math.Abs(float64(diff))))
			fmt.Println("Scaled down")
		}
		if diff < 0 {
			// scale up
		}
	} else {
		fmt.Println("No changes in service", VService.Service.Name, UpdatedService.Service.Name)
	}

	return nil
}

func ScaleDown(service *types.VService, diff int) error {
	fmt.Println("Scaling down service")
	fmt.Println(service.ContainerIDs)
	oldQuantity := service.Service.Replicas
	if oldQuantity < diff {
		return fmt.Errorf("can not scale down more than the current replicas")
	}
	// ServiceUID range to delete

	high := uint(oldQuantity)       //inclusive
	low := uint(oldQuantity - diff) //non inclusive

	fmt.Println("High: ", high)
	fmt.Println("Low: ", low)
	state.VortexContainers.Mu.Lock()
	var g errgroup.Group
	fmt.Println("helo")

	//////////////////////////////some thing wrong here v
	for _, containerID := range service.ContainerIDs {
		for _, cfg := range state.VortexContainers.List {
			fmt.Println("ServiceUID: ", cfg.ServiceUID, cfg.ID == containerID, cfg.ServiceUID <= high, cfg.ServiceUID > low)
			if cfg.ID == containerID && cfg.ServiceUID <= high && cfg.ServiceUID > low {
				// removing from vortex containers
				state.VortexContainers.List = utils.RemoveContainerConfigByID(state.VortexContainers.List, cfg.ID)

				// removing from vortexService.ContainerID
				// service.Mu.Lock()

				for i, id := range service.ContainerIDs {
					if id == cfg.ID {
						service.ContainerIDs = append(service.ContainerIDs[:i], service.ContainerIDs[i+1:]...)
						break
					}
				}
				fmt.Println("Removed containerID from list")
				// service.Mu.Unlock()
				g.Go(func() error {
					fmt.Println("Removing container: " + cfg.ID)
					return docker.DeleteContainer(cfg.ID)
				})
				break
			}
		}
	}
	service.Service.Replicas -= diff
	state.VortexContainers.Mu.Unlock()
	if err := g.Wait(); err != nil {
		fmt.Println("Error in deleting container: " + err.Error())
	}
	return nil
}

// Returns truE if the provideD maps are different
func compareEnvs(env1 map[string]interface{}, env2 map[string]interface{}) bool {
	if len(env1) != len(env2) {
		return true
	}

	for i := range env1 {
		if env1[i] != env2[i] {
			return true
		}
	}
	return false
}

func compareHealthChecks(hc1 *types.HealthCheck, hc2 *types.HealthCheck) bool {

	if hc1 == nil && hc2 == nil {
		return false
	}

	if (hc1 == nil && hc2 != nil) || (hc1 != nil && hc2 == nil) {
		return true
	}

	if hc1.Command != hc2.Command || hc1.Interval != hc2.Interval || hc1.Retries != hc2.Retries || hc1.Timeout != hc2.Timeout {
		return true
	}

	return false
}
