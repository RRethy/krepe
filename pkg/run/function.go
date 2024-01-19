package run

type function struct {
	pkg  string
	name string
}

func newFunction(pkg, name string) *function {
	return &function{pkg: pkg, name: name}
}

func (f *function) run() error {
	return nil
}
