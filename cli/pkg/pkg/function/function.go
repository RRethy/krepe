package function

import (
	"fmt"

	"github.com/Shopify/krepe/cli/pkg/pkg/resource"
)

var functions = map[string]Function{
	"set-annotations": &SetAnnotations{},
	"set-labels":      &SetLabels{},
	"set-namespace":   &SetNamespace{},
	"set-name":        &SetName{},
	"validate-schema": &ValidateSchema{},
	"tag-image":       &TagImage{},
	"jsonpatch":       &JsonPatch{},
}

type Function interface {
	WithConfigMap(configMap map[string]any) (Function, error)
	Run(res *resource.Resource) error
}

func NewFunction(name string, configMap map[string]any) (Function, error) {
	f, ok := functions[name]
	if !ok {
		return nil, fmt.Errorf("function `%s` not found", name)
	}

	f, err := f.WithConfigMap(configMap)
	if err != nil {
		return nil, fmt.Errorf("initializing function `%s`: %w", name, err)
	}

	return f, nil
}
