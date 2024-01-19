package pipeline

type Step struct {
	Function  string         `yaml:"function,omitempty"`
	ConfigMap map[string]any `yaml:"configMap,omitempty"`
}
