package pkg

type Pipeline []*Step

func (p Pipeline) Run(res *Resource) error {
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
