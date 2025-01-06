package types

type State []ServiceState

type ServiceState struct {
	Deployment   string
	Service      Service
	ContainerIDs []Cntr
	HealthCheck  HealthCheck
}

type Cntr struct {
	ID         string
	ServiceUID uint
}
