package run

import (
	"errors"

	"github.com/Shopify/krepe/cli/pkg/pkg"
)

type runnable interface {
	run() error
}

func newRunnable(pkgPath, pipeline, function string) (runnable, error) {
	pkg, err := pkg.NewPkgFromPath(pkgPath)
	if err != nil {
		return nil, err
	}

	if pipeline != "." && function != "" {
		return nil, errors.New("cannot specify both pipeline and function")
	}
	if function != "" {
		return newFunction(pkg, function), nil
	}
	return newPipeline(pkg, pipeline), nil
}
