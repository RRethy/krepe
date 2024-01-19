package function

import (
	"github.com/Shopify/krepe/pkg/pkg/resource"
)

var Functions = map[string]Function{
	"set-annotations": &SetAnnotations{},
	"set-labels":      &SetLabels{},
	"set-namespace":   &SetNamespace{},
	"set-name":        &SetName{},
	"set-field":       &SetField{},
	"validate-schema": &ValidateSchema{},
	"tag-image":       &TagImage{},
}

type Function interface {
	Run(res *resource.Resource, configMap map[string]any) error
}
