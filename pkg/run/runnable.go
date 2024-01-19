package run

import (
	"errors"
)

type runnable interface {
	run() error
}

func newRunnable(pkg, pipeline, function string) (runnable, error) {
	if pipeline != "." && function != "" {
		return nil, errors.New("cannot specify both pipeline and function")
	}
	if function != "" {
		return newFunction(pkg, function), nil
	}
	return newPipeline(pkg, pipeline), nil
}
