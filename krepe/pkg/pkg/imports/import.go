package imports

type Imports struct {
	Files    []string `yaml:"files,omitempty"`
	Packages []*Pkg   `yaml:"packages,omitempty"`
}

func (i *Imports) AddFile(file string) {
	i.Files = append(i.Files, file)
}

func (i *Imports) AddPackage(pkg *Pkg) {
	i.Packages = append(i.Packages, pkg)
}
