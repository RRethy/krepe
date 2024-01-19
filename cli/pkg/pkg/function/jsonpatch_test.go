package function

import (
	"testing"

	"github.com/Shopify/krepe/jsonpatch"
)

func TestJsonPatchWithConfigMap(t *testing.T) {
	validJsonPatchAdd, _ := jsonpatch.NewAdd("/foo", "bar")
	validJsonPatchMove, _ := jsonpatch.NewMove("/foo/bar", "/foo/baz")
	runWithConfigMapTests(t, Function(&JsonPatch{}), []configMapTest{
		{
			name: "succeeds with valid `add` config map",
			configMap: map[string]any{
				"op":    "add",
				"path":  "/foo",
				"value": "bar",
			},
			wantFn: &JsonPatch{
				op:    "add",
				path:  "/foo",
				value: "bar",
				patch: validJsonPatchAdd,
			},
			wantErr: false,
		},
		{
			name: "succeeds with valid `move` config map",
			configMap: map[string]any{
				"op":   "move",
				"from": "/foo/bar",
				"path": "/foo/baz",
			},
			wantFn: &JsonPatch{
				op:    "move",
				from:  "/foo/bar",
				path:  "/foo/baz",
				patch: validJsonPatchMove,
			},
			wantErr: false,
		},
		{
			name: "fails with invalid config map",
			configMap: map[string]any{
				"foo": 1,
			},
			wantFn:  nil,
			wantErr: true,
		},
	})
}
