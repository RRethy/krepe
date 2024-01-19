package pkg

import (
	"testing"

	"github.com/RRethy/krepe/krepe/pkg/pkg/functions"
	"github.com/stretchr/testify/assert"
)

func TestPipelineRun(t *testing.T) {
	setLabelsFn, _ := functions.NewFunction("set-labels", map[string]any{
		"foo": "bar",
	})
	setAnnotationsFn, _ := functions.NewFunction("set-annotations", map[string]any{
		"baz": "qux",
	})
	badjsonPatchFn, _ := functions.NewFunction("jsonpatch", map[string]any{
		"op":    "add",
		"path":  "/spec/foo/bar",
		"value": 1,
	})
	startResource, _ := NewResourceFromBytes("foo.yaml", []byte(``))

	tests := []struct {
		name          string
		pipeline      Pipeline
		inputResource *Resource
		wantObject    map[string]any
		wantErr       bool
	}{
		{
			name: "succeeds with valid pipeline steps",
			pipeline: []*Step{
				{
					name: "set-labels",
					fn:   setLabelsFn,
				},
				{
					name: "set-annotations",
					fn:   setAnnotationsFn,
				},
			},
			inputResource: startResource,
			wantObject: map[string]any{
				"metadata": map[string]any{
					"labels": map[string]any{
						"foo": "bar",
					},
					"annotations": map[string]any{
						"baz": "qux",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "fails with any invalid pipeline step",
			pipeline: []*Step{
				{
					name: "set-labels",
					fn:   setLabelsFn,
				},
				{
					name: "jsonpatch",
					fn:   badjsonPatchFn,
				},
			},
			inputResource: startResource,
			wantObject:    nil,
			wantErr:       true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.pipeline.Run(test.inputResource)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.wantObject, test.inputResource.Object)
			}
		})
	}
}
