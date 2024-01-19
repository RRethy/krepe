package function

import (
	"testing"

	"github.com/Shopify/krepe/krepe/pkg/pkg/resource"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestSetNameWithConfigMap(t *testing.T) {
	runWithConfigMapTests(t, Function(&SetName{}), []configMapTest{
		{
			name:      "succeeds with valid config map",
			configMap: map[string]any{"name": "bar"},
			wantFn: &SetName{
				name: "bar",
			},
			wantErr: false,
		},
		{
			name:      "fails with invalid config map",
			configMap: map[string]any{"foo": 1},
			wantFn:    nil,
			wantErr:   true,
		},
	})
}

func TestSetNameRun(t *testing.T) {
	runRunTests(t, Function(&SetName{}), []runTest{
		{
			name: "succeeds with valid set name",
			configMap: map[string]any{
				"name": "foo",
			},
			res: &resource.Resource{
				Unstructured: unstructured.Unstructured{
					Object: map[string]any{
						"name": "bar",
					},
				},
			},
			validate: func(t *testing.T, res *resource.Resource) {
				assert.Equal(t, "foo", res.GetName())
			},
			wantErr: false,
		},
	})
}
