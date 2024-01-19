package run

type pipeline struct {
	pkg  string
	name string
}

func newPipeline(pkg, name string) *pipeline {
	return &pipeline{pkg: pkg, name: name}
}

func (p *pipeline) run() error {
	return nil
}
