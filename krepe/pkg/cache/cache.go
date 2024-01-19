package cache

import (
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

const (
	dirName = "krepe"
)

func Path() string {
	path := filepath.Join(xdg.CacheHome, dirName)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
	}

	return path
}
