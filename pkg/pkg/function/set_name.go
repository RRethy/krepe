package function

import (
	"fmt"

	"github.com/Shopify/krepe/pkg/pkg/resource"
)

type SetName struct{}

func (f *SetName) Run(res *resource.Resource, configMap map[string]any) error {
	if v, ok := configMap["name"]; ok {
		if ns, ok := v.(string); ok {
			res.SetName(ns)
		} else {
			return fmt.Errorf("invalid type for name key in configMap: %T", v)
		}
	} else {
		return fmt.Errorf("missing name key in configMap")
	}
	return nil
}
