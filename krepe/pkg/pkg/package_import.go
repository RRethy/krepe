package pkg

import (
	"github.com/RRethy/krepe/krepe/pkg/git"
	"github.com/RRethy/krepe/krepe/pkg/yaml"
)

var (
	_ yaml.InterfaceUnmarshaler = &PackageImport{}
	_ yaml.InterfaceMarshaler   = &PackageImport{}
)

type PackageImport struct {
	Ref  *git.PkgRef
	name string
}

type RawPkg struct {
	Ref  string `yaml:"ref,omitempty"`
	Name string `yaml:"name,omitempty"`
}

func NewPackageImport(ref *git.PkgRef, name string) *PackageImport {
	return &PackageImport{
		Ref:  ref,
		name: name,
	}
}

func (p *PackageImport) UnmarshalYAML(unmarshal func(interface{}) error) error {
	raw := RawPkg{}
	if err := unmarshal(&raw); err != nil {
		return err
	}

	ref, err := git.NewPkgRefFromString(raw.Ref)
	if err != nil {
		return err
	}

	newPkg := NewPackageImport(ref, raw.Name)
	if err != nil {
		return err
	}

	*p = *newPkg
	return nil
}

func (p *PackageImport) MarshalYAML() (interface{}, error) {
	return RawPkg{
		Ref:  p.Ref.String(),
		Name: p.name,
	}, nil
}

func (p *PackageImport) Name() string {
	if p.name == "" {
		return p.Ref.Repo
	}

	return p.name
}
