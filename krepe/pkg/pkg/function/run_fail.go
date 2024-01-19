package function

import (
	"errors"

	"github.com/RRethy/krepe/krepe/pkg/pkg/resource"
)

type RunFail struct{}

func (rf *RunFail) WithConfigMap(configMap map[string]any) (Function, error) {
	return rf, nil
}

func (rf *RunFail) Run(res *resource.Resource) error {
	return errors.New("run fail")
}
