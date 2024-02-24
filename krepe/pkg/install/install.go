package install

import (
	"os"

	"github.com/RRethy/krepe/krepe/pkg/writer"
)

func Install(pkgPath, name string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	return (&Installer{Writer: &writer.Disk{}}).Install(dir, pkgPath, name)
}
