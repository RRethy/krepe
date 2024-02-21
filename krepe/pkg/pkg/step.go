package pkg

import (
	"fmt"

	"github.com/RRethy/krepe/krepe/pkg/pkg/functions"
	"github.com/RRethy/krepe/krepe/pkg/types"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Step struct {
	Name   string
	Fn     functions.Function
	Target Target
	Config map[string]any
}

func NewStep(step types.Step) (*Step, error) {
	fn, err := functions.NewFunction(step.Function, step.Config)
	if err != nil {
		return nil, err
	}

	target, err := NewTarget(step.Target)
	if err != nil {
		return nil, err
	}

	return &Step{
		Name:   step.Function,
		Fn:     fn,
		Target: target,
		Config: step.Config,
	}, nil
}

func (s *Step) Run(resource *unstructured.Unstructured) error {
	if !s.Target.Matches(resource) {
		return nil
	}

	err := s.Fn.Run(resource)
	if err != nil {
		return fmt.Errorf(
			"running functions `%s` on resource `%s`: %w",
			s.Name,
			resource.GetName(),
			err,
		)
	}

	return nil
}
