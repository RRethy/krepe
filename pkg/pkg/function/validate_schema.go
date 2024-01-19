package function

import (
	_ "fmt"

	"github.com/Shopify/krepe/pkg/pkg/resource"
)

type ValidateSchema struct{}

func (f *ValidateSchema) Run(res *resource.Resource, configMap map[string]any) error {
	return nil
}
