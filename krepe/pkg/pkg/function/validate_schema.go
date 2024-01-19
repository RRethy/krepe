package function

import (
	_ "fmt"

	"github.com/RRethy/krepe/krepe/pkg/pkg/resource"
)

var (
	_ Function = &ValidateSchema{}
)

type ValidateSchema struct{}

func (f *ValidateSchema) WithConfigMap(configMap map[string]any) (Function, error) {
	return f, nil
}

// use https://github.com/goccy/go-yaml#51-print-customized-error-with-yaml-source-code
func (f *ValidateSchema) Run(res *resource.Resource) error {
	return nil
}
