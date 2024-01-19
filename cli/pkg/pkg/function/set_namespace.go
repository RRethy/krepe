package function

import (
	"fmt"

	"github.com/Shopify/krepe/cli/pkg/pkg/resource"
)

type SetNamespace struct{}

func (f *SetNamespace) Run(res *resource.Resource, configMap map[string]any) error {
	ns, ok := configMap["namespace"].(string)
	if !ok {
		return fmt.Errorf("failed to get a key `namespace` with type `string` form configMap")
	}
	res.SetNamespace(ns)
	return nil
}
