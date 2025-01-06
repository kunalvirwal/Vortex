package utils

import (
	"errors"

	"github.com/kunalvirwal/Vortex/internal/state"
	"github.com/kunalvirwal/Vortex/types"
)

func PopServiceByName(slice *state.ServiceList, name string) *types.VService {
	var rmService *types.VService
	slice.Mu.Lock()
	for i, v := range slice.List {
		if v.Service.Name == name {
			rmService = slice.List[i]
			slice.List = append(slice.List[:i], slice.List[i+1:]...)
		}
	}
	slice.Mu.Unlock()
	return rmService
}

func PopDeployment(deployments *state.DeploymentList, version string) *types.Deployment {
	var rmdep *types.Deployment
	for i, dep := range deployments.List {
		if dep.Version == version {
			deployments.Mu.Lock()
			rmdep = deployments.List[i]
			deployments.List = append(deployments.List[:i], deployments.List[i+1:]...)
			deployments.Mu.Unlock()
		}
	}
	return rmdep
}

func PopServicesByDepVersion(Vservices *state.ServiceList, version string) []*types.VService {
	var rmServices []*types.VService
	for i, service := range Vservices.List {
		if service.Deployment == version {
			rmServices = append(rmServices, Vservices.List[i])
		}
	}
	for _, service := range rmServices {
		Vservices.Mu.Lock()
		for i, v := range Vservices.List {
			if v.Service.Name == service.Service.Name && v.Deployment == service.Deployment {
				Vservices.List = append(Vservices.List[:i], Vservices.List[i+1:]...)
			}
		}
		Vservices.Mu.Unlock()
	}
	return rmServices
}

func RemoveContainerConfigsByService(Vcontainers *state.ContainerList, service *types.VService) {
	Vcontainers.Mu.Lock()
	containerIDs := service.ContainerIDs
	for _, id := range containerIDs {
		for j, cfg := range Vcontainers.List {
			if id == cfg.ID {
				Vcontainers.List = append(Vcontainers.List[:j], Vcontainers.List[j+1:]...)
				break
			}
		}
	}
	Vcontainers.Mu.Unlock()

}

func GetServiceByName(name string) (*types.VService, error) {
	state.VortexServices.Mu.RLock()
	defer state.VortexServices.Mu.RUnlock()
	for _, service := range state.VortexServices.List {
		if service.Service.Name == name {
			return service, nil
		}
	}
	return nil, errors.New("service not found")
}

func RemoveContainerConfigByID(Vcontainers []*types.ContainerConfig, id string) []*types.ContainerConfig {
	for j, cfg := range Vcontainers {
		if id == cfg.ID {
			Vcontainers = append(Vcontainers[:j], Vcontainers[j+1:]...)
			break
		}
	}
	return Vcontainers
}
