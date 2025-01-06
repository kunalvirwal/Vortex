package types

type ServiceState struct {
	Deployment   string
	Service      Service
	ContainerIDs []string
	HealthCheck  HealthCheck
}

type State []ServiceState
