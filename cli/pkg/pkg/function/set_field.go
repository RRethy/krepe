package function

import (
	"github.com/Shopify/krepe/cli/pkg/pkg/resource"
)

type SetField struct{}

func (f *SetField) Run(res *resource.Resource, configMap map[string]any) error {
	return nil
	// op, ok := configMap["op"].(string)
	// if !ok {
	// 	return fmt.Errorf("failed to get a key `op` with type `string` form configMap")
	// }
	// path, ok := configMap["path"].(string)
	// if !ok {
	// 	return fmt.Errorf("failed to get a key `path` with type `string` form configMap")
	// }
	// value, ok := configMap["value"].(string)
	// if !ok {
	// 	return fmt.Errorf("failed to get a key `value` with type `string` form configMap")
	// }
	//
	// return doJsonPatch(res, op, path, value)
}

// TODO(RRethy): This needs to be cleaned up and put into its own module
// func doJsonPatch(res *resource.Resource, op string, rawPath string, value string) error {
// 	path := strings.Split(rawPath, "/")
// 	if len(path) == 0 {
// 		return fmt.Errorf("empty path is unsupported")
// 	}
// 	for i, p := range path {
// 		p = strings.ReplaceAll(p, "~1", "/")
// 		p = strings.ReplaceAll(p, "~0", "~")
// 		path[i] = p
// 	}
//
// 	obj := res.Object
//
// 	switch op {
// 	case "add":
// 		newObj, err := addJsonPatch(obj, path, value)
// 		if err != nil {
// 			return err
// 		}
// 		res.Object = newObj
// 		return nil
// 	case "remove":
// 		return removeJsonPatch(res, path)
// 	case "replace":
// 		return replaceJsonPatch(res, path, value)
// 	case "move":
// 		return moveJsonPatch(res, path, value)
// 	case "copy":
// 		return copyJsonPatch(res, path, value)
// 	case "test":
// 		return testJsonPatch(res, path, value)
// 	default:
// 		return fmt.Errorf("invalid op: %s", op)
// 	}
// }
//
// func addJsonPatch(obj map[string]any, path []string, value string) (map[string]any, error) {
// 	n, err := strconv.ParseFloat(path[0], 64)
// 	if err == nil {
//
// 	} else {
// 		if len(path) == 1 {
// 			obj[path[0]] = value
// 			return obj, nil
// 		}
// 		if _, ok := obj[path[0]]; !ok {
// 			obj[path[0]] = make(map[string]any, 1)
// 		}
// 		if v, ok := obj[path[0]].(map[string]any); ok {
// 			newObj, err := addJsonPatch(v, path[1:], value)
// 			if err != nil {
// 				return nil, err
// 			}
// 			obj[path[0]] = newObj
// 			return obj, nil
// 		} else {
// 			return nil, fmt.Errorf("failed to add a new field to a non-object %s(%s)", path[0], obj[path[0]])
// 		}
// 	}
// }
//
// func removeJsonPatch(res *resource.Resource, path []string) error {
// 	// TODO
// 	return nil
// }
//
// func replaceJsonPatch(res *resource.Resource, path []string, value string) error {
// 	// TODO
// 	return nil
// }
//
// func moveJsonPatch(res *resource.Resource, path []string, value string) error {
// 	// TODO
// 	return nil
// }
//
// func copyJsonPatch(res *resource.Resource, path []string, value string) error {
// 	// TODO
// 	return nil
// }
//
// func testJsonPatch(res *resource.Resource, path []string, value string) error {
// 	// TODO
// 	return nil
// }
