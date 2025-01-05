package utils

import (
	"errors"

	"github.com/asaskevich/govalidator"
	"github.com/kunalvirwal/Vortex/types"
)

func ValidateDeployment(dep *types.Deployment) error {

	if result, err := govalidator.ValidateStruct(dep); !result {
		return err
	}

	for _, service := range dep.Services {

		if service.Replicas < 1 {
			return errors.New("invalid number of replicas in service " + service.Name)
		}

		if service.HealthCheck != nil {
			if service.HealthCheck.Interval < 1 || service.HealthCheck.Timeout < 1 || service.HealthCheck.Retries < 1 {
				return errors.New("invalid health check configuration in service " + service.Name)

			}
		}

		if service.RestartPolicy != "" && service.RestartPolicy != "Always" && service.RestartPolicy != "OnFailure" && service.RestartPolicy != "Never" {
			return errors.New("invalid restart policy in service " + service.Name)
		}
		if service.RestartPolicy == "" {
			service.RestartPolicy = "Always"
		}
	}

	return nil
}
