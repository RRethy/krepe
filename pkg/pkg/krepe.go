package pkg

import (
	"os"

	"github.com/Shopify/krepe/pkg/pkg/imports"
	"github.com/Shopify/krepe/pkg/pkg/pipeline"
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Krepe struct {
	metav1.TypeMeta   `yaml:",inline"`
	metav1.ObjectMeta `yaml:"metadata,omitempty"`

	Imports   *imports.Imports     `yaml:"imports,omitempty"`
	Pipelines []*pipeline.Pipeline `yaml:"pipelines,omitempty"`
}

func NewKrepeFromPath(krepePath string) (*Krepe, error) {
	data, err := os.ReadFile(krepePath)
	if err != nil {
		return nil, err
	}

	k := &Krepe{}
	err = yaml.Unmarshal(data, k)
	if err != nil {
		return nil, err
	}

	return k, nil
}
