package resource

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Resource struct {
	unstructured.Unstructured

	fname string
	raw   []byte
}

func NewResourceFromPath(path string) (*Resource, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	r := &Resource{
		fname: filepath.Base(path),
		raw:   raw,
	}

	m := make(map[string]any)
	if err := yaml.Unmarshal([]byte(r.raw), &m); err != nil {
		panic(err)
	}
	r.Object = m
	return r, nil
}

func (r *Resource) Fname() string {
	return r.fname
}
