package pkg

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/RRethy/krepe/krepe/pkg/yaml"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Resource struct {
	unstructured.Unstructured

	Filename string
	Raw      []byte
}

func NewResourceFromPath(path string) (*Resource, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return NewResourceFromBytes(filepath.Base(path), raw)
}

func NewResourceFromBytes(fname string, bytes []byte) (*Resource, error) {
	r := &Resource{
		Filename: filepath.Base(fname),
		Raw:      bytes,
	}

	m := make(map[string]any)
	if err := yaml.Unmarshal([]byte(r.Raw), &m); err != nil {
		return nil, fmt.Errorf("unmarshalling resource `%s`: %w", fname, err)
	}

	r.Object = m
	return r, nil
}
