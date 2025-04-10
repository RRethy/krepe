package functions

import (
	_ "fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var _ Function = &TagImage{}

type TagImage struct{}

func (f *TagImage) WithConfigMap(configMap map[string]any) (Function, error) {
	return f, nil
}

func (f *TagImage) Run(res *unstructured.Unstructured) error {
	return nil
}
