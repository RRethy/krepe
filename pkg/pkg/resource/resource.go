package resource

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Resource struct {
	name string
	raw  []byte

	metav1.TypeMeta   `yaml:",inline"`
	metav1.ObjectMeta `yaml:"metadata,omitempty"`

	Data map[string]any `yaml:",inline,omitempty"`
}

func NewResourceFromPath(path string) (*Resource, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	r := &Resource{
		name: filepath.Base(path),
		raw:  raw,
	}
	err = yaml.Unmarshal(raw, r)
	if err != nil {
		return nil, err
	}

	return r, nil
}
