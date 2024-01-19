package functions

import (
	"errors"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var (
	_ Function = &RunFail{}
)

type RunFail struct{}

func (rf *RunFail) WithConfigMap(configMap map[string]any) (Function, error) {
	return rf, nil
}

func (rf *RunFail) Run(res *unstructured.Unstructured) error {
	return errors.New("run fail")
}
