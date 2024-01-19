package pipeline

import (
	"github.com/Shopify/krepe/cli/pkg/pkg/resource"
)

type Pipeline []*Step

func (p Pipeline) Run(res *resource.Resource) error {
	for _, step := range p {
		err := step.Run(res)
		if err != nil {
			return err
		}
	}
	return nil
}
