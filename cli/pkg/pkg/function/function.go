package function

import (
	"github.com/Shopify/krepe/cli/pkg/pkg/resource"
)

var Functions = map[string]Function{
	"set-annotations": &SetAnnotations{},
	"set-labels":      &SetLabels{},
	"set-namespace":   &SetNamespace{},
	"set-name":        &SetName{},
	"validate-schema": &ValidateSchema{},
	"tag-image":       &TagImage{},
	"jsonpatch":       &JsonPatch{},
}

// TODO: parsing the configMap should happen BEFORE we execute the function
type Function interface {
	Run(res *resource.Resource, configMap map[string]any) error
}
