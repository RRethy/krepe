package imports

type Pkg struct {
	Url *string
}

func (p *Pkg) Name() string {
	return *p.Url
}
