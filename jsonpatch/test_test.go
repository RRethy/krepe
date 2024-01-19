package jsonpatch

import (
	"testing"
)

func TestTest(t *testing.T) {
	runPatchTests(t, []*patchTest{
		{
			name: "test with valid path and valid value",
			obj: map[string]any{
				"foo": "bar",
				"baz": map[string]any{
					"baz2": map[string]any{
						"baz3": "qux",
					},
				},
			},
			patch: &Test{path: []string{"baz", "baz2"}, value: map[string]any{"baz3": "qux"}},
			want: map[string]any{
				"foo": "bar",
				"baz": map[string]any{
					"baz2": map[string]any{
						"baz3": "qux",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "test with valid path and invalid value",
			obj: map[string]any{
				"foo": "bar",
				"baz": map[string]any{
					"baz2": map[string]any{
						"baz3": "qux",
					},
				},
			},
			patch:   &Test{path: []string{"baz", "baz2"}, value: map[string]any{"baz3": "quux"}},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test with invalid path",
			obj: map[string]any{
				"foo": "bar",
				"baz": map[string]any{
					"baz2": map[string]any{
						"baz3": "qux",
					},
				},
			},
			patch:   &Test{path: []string{"baz", "baz3"}, value: "qux"},
			want:    nil,
			wantErr: true,
		},
	})
}
