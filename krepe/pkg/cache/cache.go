package cache

import (
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

const (
	dirName = "krepe"
)

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Path() string {
	path := filepath.Join(xdg.CacheHome, dirName)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
	}

	return path
}
