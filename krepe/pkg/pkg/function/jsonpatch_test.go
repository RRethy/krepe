package function

import (
	"testing"

	"github.com/Shopify/krepe/jsonpatch"
	"github.com/Shopify/krepe/krepe/pkg/pkg/resource"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
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

func TestJsonPatchRun(t *testing.T) {
	runRunTests(t, Function(&JsonPatch{}), []runTest{
		{
			name: "succeeds with valid `add` json patch",
			configMap: map[string]any{
				"op":    "add",
				"path":  "/foo",
				"value": "bar",
			},
			res: &resource.Resource{
				Unstructured: unstructured.Unstructured{
					Object: map[string]any{
						"foo": "baz",
					},
				},
			},
			validate: func(t *testing.T, res *resource.Resource) {
				assert.Equal(t, map[string]any{
					"foo": "bar",
				}, res.Object)
			},
			wantErr: false,
		},
		{
			name: "fails with invalid path for `add` json patch",
			configMap: map[string]any{
				"op":    "add",
				"path":  "/foo/bar",
				"value": "baz",
			},
			res: &resource.Resource{
				Unstructured: unstructured.Unstructured{
					Object: map[string]any{
						"foo": "baz",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "fails with invalid value for `add` json patch",
			configMap: map[string]any{
				"op":    "add",
				"path":  "",
				"value": 1,
			},
			res: &resource.Resource{
				Unstructured: unstructured.Unstructured{
					Object: map[string]any{
						"foo": "baz",
					},
				},
			},
			wantErr: true,
		},
	})
}
