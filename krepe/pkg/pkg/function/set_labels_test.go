package function

import (
	"testing"
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
