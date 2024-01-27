package pkg

import (
	"github.com/RRethy/krepe/krepe/pkg/types"
)

type Pipeline struct {
	Steps []*Step
}

func NewPipeline(steps []types.Step) (Pipeline, error) {
	var err error
	pipeline := Pipeline{
		Steps: make([]*Step, len(steps)),
	}

	for i, step := range steps {
		pipeline.Steps[i], err = NewStep(step)
		if err != nil {
			return Pipeline{}, err
		}
	}

	return pipeline, nil
}

func (p Pipeline) Run(res *Resource) error {
	for _, step := range p.Steps {
		if step.Matches(res) {
			err := step.Run(res)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
