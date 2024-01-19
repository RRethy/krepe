package pipeline

import (
	"fmt"

	"github.com/Shopify/krepe/krepe/pkg/pkg/function"
	"github.com/Shopify/krepe/krepe/pkg/pkg/resource"
	"github.com/Shopify/krepe/krepe/pkg/yaml"
)

var (
	_ yaml.InterfaceUnmarshaler = &Step{}
	_ yaml.InterfaceMarshaler   = &Step{}
)

type Step struct {
	name      string
	fn        function.Function
	target    *Target
	configMap map[string]any
}

type RawStep struct {
	Function  string         `yaml:"function,omitempty"`
	Target    *Target        `yaml:"target,omitempty"`
	ConfigMap map[string]any `yaml:"configMap,omitempty"`
}

func (m *Step) UnmarshalYAML(unmarshal func(interface{}) error) error {
	raw := RawStep{}
	if err := unmarshal(&raw); err != nil {
		return err
	}

	f, err := function.NewFunction(raw.Function, raw.ConfigMap)
	if err != nil {
		return err
	}

	m.name = raw.Function
	m.fn = f
	m.target = raw.Target
	m.configMap = raw.ConfigMap
	return nil
}

func (m *Step) MarshalYAML() (interface{}, error) {
	return RawStep{
		Function:  m.name,
		Target:    m.target,
		ConfigMap: m.configMap,
	}, nil
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

func (s *Step) Name() string {
	return s.name
}
