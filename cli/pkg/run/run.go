package run

import (
	"fmt"
)

func Run(pkg, pipeline, function string) error {
	r, err := newRunnable(pkg, pipeline, function)
	if err != nil {
		return fmt.Errorf("failed to create runnable: %w", err)
	}

	err = r.run()
	if err != nil {
		return fmt.Errorf("failed to run in pkg `%s`: %w", pkg, err)
	}

	return nil
}
