package function

import (
	"fmt"

	"github.com/RRethy/krepe/krepe/pkg/pkg/resource"
)

var (
	_ Function = &AddAnnotations{}
)

type AddAnnotations struct {
	annotations map[string]string
}

func (aa *AddAnnotations) WithConfigMap(configMap map[string]any) (Function, error) {
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

	return &AddAnnotations{
		annotations: annotations,
	}, nil
}

func (aa *AddAnnotations) Run(res *resource.Resource) error {
	annotations := res.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}

	for k, v := range aa.annotations {
		annotations[k] = v
	}

	res.SetAnnotations(annotations)
	return nil
}
