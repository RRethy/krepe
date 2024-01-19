package pkg

type Imports struct {
	Files    []string         `yaml:"files,omitempty"`
	Packages []*PackageImport `yaml:"packages,omitempty"`
}

func (i *Imports) AddFile(file string) {
	i.Files = append(i.Files, file)
}

func (i *Imports) AddPackage(pkg *PackageImport) {
	i.Packages = append(i.Packages, pkg)
}
