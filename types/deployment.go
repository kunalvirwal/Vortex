package types

type Deployment struct {
	Version  string    `json:"version" yaml:"version" valid:"required~Field Version is required,matches(^.*\\S+.*$)~Version cannot contain only spaces"`
	Services []Service `json:"services" yaml:"services" valid:"required~Atleast 1 service is required"`
}
