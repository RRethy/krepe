package writer

import (
	"github.com/RRethy/krepe/krepe/pkg/pkg"
	"github.com/stretchr/testify/assert"
)

type Mock struct {
	Success bool
	Cnt     int
}

func (m *Mock) Write(pkg *pkg.Package) error {
	if !m.Success {
		return assert.AnError
	}

	m.Cnt++
	return nil
}
