package pkg

import (
	"os"

	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type resource struct {
	metav1.TypeMeta   `yaml:",inline"`
	metav1.ObjectMeta `yaml:"metadata,omitempty"`

	data map[string]any `yaml:",inline,omitempty"`
}

func newResourceFromPath(path string) (*resource, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	r := &resource{}
	err = yaml.Unmarshal(data, r)
	if err != nil {
		return nil, err
	}

	return r, nil
}
