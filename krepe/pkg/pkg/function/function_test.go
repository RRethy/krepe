package function

import (
	"testing"

	"github.com/Shopify/krepe/krepe/pkg/pkg/resource"
	"github.com/stretchr/testify/assert"
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
	res       *resource.Resource
	validate  func(t *testing.T, res *resource.Resource)
	wantErr   bool
}

func runWithConfigMapTests(t *testing.T, fn Function, tests []configMapTest) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn, err := fn.WithConfigMap(tt.configMap)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, fn)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantFn, fn)
			}
		})
	}
}

func runRunTests(t *testing.T, fn Function, tests []runTest) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn, err := fn.WithConfigMap(tt.configMap)
			assert.NoError(t, err)
			err = fn.Run(tt.res)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				tt.validate(t, tt.res)
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn, err := NewFunction(tt.fnName, tt.configMap)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, fn)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantFn, fn)
			}
		})
	}
}
