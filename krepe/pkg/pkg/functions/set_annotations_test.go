package functions

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestSetAnnotationsWithConfigMap(t *testing.T) {
	runWithConfigMapTests(t, Function(&SetAnnotations{}), []configMapTest{
		{
			name: "succeeds with valid config map",
			configMap: map[string]any{
				"foo": "bar",
			},
			wantFn: &SetAnnotations{
				annotations: map[string]string{
					"foo": "bar",
				},
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

func TestSetAnnotationsRun(t *testing.T) {
	runRunTests(t, Function(&SetAnnotations{}), []runTest{
		{
			name: "succeeds with valid set annotations",
			configMap: map[string]any{
				"foo": "bar",
			},
			res: &unstructured.Unstructured{
				Object: map[string]any{
					"metadata": map[string]any{
						"annotations": map[string]any{
							"foo": "baz",
							"bar": "baz",
						},
					},
				},
			},
			validate: func(t *testing.T, res *unstructured.Unstructured) {
				assert.Equal(t, map[string]string{
					"foo": "bar",
				}, res.GetAnnotations())
			},
			wantErr: false,
		},
	})
}
