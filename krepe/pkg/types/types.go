package types

import (
	"github.com/RRethy/krepe/krepe/pkg/yaml"
	"github.com/wk8/go-ordered-map/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	_ yaml.BytesUnmarshaler = &Pipelines{}
	_ yaml.BytesMarshaler   = &Pipelines{}
)

type Krepe struct {
	metav1.TypeMeta   `yaml:",inline"`
	metav1.ObjectMeta `yaml:"metadata,omitempty"`

	Imports   Imports    `yaml:"imports,omitempty"`
	Pipelines *Pipelines `yaml:"pipelines,omitempty"`
}

type Imports struct {
	Files    []string        `yaml:"files,omitempty"`
	Packages []PackageImport `yaml:"packages,omitempty"`
}

type PackageImport struct {
	Ref  string `yaml:"ref,omitempty"`
	Name string `yaml:"name,omitempty"`
}

type Pipelines struct {
	orderedmap.OrderedMap[string, []Step] `yaml:",inline"`
}

func (p *Pipelines) UnmarshalYAML(data []byte) error {
	*p = Pipelines{orderedmap.OrderedMap[string, []Step]{}}
	return yaml.UnmarshalCompatibilityShim(data, p)
}

func (p *Pipelines) MarshalYAML() ([]byte, error) {
	return yaml.MarshalCompatibilityShim(&p.OrderedMap)
}

type Step struct {
	Function  string         `yaml:"function,omitempty"`
	Target    Target         `yaml:"target,omitempty"`
	ConfigMap map[string]any `yaml:"configMap,omitempty"`
}

type Target struct {
	APIVersion string `yaml:"apiVersion,omitempty"`
	Group      string `yaml:"group,omitempty"`
	Version    string `yaml:"version,omitempty"`
	Kind       string `yaml:"kind,omitempty"`
	Name       string `yaml:"name,omitempty"`
	Namespace  string `yaml:"namespace,omitempty"`
}
