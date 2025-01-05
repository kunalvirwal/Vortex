package types

type HealthCheck struct {
	Command  string `json:"command" yaml:"command" valid:"matches(^.*\\S+.*$)~Command cannot contain only spaces"`
	Interval int    `json:"interval" yaml:"interval" valid:"numeric"`
	Timeout  int    `json:"timeout" yaml:"timeout" valid:"numeric"`
	Retries  int    `json:"retries" yaml:"retries" valid:"numeric"`
}

type ContainerConfig struct {
	Image       string
	Name        string
	ID          string
	Service     string
	Deployment  string
	Env         map[string]interface{}
	HealthCheck *HealthCheck
}
