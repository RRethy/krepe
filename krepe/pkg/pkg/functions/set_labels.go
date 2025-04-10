package functions

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var _ Function = &SetLabels{}

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

func (sl *SetLabels) Run(res *unstructured.Unstructured) error {
	res.SetLabels(sl.labels)
	return nil
}
