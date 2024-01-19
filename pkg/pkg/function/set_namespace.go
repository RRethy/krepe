package function

import (
	"fmt"

	"github.com/Shopify/krepe/pkg/pkg/resource"
)

type SetNamespace struct{}

func (f *SetNamespace) Run(res *resource.Resource, configMap map[string]any) error {
	if v, ok := configMap["namespace"]; ok {
		if ns, ok := v.(string); ok {
			res.SetNamespace(ns)
		} else {
			return fmt.Errorf("invalid type for namespace key in configMap: %T", v)
		}
	} else {
		return fmt.Errorf("missing namespace key in configMap")
	}
	return nil
}
