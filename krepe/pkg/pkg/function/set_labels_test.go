package function

import (
	"testing"

	"github.com/Shopify/krepe/krepe/pkg/pkg/resource"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestSetLabelsWithConfigMap(t *testing.T) {
	runWithConfigMapTests(t, Function(&SetLabels{}), []configMapTest{
		{
			name:      "succeeds with valid config map",
			configMap: map[string]any{"foo": "bar"},
			wantFn: &SetLabels{
				labels: map[string]string{"foo": "bar"},
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

func TestSetLabelsRun(t *testing.T) {
	runRunTests(t, Function(&SetLabels{}), []runTest{
		{
			name: "succeeds with valid set labels",
			configMap: map[string]any{
				"foo": "bar",
			},
			res: &resource.Resource{
				Unstructured: unstructured.Unstructured{
					Object: map[string]any{
						"metadata": map[string]any{
							"labels": map[string]any{
								"foo": "baz",
								"bar": "baz",
							},
						},
					},
				},
			},
			validate: func(t *testing.T, res *resource.Resource) {
				assert.Equal(t, map[string]string{
					"foo": "bar",
				}, res.GetLabels())
			},
			wantErr: false,
		},
	})
}
