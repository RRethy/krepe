package imports

import (
	"github.com/RRethy/krepe/krepe/pkg/reporef"
	"github.com/RRethy/krepe/krepe/pkg/yaml"
)

var (
	_ yaml.InterfaceUnmarshaler = &Pkg{}
	_ yaml.InterfaceMarshaler   = &Pkg{}
)

type Pkg struct {
	ref  *reporef.RepoRef
	name string
}

type RawPkg struct {
	Ref  string `yaml:"ref,omitempty"`
	Name string `yaml:"name,omitempty"`
}

func NewPkg(ref *reporef.RepoRef, name string) *Pkg {
	return &Pkg{
		ref:  ref,
		name: name,
	}
}

func (p *Pkg) UnmarshalYAML(unmarshal func(interface{}) error) error {
	raw := RawPkg{}
	if err := unmarshal(&raw); err != nil {
		return err
	}

	ref, err := reporef.NewRepoRefFromString(raw.Ref)
	if err != nil {
		return err
	}

	newPkg := NewPkg(ref, raw.Name)
	if err != nil {
		return err
	}

	*p = *newPkg
	return nil
}

func (p *Pkg) MarshalYAML() (interface{}, error) {
	return RawPkg{
		Ref:  p.ref.String(),
		Name: p.name,
	}, nil
}

func (p *Pkg) Name() string {
	if p.name == "" {
		return p.ref.Name
	}

	return p.name
}
