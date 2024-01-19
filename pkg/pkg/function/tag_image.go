package function

import (
	_ "fmt"

	"github.com/Shopify/krepe/pkg/pkg/resource"
)

type TagImage struct{}

func (f *TagImage) Run(res *resource.Resource, configMap map[string]any) error {
	return nil
}
