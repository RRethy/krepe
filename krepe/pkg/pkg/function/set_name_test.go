package function

import (
	"testing"
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
