package functions

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var _ Function = &Succeed{}

type Succeed struct{}

func (s *Succeed) WithConfigMap(_ map[string]any) (Function, error) {
	return s, nil
}

func (s *Succeed) Run(_ *unstructured.Unstructured) error {
	return nil
}
