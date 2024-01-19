package functions

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type configMapTest struct {
	name      string
	configMap map[string]any
	wantFn    Function
	wantErr   bool
}

type runTest struct {
	name      string
	configMap map[string]any
	res       *unstructured.Unstructured
	validate  func(t *testing.T, res *unstructured.Unstructured)
	wantErr   bool
}

func runWithConfigMapTests(t *testing.T, fn Function, tests []configMapTest) {
	t.Helper()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fn, err := fn.WithConfigMap(test.configMap)
			if test.wantErr {
				assert.Error(t, err)
				assert.Nil(t, fn)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.wantFn, fn)
			}
		})
	}
}

func runRunTests(t *testing.T, fn Function, tests []runTest) {
	t.Helper()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fn, err := fn.WithConfigMap(test.configMap)
			assert.NoError(t, err)
			err = fn.Run(test.res)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				test.validate(t, test.res)
			}
		})
	}
}

func TestNewFunction(t *testing.T) {
	tests := []struct {
		name      string
		fnName    string
		configMap map[string]any
		wantFn    Function
		wantErr   bool
	}{
		{
			name:   "succeeds with valid config map",
			fnName: "set-annotations",
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
			name:   "fails with invalid function name",
			fnName: "invalid",
			configMap: map[string]any{
				"foo": "bar",
			},
			wantFn:  nil,
			wantErr: true,
		},
		{
			name:   "fails with invalid config map",
			fnName: "set-annotations",
			configMap: map[string]any{
				"foo": 1,
			},
			wantFn:  nil,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fn, err := NewFunction(test.fnName, test.configMap)
			if test.wantErr {
				assert.Error(t, err)
				assert.Nil(t, fn)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.wantFn, fn)
			}
		})
	}
}
