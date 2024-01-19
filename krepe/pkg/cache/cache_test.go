package cache

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPath(t *testing.T) {
	path := t.TempDir()
	cache := NewCache(WithDir(path))
	cachePath, err := cache.Path()
	assert.NoError(t, err)
	assert.NotEmpty(t, cachePath)
	assert.DirExists(t, cachePath)
	assert.Equal(t, filepath.Join(path, "krepe"), cachePath)
}
