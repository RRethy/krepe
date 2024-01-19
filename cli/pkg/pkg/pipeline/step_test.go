package pipeline

import (
	"testing"

	"github.com/Shopify/krepe/cli/pkg/pkg/function"
	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert"
)

func TestStepUnmarshallYAML(t *testing.T) {
	setLabelsFn, _ := function.NewFunction("set-labels", map[string]any{
		"foo": "bar",
	})

	tests := []struct {
		name     string
		yaml     string
		wantName string
		wantFn   function.Function
		wantErr  bool
	}{
		{
			name: "succeeds with valid step",
			yaml: `
function: set-labels
configMap:
  foo: bar
			`,
			wantName: "set-labels",
			wantFn:   setLabelsFn,
			wantErr:  false,
		},
		{
			name:     "fails with invalid step yaml",
			yaml:     `foo: bar`,
			wantName: "",
			wantFn:   nil,
			wantErr:  true,
		},
		{
			name: "fails with invalid function",
			yaml: `
function: unknown-fn
configMap:
  foo: bar
			`,
			wantName: "",
			wantFn:   nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Step{}
			err := yaml.Unmarshal([]byte(tt.yaml), s)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, s.fn)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantName, s.name)
				assert.Equal(t, tt.wantFn, s.fn)
			}
		})
	}
}
