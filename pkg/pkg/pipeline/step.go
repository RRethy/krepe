package pipeline

import (
	"fmt"

	"github.com/Shopify/krepe/pkg/pkg/function"
	"github.com/Shopify/krepe/pkg/pkg/resource"
)

type Step struct {
	Function  string         `yaml:"function,omitempty"`
	ConfigMap map[string]any `yaml:"configMap,omitempty"`
}

func (s *Step) Run(res *resource.Resource) error {
	f, ok := function.Functions[s.Function]
	if !ok {
		return fmt.Errorf("function `%s` not found", s.Function)
	}

	err := f.Run(res, s.ConfigMap)
	if err != nil {
		return fmt.Errorf("failed to run function `%s`: %w", s.Function, err)
	}
	return nil
}
