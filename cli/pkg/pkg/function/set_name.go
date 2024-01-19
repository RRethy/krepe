package function

import (
	"fmt"

	"github.com/Shopify/krepe/cli/pkg/pkg/resource"
)

type SetName struct{}

func (f *SetName) Run(res *resource.Resource, configMap map[string]any) error {
	name, ok := configMap["name"].(string)
	if !ok {
		return fmt.Errorf("failed to get a key `name` with type `string` form configMap")
	}
	res.SetName(name)
	return nil
}
