package writer

import (
	"github.com/RRethy/krepe/krepe/pkg/pkg"
	"github.com/stretchr/testify/assert"
)

var _ Writer = &Mock{}

type Mock struct {
	Success bool
	Cnt     int
	Pkg     *pkg.Package
}

func (m *Mock) Write(pkg *pkg.Package, _ string) error {
	if !m.Success {
		return assert.AnError
	}

	m.Pkg = pkg

	m.Cnt++
	return nil
}
