package pipeline

import (
	"fmt"

	"github.com/Shopify/krepe/krepe/pkg/pkg/resource"
	"github.com/Shopify/krepe/krepe/pkg/yaml"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	_ yaml.InterfaceUnmarshaler = &Target{}
	_ yaml.InterfaceMarshaler   = &Target{}
)

type RawTarget struct {
	APIVersion string `yaml:"apiVersion,omitempty"`
	Group      string `yaml:"group,omitempty"`
	Version    string `yaml:"version,omitempty"`
	Kind       string `yaml:"kind,omitempty"`
	Name       string `yaml:"name,omitempty"`
	Namespace  string `yaml:"namespace,omitempty"`
}

type Target struct {
	APIVersion string
	Group      string
	Version    string
	Kind       string
	Name       string
	Namespace  string
}

func (t *Target) UnmarshalYAML(unmarshal func(interface{}) error) error {
	raw := RawTarget{}
	if err := unmarshal(&raw); err != nil {
		return err
	}

	if raw.APIVersion != "" && (raw.Group != "" || raw.Version != "") {
		return fmt.Errorf("cannot specify both apiVersion and group/version")
	}

	if raw.APIVersion != "" {
		gv, err := schema.ParseGroupVersion(raw.APIVersion)
		if err != nil {
			return fmt.Errorf("parsing apiVersion: %w", err)
		}
		raw.Group = gv.Group
		raw.Version = gv.Version
	}

	t.APIVersion = raw.APIVersion
	t.Group = raw.Group
	t.Version = raw.Version
	t.Kind = raw.Kind
	t.Name = raw.Name
	t.Namespace = raw.Namespace
	return nil
}

func (t *Target) MarshalYAML() (interface{}, error) {
	raw := &RawTarget{
		Kind:      t.Kind,
		Name:      t.Name,
		Namespace: t.Namespace,
	}

	if t.APIVersion != "" {
		raw.APIVersion = t.APIVersion
	} else {
		raw.Group = t.Group
		raw.Version = t.Version
	}

	return raw, nil
}

func (t *Target) Matches(res *resource.Resource) bool {
	if t.APIVersion != "" && t.APIVersion != res.GetAPIVersion() {
		return false
	}
	gvk := res.GetObjectKind().GroupVersionKind()
	if t.Group != "" && t.Group != gvk.Group {
		return false
	}
	if t.Version != "" && t.Version != gvk.Version {
		return false
	}
	if t.Kind != "" && t.Kind != res.GetKind() {
		return false
	}
	if t.Name != "" && t.Name != res.GetName() {
		return false
	}
	if t.Namespace != "" && t.Namespace != res.GetNamespace() {
		return false
	}
	return true
}
