package cache

import (
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

const (
	name = "krepe"
)

type Cache struct {
	dir string
}

type Option func(*Cache)

func WithDir(dir string) Option {
	return func(c *Cache) {
		c.dir = dir
	}
}

func NewCache(options ...Option) *Cache {
	cache := &Cache{
		dir: xdg.CacheHome,
	}
	for _, option := range options {
		option(cache)
	}
	return cache
}

func (c *Cache) Path() (string, error) {
	path := filepath.Join(c.dir, name)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, 0755)
		if err != nil {
			return "", err
		}
	}

	return path, nil
}
