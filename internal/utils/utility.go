package utils

import (
	"github.com/kunalvirwal/Vortex/types"
)

func PopServiceByName(slice []*types.VService, name string) ([]*types.VService, *types.VService) {
	var rmService *types.VService
	for i, v := range slice {
		if v.Service.Name == name {
			rmService = slice[i]
			slice = append(slice[:i], slice[i+1:]...)
		}
	}
	return slice, rmService
}

func PopDeployment(deployments []*types.Deployment, version string) ([]*types.Deployment, *types.Deployment) {
	var rmdep *types.Deployment
	for i, dep := range deployments {
		if dep.Version == version {
			rmdep = deployments[i]
			deployments = append(deployments[:i], deployments[i+1:]...)
		}
	}
	return deployments, rmdep
}

func PopServicesByDepVersion(Vservices []*types.VService, version string) ([]*types.VService, []*types.VService) {
	var rmServices []*types.VService
	for i, service := range Vservices {
		if service.Deployment == version {
			rmServices = append(rmServices, Vservices[i])
		}
	}
	for _, service := range rmServices {
		Vservices, _ = PopServiceByName(Vservices, service.Service.Name)
	}
	return Vservices, rmServices
}

func RemoveContainerConfigByService(Vcontainers []*types.ContainerConfig, service *types.VService) []*types.ContainerConfig {

	containerIDs := service.ContainerIDs
	for _, id := range containerIDs {
		for j, cfg := range Vcontainers {
			if id == cfg.ID {
				Vcontainers = append(Vcontainers[:j], Vcontainers[j+1:]...)
				break
			}
		}
	}

	return Vcontainers
}
