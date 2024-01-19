package functions

import (
	_ "fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var (
	_ Function = &ValidateSchema{}
)

type ValidateSchema struct{}

func (f *ValidateSchema) WithConfigMap(configMap map[string]any) (Function, error) {
	return f, nil
}

// use https://github.com/goccy/go-yaml#51-print-customized-error-with-yaml-source-code
func (f *ValidateSchema) Run(res *unstructured.Unstructured) error {
	return nil
}
