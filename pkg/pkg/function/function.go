package function

import (
	"github.com/Shopify/krepe/pkg/pkg/resource"
)

type Function interface {
	Run(res *resource.Resource, configMap map[string]any) error
}
