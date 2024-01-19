package function

import (
	"fmt"

	"github.com/RRethy/krepe/krepe/pkg/pkg/resource"
)

type SetName struct {
	name string
}

func (sn *SetName) WithConfigMap(configMap map[string]any) (Function, error) {
	name, ok := configMap["name"].(string)
	if !ok {
		return nil, fmt.Errorf("getting a key `name` with type `string` form configMap: %s", name)
	}

	return &SetName{
		name: name,
	}, nil
}

func (sn *SetName) Run(res *resource.Resource) error {
	res.SetName(sn.name)
	return nil
}
