package pipeline

import (
	"testing"

	"github.com/RRethy/krepe/krepe/pkg/pkg/function"
	"github.com/RRethy/krepe/krepe/pkg/pkg/resource"
	"github.com/RRethy/krepe/krepe/pkg/yaml"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestStepUnmarshallYAML(t *testing.T) {
	setLabelsFn, _ := function.NewFunction("set-labels", map[string]any{
		"foo": "bar",
	})

	tests := []struct {
		name       string
		yaml       string
		wantName   string
		wantFn     function.Function
		wantTarget *Target
		wantErr    bool
	}{
		{
			name: "succeeds with valid step",
			yaml: `
function: set-labels
target:
  kind: Deployment
  name: foo
configMap:
  foo: bar
`,
			wantName: "set-labels",
			wantFn:   setLabelsFn,
			wantTarget: &Target{
				Kind: "Deployment",
				Name: "foo",
			},
			wantErr: false,
		},
		{
			name: "succeeds with empty target",
			yaml: `
function: set-labels
configMap:
  foo: bar
`,
			wantName:   "set-labels",
			wantFn:     setLabelsFn,
			wantTarget: nil,
			wantErr:    false,
		},
		{
			name:       "fails with invalid step yaml",
			yaml:       `foo: bar`,
			wantName:   "",
			wantFn:     nil,
			wantTarget: nil,
			wantErr:    true,
		},
		{
			name: "fails with invalid function",
			yaml: `
function: unknown-fn
configMap:
  foo: bar
			`,
			wantName:   "",
			wantFn:     nil,
			wantTarget: nil,
			wantErr:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &Step{}
			err := yaml.Unmarshal([]byte(test.yaml), s)
			if test.wantErr {
				assert.Error(t, err)
				assert.Nil(t, s.fn)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.wantName, s.name)
				assert.Equal(t, test.wantFn, s.fn)
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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.step.Run(test.inputResource)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want, test.inputResource.Object)
			}
		})
	}
}

func TestStepMatches(t *testing.T) {
	tests := []struct {
		name   string
		step   *Step
		resObj map[string]any
		want   bool
	}{
		{
			name: "succeeds with matching target",
			step: &Step{
				target: &Target{
					Kind: "Pod",
				},
			},
			resObj: map[string]any{
				"apiVersion": "v1",
				"kind":       "Pod",
				"metadata": map[string]any{
					"name":      "test-pod",
					"namespace": "default",
				},
			},
			want: true,
		},
		{
			name: "fails with non-matching target",
			step: &Step{
				target: &Target{
					Kind: "Pod",
				},
			},
			resObj: map[string]any{
				"apiVersion": "v1",
				"kind":       "Deployment",
			},
			want: false,
		},
		{
			name: "succeeds with empty target",
			step: &Step{
				target: nil,
			},
			resObj: map[string]any{
				"apiVersion": "v1",
				"kind":       "Deployment",
			},
			want: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.step.Matches(&resource.Resource{
				Unstructured: unstructured.Unstructured{
					Object: test.resObj,
				},
			})
			assert.Equal(t, test.want, got)
		})
	}
}
