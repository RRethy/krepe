package pipeline

import (
	"github.com/RRethy/krepe/krepe/pkg/pkg/resource"
)

type Pipeline []*Step

func (p Pipeline) Run(res *resource.Resource) error {
	for _, step := range p {
		if step.Matches(res) {
			err := step.Run(res)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
