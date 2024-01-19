package function

import (
	"testing"
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
