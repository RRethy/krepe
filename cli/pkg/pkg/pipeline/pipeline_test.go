package pipeline

import (
	"testing"

	"github.com/Shopify/krepe/cli/pkg/pkg/function"
	"github.com/Shopify/krepe/cli/pkg/pkg/resource"
	"github.com/stretchr/testify/assert"
)

func TestPipelineRun(t *testing.T) {
	setLabelsFn, _ := function.NewFunction("set-labels", map[string]any{
		"foo": "bar",
	})
	setAnnotationsFn, _ := function.NewFunction("set-annotations", map[string]any{
		"baz": "qux",
	})
	badjsonPatchFn, _ := function.NewFunction("jsonpatch", map[string]any{
		"op":    "add",
		"path":  "/spec/foo/bar",
		"value": 1,
	})
	startResource, _ := resource.NewResourceFromBytes("foo.yaml", []byte(``))

	tests := []struct {
		name          string
		pipeline      Pipeline
		inputResource *resource.Resource
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.pipeline.Run(tt.inputResource)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantObject, tt.inputResource.Object)
			}
		})
	}
}
