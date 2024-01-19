package imports

type Imports struct {
	Files    []string `yaml:"files,omitempty"`
	Packages []*Pkg   `yaml:"packages,omitempty"`
}
