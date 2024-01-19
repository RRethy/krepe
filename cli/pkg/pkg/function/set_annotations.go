package function

import (
	"fmt"

	"github.com/Shopify/krepe/cli/pkg/pkg/resource"
)

type SetAnnotations struct{}

func (f *SetAnnotations) Run(res *resource.Resource, configMap map[string]any) error {
	annotations := res.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}

	for k, v := range configMap {
		if s, ok := v.(string); ok {
			annotations[k] = s
		} else {
			return fmt.Errorf("invalid type for annotation `%s`: %T", k, v)
		}
	}

	res.SetAnnotations(annotations)
	return nil

}
