package functions

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestSetNamespaceWithConfigMap(t *testing.T) {
	runWithConfigMapTests(t, Function(&SetNamespace{}), []configMapTest{
		{
			name: "succeeds with valid config map",
			configMap: map[string]any{
				"namespace": "foo",
			},
			wantFn: &SetNamespace{
				namespace: "foo",
			},
			wantErr: false,
		},
		{
			name: "fails with invalid config map",
			configMap: map[string]any{
				"namespace": 1,
			},
			wantFn:  nil,
			wantErr: true,
		},
	})
}

func TestSetNamespaceRun(t *testing.T) {
	runRunTests(t, Function(&SetNamespace{}), []runTest{
		{
			name: "succeeds with valid set namespace",
			configMap: map[string]any{
				"namespace": "foo",
			},
			res: &unstructured.Unstructured{
				Object: map[string]any{
					"namespace": "bar",
				},
			},
			validate: func(t *testing.T, res *unstructured.Unstructured) {
				assert.Equal(t, "foo", res.GetNamespace())
			},
			wantErr: false,
		},
	})
}
