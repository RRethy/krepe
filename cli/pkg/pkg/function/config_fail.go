package function

import (
	"errors"

	"github.com/Shopify/krepe/cli/pkg/pkg/resource"
)

type ConfigFail struct{}

func (cf *ConfigFail) WithConfigMap(configMap map[string]any) (Function, error) {
	return nil, errors.New("config fail")
}

func (cf *ConfigFail) Run(res *resource.Resource) error {
	return nil
}
