package pkg

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/RRethy/krepe/krepe/pkg/yaml"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type FileImport struct {
	Resource *unstructured.Unstructured
	Name     string
	Raw      []byte
}

func NewFileImportFromPath(path string) (FileImport, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return FileImport{}, err
	}

	return NewFileImportFromBytes(filepath.Base(path), raw)
}

func NewFileImportFromBytes(filename string, bytes []byte) (FileImport, error) {
	r := &FileImport{
		Resource: &unstructured.Unstructured{
			Object: make(map[string]any),
		},
		Name: filename,
		Raw:  bytes,
	}

	if err := yaml.Unmarshal(bytes, &r.Resource.Object); err != nil {
		return FileImport{}, fmt.Errorf("unmarshalling resource `%s`: %w", filename, err)
	}

	return *r, nil
}
