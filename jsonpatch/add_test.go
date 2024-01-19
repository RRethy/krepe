package jsonpatch

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJSONPatch(t *testing.T) {
	tests := []struct {
		name    string
		obj     map[string]any
		op      string
		path    string
		value   any
		want    map[string]any
		wantErr bool
	}{
		{
			name: "add to map",
			obj: map[string]any{
				"foo": "bar",
			},
			op:    "add",
			path:  "/baz",
			value: "qux",
			want: map[string]any{
				"foo": "bar",
				"baz": "qux",
			},
			wantErr: false,
		},
		{
			name: "add to map long path",
			obj: map[string]any{
				"foo": "bar",
				"baz": map[string]any{
					"baz2": "qux",
					"baz3": map[string]any{
						"baz4": "quux",
					},
				},
			},
			op:    "add",
			path:  "/baz/baz3/baz5",
			value: "corge",
			want: map[string]any{
				"foo": "bar",
				"baz": map[string]any{
					"baz2": "qux",
					"baz3": map[string]any{
						"baz4": "quux",
						"baz5": "corge",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "add to map with non existent ptr",
			obj: map[string]any{
				"foo": "bar",
			},
			op:    "add",
			path:  "/baz/baz2",
			value: "qux",
			want: map[string]any{
				"foo": "bar",
				"baz": map[string]any{
					"baz2": "qux",
				},
			},
			wantErr: false,
		},
		{
			name: "add to map with non existent ptr long path",
			obj: map[string]any{
				"foo": "bar",
			},
			op:    "add",
			path:  "/baz/baz2/baz3/baz4/baz5",
			value: "qux",
			want: map[string]any{
				"foo": "bar",
				"baz": map[string]any{
					"baz2": map[string]any{
						"baz3": map[string]any{
							"baz4": map[string]any{
								"baz5": "qux",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "add to array",
			obj: map[string]any{
				"foo": []any{
					"bar",
				},
			},
			op:    "add",
			path:  "/foo/1",
			value: "qux",
			want: map[string]any{
				"foo": []any{
					"bar",
					"qux",
				},
			},
			wantErr: false,
		},
		{
			name: "add to array with -",
			obj: map[string]any{
				"foo": []any{
					"bar",
				},
			},
			op:    "add",
			path:  "/foo/-",
			value: "qux",
			want: map[string]any{
				"foo": []any{
					"bar",
					"qux",
				},
			},
			wantErr: false,
		},
		{
			name: "add to array overwrite first element",
			obj: map[string]any{
				"foo": []any{
					"bar",
				},
			},
			op:    "add",
			path:  "/foo/0",
			value: "qux",
			want: map[string]any{
				"foo": []any{
					"qux",
				},
			},
			wantErr: false,
		},
		{
			name: "add to array with index out of range",
			obj: map[string]any{
				"foo": []any{
					"bar",
				},
			},
			op:    "add",
			path:  "/foo/2",
			value: "qux",
			want: map[string]any{
				"foo": []any{
					"bar",
					nil,
					"qux",
				},
			},
			wantErr: false,
		},
		{
			name: "add to array with large index out of range",
			obj: map[string]any{
				"foo": []any{
					"bar",
				},
			},
			op:    "add",
			path:  "/foo/3",
			value: "qux",
			want: map[string]any{
				"foo": []any{
					"bar",
					nil,
					nil,
					"qux",
				},
			},
			wantErr: false,
		},
		{
			name: "add with arrays and maps",
			obj: map[string]any{
				"foo": []any{
					"bar",
					map[string]any{
						"baz": "qux",
						"quux": []any{
							"corge",
							"grault",
						},
					},
					"baz",
				},
			},
			op:    "add",
			path:  "/foo/1/quux/1",
			value: "qux",
			want: map[string]any{
				"foo": []any{
					"bar",
					map[string]any{
						"baz": "qux",
						"quux": []any{
							"corge",
							"qux",
						},
					},
					"baz",
				},
			},
			wantErr: false,
		},
		{
			name: "add to array fails with invalid path",
			obj: map[string]any{
				"foo": 3,
			},
			op:    "add",
			path:  "/foo/bar",
			value: "baz",
			want: map[string]any{
				"foo": 3,
			},
			wantErr: true,
		},
		{
			name: "add to array fails with invalid path with array",
			obj: map[string]any{
				"foo": []any{
					"bar",
				},
			},
			op:    "add",
			path:  "/foo/bar",
			value: "baz",
			want: map[string]any{
				"foo": []any{
					"bar",
				},
			},
			wantErr: true,
		},
		{
			name: "add fails if expecting map",
			obj: map[string]any{
				"foo": []any{
					"bar",
				},
			},
			op:    "add",
			path:  "/foo/bar",
			value: "baz",
			want: map[string]any{
				"foo": []any{
					"bar",
				},
			},
			wantErr: true,
		},
		{
			name: "add fails if expecting array",
			obj: map[string]any{
				"foo": map[string]any{
					"bar": "baz",
				},
			},
			op:    "add",
			path:  "/foo/1",
			value: "baz",
			want: map[string]any{
				"foo": map[string]any{
					"bar": "baz",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ptrs, _ := pathToJsonPtrs(tt.path)
			newObj, err := add(tt.obj, ptrs, tt.value)
			assert.Equal(t, tt.want, newObj)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
