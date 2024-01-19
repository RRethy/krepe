package imports

type Pkg struct {
	url  string
	name string
}

type RawPkg struct {
	Url  string `yaml:"url,omitempty"`
	Name string `yaml:"name,omitempty"`
}

func NewPkg(url string, name string) (*Pkg, error) {
	return &Pkg{
		url:  url,
		name: name,
	}, nil // TODO
}

func (p *Pkg) UnmarshalYAML(unmarshal func(interface{}) error) error {
	raw := RawPkg{}
	if err := unmarshal(&raw); err != nil {
		return err
	}

	newPkg, err := NewPkg(raw.Url, raw.Name)
	if err != nil {
		return err
	}

	*p = *newPkg
	return nil
}

func (p *Pkg) MarshalYAML() (interface{}, error) {
	return RawPkg{
		Url:  p.url,
		Name: p.name,
	}, nil
}

func (p *Pkg) Name() string {
	return p.name
}
