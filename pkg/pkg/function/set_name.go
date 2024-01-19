package function

import (
	_ "fmt"

	"github.com/Shopify/krepe/pkg/pkg/resource"
)

type SetName struct{}

func (f *SetName) Run(res *resource.Resource, configMap map[string]any) error {
	return nil
}
