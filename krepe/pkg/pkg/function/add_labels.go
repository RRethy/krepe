package function

import (
	"fmt"

	"github.com/RRethy/krepe/krepe/pkg/pkg/resource"
)

type AddLabels struct {
	labels map[string]string
}

func (al *AddLabels) WithConfigMap(configMap map[string]any) (Function, error) {
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

	return &AddLabels{
		labels: labels,
	}, nil
}

func (al *AddLabels) Run(res *resource.Resource) error {
	labels := res.GetLabels()
	if labels == nil {
		labels = make(map[string]string)
	}

	for k, v := range al.labels {
		labels[k] = v
	}

	res.SetLabels(labels)
	return nil
}
