package resource

import (
	"os"

	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Resource struct {
	metav1.TypeMeta   `yaml:",inline"`
	metav1.ObjectMeta `yaml:"metadata,omitempty"`

	Data map[string]any `yaml:",inline,omitempty"`
}

func NewResourceFromPath(path string) (*Resource, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	r := &Resource{}
	err = yaml.Unmarshal(data, r)
	if err != nil {
		return nil, err
	}

	return r, nil
}
