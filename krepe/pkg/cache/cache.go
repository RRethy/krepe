package cache

import (
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

const (
	dirName = "krepe"
)

type Cache interface {
	Path() string
}

type xdgCache struct{}

func NewCache() Cache {
	return &xdgCache{}
}

func (c *xdgCache) Path() string {
	path := filepath.Join(xdg.CacheHome, dirName)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
	}

	return path
}
