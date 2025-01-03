package docker

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/kunalvirwal/Vortex/types"
)

func CreateContainer(cfg types.ContainerConfig) (string, error) {
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

	containerConfig := &container.Config{
		Image: cfg.Image,
		// ExposedPorts
		Env: envArray,
		// volume
		Healthcheck: &container.HealthConfig{
			Test:     []string{"CMD-SHELL", cfg.HealthCheck.Command},
			Interval: cfg.HealthCheck.Interval * time.Second,
			Timeout:  cfg.HealthCheck.Timeout * time.Second,
			Retries:  cfg.HealthCheck.Retries,
		},
	}

	hostConfig := &container.HostConfig{
		//port mapping
		//resources
	}

	containerConf, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, cfg.Name)
	if err != nil {
		fmt.Println("Error in creating container", err)
	}

	return containerConf.ID, nil
}
