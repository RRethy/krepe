package functions

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestAddLabelsWithConfigMap(t *testing.T) {
	runWithConfigMapTests(t, Function(&AddLabels{}), []configMapTest{
		{
			name:      "succeeds with valid config map",
			configMap: map[string]any{"foo": "bar"},
			wantFn: &AddLabels{
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

func TestAddLabelsRun(t *testing.T) {
	runRunTests(t, Function(&AddLabels{}), []runTest{
		{
			name: "succeeds with valid add labels",
			configMap: map[string]any{
				"foo": "bar",
			},
			res: &unstructured.Unstructured{
				Object: map[string]any{
					"metadata": map[string]any{
						"labels": map[string]any{
							"foo": "baz",
							"bar": "baz",
						},
					},
				},
			},
			validate: func(t *testing.T, res *unstructured.Unstructured) {
				assert.Equal(t, map[string]string{
					"foo": "bar",
					"bar": "baz",
				}, res.GetLabels())
			},
			wantErr: false,
		},
		{
			name: "succeeds with valid add labels and no labels",
			configMap: map[string]any{
				"foo": "bar",
			},
			res: &unstructured.Unstructured{
				Object: map[string]any{
					"metadata": map[string]any{
						"labels": nil,
					},
				},
			},
			validate: func(t *testing.T, res *unstructured.Unstructured) {
				assert.Equal(t, map[string]string{
					"foo": "bar",
				}, res.GetLabels())
			},
		},
	})
}
