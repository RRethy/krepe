package pipeline

import (
	"testing"

	"github.com/Shopify/krepe/cli/pkg/pkg/function"
	"github.com/Shopify/krepe/cli/pkg/pkg/resource"
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

func TestStepMarshalYAML(t *testing.T) {
	setLabelsFn, _ := function.NewFunction("set-labels", map[string]any{
		"foo": "bar",
	})
	step := &Step{
		name: "set-labels",
		fn:   setLabelsFn,
		configMap: map[string]any{
			"foo": "bar",
		},
	}

	got, err := yaml.Marshal(step)
	assert.NoError(t, err)
	assert.Equal(t, `function: set-labels
configMap:
  foo: bar
`, string(got))
}

func TestStepRun(t *testing.T) {
	setLabelsFn, _ := function.NewFunction("set-labels", map[string]any{
		"foo": "bar",
	})
	badjsonPatchFn, _ := function.NewFunction("jsonpatch", map[string]any{
		"op":    "add",
		"path":  "/foo/bar",
		"value": "baz",
	})
	inputResource, _ := resource.NewResourceFromBytes("foo.yaml", []byte(``))

	tests := []struct {
		name          string
		step          *Step
		inputResource *resource.Resource
		want          map[string]any
		wantErr       bool
	}{
		{
			name: "succeeds with valid step",
			step: &Step{
				name: "set-labels",
				fn:   setLabelsFn,
			},
			inputResource: inputResource,
			want: map[string]any{
				"metadata": map[string]any{
					"labels": map[string]any{
						"foo": "bar",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "fails with run error step",
			step: &Step{
				name: "jsonpatch",
				fn:   badjsonPatchFn,
			},
			inputResource: inputResource,
			want:          nil,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.step.Run(tt.inputResource)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, tt.inputResource.Object)
			}
		})
	}
}
