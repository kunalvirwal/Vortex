package utils

import "github.com/kunalvirwal/Vortex/types"

func RemoveByServiceName(slice []*types.VService, elem *types.VService) []*types.VService {
	for i, v := range slice {
		if v.Service.Name == elem.Service.Name {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}
