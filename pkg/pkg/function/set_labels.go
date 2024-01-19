package function

import (
	"fmt"

	"github.com/Shopify/krepe/pkg/pkg/resource"
)

type SetLabels struct{}

func (f *SetLabels) Run(res *resource.Resource, configMap map[string]any) error {
	labels := res.GetLabels()
	if labels == nil {
		labels = make(map[string]string)
	}

	for k, v := range configMap {
		if s, ok := v.(string); ok {
			labels[k] = s
		} else {
			return fmt.Errorf("invalid type for label `%s`: %T", k, v)
		}
	}

	res.SetLabels(labels)
	return nil
}
