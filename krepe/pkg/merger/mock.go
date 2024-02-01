package merger

import (
	"testing"

	"github.com/RRethy/krepe/krepe/pkg/pkg"
	"github.com/stretchr/testify/assert"
)

var (
	_ Merger = &Mock{}
)

type Mock struct {
	Origin   *pkg.Package
	Local    *pkg.Package
	Upstream *pkg.Package
	Cnt      int
	Success  bool
}

func (m *Mock) Merge(origin, local, upstream *pkg.Package) (*pkg.Package, error) {
	if !m.Success {
		return nil, assert.AnError
	}

	m.Origin = origin
	m.Local = local
	m.Upstream = upstream
	m.Cnt++
	return local, nil
}

func (m *Mock) Assert(t *testing.T, origin, local, upstream string, cnt int) {
	assert.Equal(t, origin, m.Origin.Labels["version"])
	assert.Equal(t, local, m.Local.Labels["version"])
	assert.Equal(t, upstream, m.Upstream.Labels["version"])
	assert.Equal(t, cnt, m.Cnt)
}
