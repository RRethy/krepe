package function

import (
	"fmt"

	"github.com/Shopify/krepe/cli/pkg/pkg/resource"
	"github.com/Shopify/krepe/jsonpatch"
)

type JsonPatch struct{}

func (f *JsonPatch) Run(res *resource.Resource, configMap map[string]any) error {
	var op, from, path string
	var value any

	if _, ok := configMap["op"]; ok {
		op, ok = configMap["op"].(string)
		if !ok {
			return fmt.Errorf("invalid op type")
		}
	}

	if _, ok := configMap["from"]; ok {
		from, ok = configMap["from"].(string)
		if !ok {
			return fmt.Errorf("invalid from type")
		}
	}

	if _, ok := configMap["path"]; ok {
		path, ok = configMap["path"].(string)
		if !ok {
			return fmt.Errorf("invalid path type")
		}
	}

	value = configMap["value"]

	patch, err := jsonpatch.NewJsonPatch(op, from, path, value)
	if err != nil {
		return err
	}

	obj, err := patch.Apply(res.Object)
	if err != nil {
		return err
	}

	objMap, ok := obj.(map[string]any)
	if !ok {
		return fmt.Errorf("invalid object type from patch")
	}

	res.Object = objMap
	return nil
}
