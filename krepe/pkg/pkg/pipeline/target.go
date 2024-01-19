package pipeline

type Target struct {
	ApiVersion string `yaml:"apiVersion,omitempty"`
	Group      string `yaml:"group,omitempty"`
	Version    string `yaml:"version,omitempty"`
	Kind       string `yaml:"kind,omitempty"`
	Name       string `yaml:"name,omitempty"`
	Namespace  string `yaml:"namespace,omitempty"`
}
