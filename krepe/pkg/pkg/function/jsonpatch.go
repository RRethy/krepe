package function

import (
	"fmt"

	"github.com/RRethy/krepe/jsonpatch"
	"github.com/RRethy/krepe/krepe/pkg/pkg/resource"
)

var (
	_ Function = &JsonPatch{}
)

type JsonPatch struct {
	op    string
	from  string
	path  string
	value any
	patch jsonpatch.JsonPatch
}

func (jp *JsonPatch) WithConfigMap(configMap map[string]any) (Function, error) {
	var op, from, path string
	var value any

	if _, ok := configMap["op"]; ok {
		op, ok = configMap["op"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid op type")
		}
	}

	if _, ok := configMap["from"]; ok {
		from, ok = configMap["from"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid from type")
		}
	}

	if _, ok := configMap["path"]; ok {
		path, ok = configMap["path"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid path type")
		}
	}

	value = configMap["value"]

	patch, err := jsonpatch.NewJsonPatch(op, from, path, value)
	if err != nil {
		return nil, fmt.Errorf("creating json patch: %v", err)
	}

	return &JsonPatch{
		op:    op,
		from:  from,
		path:  path,
		value: value,
		patch: patch,
	}, nil
}

func (jp *JsonPatch) Run(res *resource.Resource) error {
	obj, err := jp.patch.Apply(res.Object)
	if err != nil {
		return fmt.Errorf("applying patch: %v", err)
	}

	objMap, ok := obj.(map[string]any)
	if !ok {
		return fmt.Errorf("invalid object type from patch")
	}

	res.Object = objMap
	return nil
}
