package pkg

import (
	"github.com/RRethy/krepe/krepe/pkg/types"
)

type Pipeline struct {
	Name  string
	Steps []*Step
}

func NewPipeline(typesPipeline types.Pipeline) (Pipeline, error) {
	var err error
	pipeline := Pipeline{
		Name:  typesPipeline.Name,
		Steps: make([]*Step, len(typesPipeline.Steps)),
	}

	for i, step := range typesPipeline.Steps {
		pipeline.Steps[i], err = NewStep(step)
		if err != nil {
			return Pipeline{}, err
		}
	}

	return pipeline, nil
}

func (p Pipeline) Run(res *Resource) error {
	for _, step := range p.Steps {
		err := step.Run(res)
		if err != nil {
			return err
		}
	}
	return nil
}
