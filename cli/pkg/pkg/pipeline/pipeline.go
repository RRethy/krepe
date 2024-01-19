package pipeline

import (
	"github.com/Shopify/krepe/cli/pkg/pkg/resource"
	"golang.design/x/reflect"
)

type Pipeline struct {
	Steps []*Step `yaml:"steps,omitempty"`
}

func (p *Pipeline) Run(res *resource.Resource) (*resource.Resource, error) {
	copied := reflect.DeepCopy(res)
	for _, step := range p.Steps {
		err := step.Run(copied)
		if err != nil {
			return nil, err
		}
	}
	return copied, nil
}
