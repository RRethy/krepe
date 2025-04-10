package functions

import (
	"errors"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var _ Function = &ConfigFail{}

type ConfigFail struct{}

func (cf *ConfigFail) WithConfigMap(configMap map[string]any) (Function, error) {
	return nil, errors.New("config fail")
}

func (cf *ConfigFail) Run(res *unstructured.Unstructured) error {
	return nil
}
