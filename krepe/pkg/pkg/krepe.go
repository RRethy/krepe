package pkg

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/RRethy/krepe/krepe/pkg/pkg/imports"
	"github.com/RRethy/krepe/krepe/pkg/pkg/pipeline"
	"github.com/RRethy/krepe/krepe/pkg/yaml"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Krepe struct {
	metav1.TypeMeta   `yaml:",inline"`
	metav1.ObjectMeta `yaml:"metadata,omitempty"`

	Imports   *imports.Imports    `yaml:"imports,omitempty"`
	Pipelines *pipeline.Pipelines `yaml:"pipelines,omitempty"`
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

func (k *Krepe) AddPkgImport(p *imports.Pkg) {
	k.Imports.AddPackage(p)
}

func (k *Krepe) Write(dir string) error {
	path := filepath.Join(dir, "krepe.yaml")

	data, err := yaml.Marshal(k)
	if err != nil {
		return fmt.Errorf("marshalling krepe.yaml: %w", err)
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return fmt.Errorf("writing krepe.yaml to file: %w", err)
	}

	return nil
}
