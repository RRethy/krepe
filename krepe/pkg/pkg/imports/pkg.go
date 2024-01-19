package imports

import (
	"errors"
	"strings"

	"github.com/RRethy/krepe/krepe/pkg/yaml"
)

var (
	_ yaml.InterfaceUnmarshaler = &Pkg{}
	_ yaml.InterfaceMarshaler   = &Pkg{}
)

type Pkg struct {
	url     string
	version string
	name    string
}

type RawPkg struct {
	Tag  string `yaml:"tag,omitempty"`
	Name string `yaml:"name,omitempty"`
}

func NewPkg(tag string, name string) (*Pkg, error) {
	parts := strings.Split(tag, "@")
	if len(parts) != 2 {
		return nil, errors.New("tag must be in the format <url>@<version>")
	}

	url, version := parts[0], parts[1]

	if name == "" {
		urlParts := strings.Split(url, "/")
		name = urlParts[len(urlParts)-1]
		if name == "" {
			return nil, errors.New("could not parse name from url")
		}
	}

	return &Pkg{
		url:     url,
		version: version,
		name:    name,
	}, nil
}

func (p *Pkg) UnmarshalYAML(unmarshal func(interface{}) error) error {
	raw := RawPkg{}
	if err := unmarshal(&raw); err != nil {
		return err
	}

	newPkg, err := NewPkg(raw.Tag, raw.Name)
	if err != nil {
		return err
	}

	*p = *newPkg
	return nil
}

func (p *Pkg) MarshalYAML() (interface{}, error) {
	return RawPkg{
		Tag:  p.url + "@" + p.version,
		Name: p.name,
	}, nil
}

func (p *Pkg) URL() string {
	return p.url
}

func (p *Pkg) Version() string {
	return p.version
}

func (p *Pkg) Name() string {
	return p.name
}
