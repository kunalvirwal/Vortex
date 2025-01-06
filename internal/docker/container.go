package docker

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/kunalvirwal/Vortex/internal/state"
	"github.com/kunalvirwal/Vortex/types"
)

func CreateContainer(cfg *types.ContainerConfig) error {
	// create container
	ctx := context.Background()

	// mount volume
	// port mapping
	// env
	envArray := []string{}
	for key, value := range cfg.Env {
		envArray = append(envArray, key+"="+fmt.Sprintf("%v", value))
	}
	// envArray := []string{"PORT=80", "SECRET_KEY=123456"}

	var Healthcheck *container.HealthConfig

	if cfg.HealthCheck != nil {
		Healthcheck = &container.HealthConfig{
			Test:     []string{"CMD-SHELL", cfg.HealthCheck.Command},
			Interval: time.Duration(cfg.HealthCheck.Interval) * time.Second,
			Timeout:  time.Duration(cfg.HealthCheck.Timeout) * time.Second,
			Retries:  cfg.HealthCheck.Retries,
		}
	} else {
		Healthcheck = nil
	}

	containerConfig := &container.Config{
		Image: cfg.Image,
		// ExposedPorts
		Env: envArray,
		// volume
		Healthcheck: Healthcheck,
	}

	hostConfig := &container.HostConfig{
		//port mapping
		//resources
	}

	containerConf, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, cfg.Name)
	if err != nil {
		return errors.New("Error in creating container: " + err.Error())
	}

	// appends its containerID to its VortexService.ContainerIDs
	cfg.ID = containerConf.ID
	for _, service := range state.VortexServices {
		if service.Service.Name == cfg.Service && service.Deployment == cfg.Deployment {
			service.Mu.Lock()
			service.ContainerIDs = append(service.ContainerIDs, cfg.ID)
			service.Mu.Unlock()
		}
	}

	// appends the container to VortexContainers
	state.VortexContainers = append(state.VortexContainers, cfg)
	return nil
}

func DeleteContainer(id string) error {
	fmt.Println("Deleting container: ", id)
	err := cli.ContainerRemove(context.Background(), id, container.RemoveOptions{
		Force: true,
	})
	if err != nil {
		return errors.New("Error in deleting container: " + err.Error())
	}
	return nil
}

func StartContainer(containerID string) error {
	ctx := context.Background()
	err := cli.ContainerStart(ctx, containerID, container.StartOptions{})
	if err != nil {
		return errors.New("Error in starting container: " + err.Error())
	}
	return nil
}
