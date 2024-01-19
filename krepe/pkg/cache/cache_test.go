package cache

import (
	"path/filepath"
	"testing"

	"github.com/adrg/xdg"
	"github.com/stretchr/testify/assert"
)

func TestPath(t *testing.T) {
	path := t.TempDir()
	xdg.CacheHome = path

	cache := NewClient()
	cachePath := cache.Path()
	assert.NotEmpty(t, cachePath)
	assert.DirExists(t, cachePath)
	assert.Equal(t, filepath.Join(path, dirName), cachePath)
}
