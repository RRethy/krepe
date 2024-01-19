package function

import (
	"fmt"

	"github.com/RRethy/krepe/krepe/pkg/pkg/resource"
)

var (
	_ Function = &SetNamespace{}
)

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

func (sn *SetNamespace) Run(res *resource.Resource) error {
	res.SetNamespace(sn.namespace)
	return nil
}
