package function

import (
	_ "fmt"

	"github.com/Shopify/krepe/pkg/pkg/resource"
)

type SetField struct{}

func (f *SetField) Run(res *resource.Resource, configMap map[string]any) error {
	return nil
}
