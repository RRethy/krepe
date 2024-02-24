package update

import (
	"os"

	"github.com/RRethy/krepe/krepe/pkg/merger"
	"github.com/RRethy/krepe/krepe/pkg/writer"
)

func Update(packageName string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	return (&Updater{Merger: merger.Package{}, Writer: &writer.Disk{}}).Update(dir, packageName)
}
