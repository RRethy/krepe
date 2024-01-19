package function

import (
	"github.com/RRethy/krepe/krepe/pkg/pkg/resource"
)

var (
	_ Function = &Succeed{}
)

type Succeed struct{}

func (s *Succeed) WithConfigMap(_ map[string]any) (Function, error) {
	return s, nil
}

func (s *Succeed) Run(_ *resource.Resource) error {
	return nil
}
