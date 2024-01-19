package function

import (
	"github.com/Shopify/krepe/pkg/pkg/resource"
)

type SetAnnotations struct {
}

func (f *SetAnnotations) Run(res *resource.Resource, configMap map[string]any) error {
	return nil
}
