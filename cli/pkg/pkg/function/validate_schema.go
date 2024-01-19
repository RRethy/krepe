package function

import (
	_ "fmt"

	"github.com/Shopify/krepe/cli/pkg/pkg/resource"
)

type ValidateSchema struct{}

func (f *ValidateSchema) WithConfigMap(configMap map[string]any) (Function, error) {
	return f, nil
}

func (f *ValidateSchema) Run(res *resource.Resource) error {
	return nil
}
