package function

import (
	_ "fmt"

	"github.com/RRethy/krepe/krepe/pkg/pkg/resource"
)

type TagImage struct{}

func (f *TagImage) WithConfigMap(configMap map[string]any) (Function, error) {
	return f, nil
}

func (f *TagImage) Run(res *resource.Resource) error {
	return nil
}
