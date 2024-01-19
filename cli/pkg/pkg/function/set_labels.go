package function

import (
	"fmt"

	"github.com/Shopify/krepe/cli/pkg/pkg/resource"
)

type SetLabels struct {
	labels map[string]string
}

func (sl *SetLabels) WithConfigMap(configMap map[string]any) (Function, error) {
	labels := make(map[string]string)
	for k, v := range configMap {
		if s, ok := v.(string); ok {
			labels[k] = s
		} else {
			return nil, fmt.Errorf("invalid type for label `%s`: %T", k, v)
		}
	}

	if len(labels) == 0 {
		return nil, fmt.Errorf("no labels provided")
	}

	return &SetLabels{
		labels: labels,
	}, nil
}

func (sl *SetLabels) Run(res *resource.Resource) error {
	res.SetLabels(sl.labels)
	return nil
}
