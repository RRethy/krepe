package run

import (
	"os"

	"github.com/RRethy/krepe/krepe/pkg/writer"
)

func Run(pipeline string) error {
	pkgPath, err := os.Getwd()
	if err != nil {
		return err
	}

	return (&Pipeline{Writer: &writer.Disk{}}).Run(pkgPath, pipeline)
}
