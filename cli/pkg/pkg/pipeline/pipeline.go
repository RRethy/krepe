package pipeline

import (
	"github.com/Shopify/krepe/cli/pkg/pkg/resource"
)

type Pipeline struct {
	Steps []*Step `yaml:"steps,omitempty"`
}

func (p *Pipeline) Run(res *resource.Resource) error {
	for _, step := range p.Steps {
		err := step.Run(res)
		if err != nil {
			return err
		}
	}
	return nil
}
