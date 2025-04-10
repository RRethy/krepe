package functions

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var _ Function = &SetAnnotations{}

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

func (sa *SetAnnotations) Run(res *unstructured.Unstructured) error {
	res.SetAnnotations(sa.annotations)
	return nil
}
