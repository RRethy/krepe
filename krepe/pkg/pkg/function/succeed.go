package function

import (
	"github.com/Shopify/krepe/krepe/pkg/pkg/resource"
)

type Succeed struct{}

func (s *Succeed) WithConfigMap(_ map[string]any) (Function, error) {
	return s, nil
}

func (s *Succeed) Run(_ *resource.Resource) error {
	return nil
}
