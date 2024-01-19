package function

import (
	"testing"
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
