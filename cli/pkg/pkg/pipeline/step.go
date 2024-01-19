package pipeline

import (
	"fmt"

	"github.com/Shopify/krepe/cli/pkg/pkg/function"
	"github.com/Shopify/krepe/cli/pkg/pkg/resource"
)

type Step struct {
	name string
	fn   function.Function
}

type StepRaw struct {
	Function  string         `yaml:"function,omitempty"`
	ConfigMap map[string]any `yaml:"configMap,omitempty"`
}

func (m *Step) UnmarshalYAML(unmarshal func(interface{}) error) error {
	raw := StepRaw{}
	if err := unmarshal(&raw); err != nil {
		return err
	}

	f, err := function.NewFunction(raw.Function, raw.ConfigMap)
	if err != nil {
		return err
	}

	m.fn = f
	m.name = raw.Function
	return nil
}

func (s *Step) Run(res *resource.Resource) error {
	err := s.fn.Run(res)
	if err != nil {
		return fmt.Errorf(
			"running function `%s` on resource `%s`: %w",
			s.name,
			res.Fname(),
			err,
		)
	}

	return nil
}
