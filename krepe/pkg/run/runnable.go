package run

import (
	"errors"

	"github.com/Shopify/krepe/krepe/pkg/pkg"
)

type runnable interface {
	run(dir string) error
}

func newRunnable(pkgPath, pipeline, function string) (runnable, error) {
	pkg, err := pkg.NewPkgFromPath(pkgPath)
	if err != nil {
		return nil, err
	}

	if pipeline == "" && function == "" {
		return nil, errors.New("must specify either pipeline or function")
	}

	if function == "" {
		return newPipeline(pkg, pipeline), nil
	}

	if pipeline == "" || pipeline == "." {
		return newFunction(pkg, function), nil
	}

	return nil, errors.New("cannot specify both pipeline and function")
}
