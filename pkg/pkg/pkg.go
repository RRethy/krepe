package pkg

import (
	"fmt"
	"os"
	"path/filepath"
)

type Pkg struct {
	name  string
	krepe *Krepe
}

func NewPkgFromPath(pkgPath string) (*Pkg, error) {
	if pkgPath == "/" {
		return nil, fmt.Errorf("pkg path cannot be `/`")
	}

	fileInfo, err := os.Stat(pkgPath)
	if err != nil {
		return nil, err
	}
	if !fileInfo.IsDir() {
		return nil, fmt.Errorf("pkg path is not a directory: %s", pkgPath)
	}

	name := filepath.Base(pkgPath)
	krepe, err := NewKrepeFromPath(filepath.Join(pkgPath, "krepe.yaml"))
	if err != nil {
		return nil, err
	}

	return &Pkg{
		name:  name,
		krepe: krepe,
	}, nil
}

func (p *Pkg) RunPipeline(name string) error {
	return nil
}
