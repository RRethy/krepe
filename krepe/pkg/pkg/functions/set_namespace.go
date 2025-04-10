package functions

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var _ Function = &SetNamespace{}

type SetNamespace struct {
	namespace string
}

func (f *SetNamespace) WithConfigMap(configMap map[string]any) (Function, error) {
	ns, ok := configMap["namespace"].(string)
	if !ok {
		return nil, fmt.Errorf("getting a key `namespace` with type `string` form configMap: %s", ns)
	}

	return &SetNamespace{
		namespace: ns,
	}, nil
}

func (sn *SetNamespace) Run(res *unstructured.Unstructured) error {
	res.SetNamespace(sn.namespace)
	return nil
}
