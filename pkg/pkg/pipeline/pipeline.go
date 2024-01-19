package pipeline

type Pipeline struct {
	Steps []Step `yaml:"steps,omitempty"`
}
