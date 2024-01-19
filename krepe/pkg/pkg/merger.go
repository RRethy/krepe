package pkg

import (
	"github.com/RRethy/krepe/krepe/pkg/merger"
)

var (
	_ merger.Merger[*Pkg] = (*Merger)(nil)
)

type Merger struct {
}

func NewMerger() *Merger {
	return &Merger{}
}

func (m *Merger) Merge(origin, local, upstream *Pkg) (*Pkg, error) {
	// TODO
	return nil, nil
}
