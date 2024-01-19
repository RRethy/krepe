package function

import (
	"fmt"

	"github.com/Shopify/krepe/cli/pkg/pkg/resource"
)

type SetAnnotations struct {
	annotations map[string]string
}

func (sa *SetAnnotations) WithConfigMap(configMap map[string]any) (Function, error) {
	annotations := make(map[string]string)
	for k, v := range configMap {
		if s, ok := v.(string); ok {
			annotations[k] = s
		} else {
			return nil, fmt.Errorf("invalid type for annotation `%s`: %T", k, v)
		}
	}

	if len(annotations) == 0 {
		return nil, fmt.Errorf("no annotations provided")
	}

	return &SetAnnotations{
		annotations: annotations,
	}, nil
}

func (sa *SetAnnotations) Run(res *resource.Resource) error {
	res.SetAnnotations(sa.annotations)
	return nil
}
