package jsonpatch

import (
	"testing"
)

func TestAdd(t *testing.T) {
	runPatchTests(t, []*patchTest{
		{
			name: "add with empty path",
			obj: map[string]any{
				"foo": "bar",
				"baz": map[string]any{
					"baz2": "qux",
				},
			},
			patch:   &Add{path: []string{}, value: "quux"},
			want:    "quux",
			wantErr: false,
		},
		{
			name: "add with empty json path",
			obj: map[string]any{
				"foo": "bar",
				"baz": map[string]any{
					"baz2": "qux",
				},
			},
			patch: &Add{path: []string{"baz", ""}, value: "quux"},
			want: map[string]any{
				"foo": "bar",
				"baz": map[string]any{
					"baz2": "qux",
					"":     "quux",
				},
			},
			wantErr: false,
		},
		{
			name: "add to map",
			obj: map[string]any{
				"foo": "bar",
			},
			patch: &Add{path: []string{"baz"}, value: "qux"},
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
			patch: &Add{path: []string{"baz", "baz3", "baz5"}, value: "corge"},
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
			name: "add to map with non existent path",
			obj: map[string]any{
				"foo": "bar",
			},
			patch:   &Add{path: []string{"baz", "baz2"}, value: "qux"},
			want:    nil,
			wantErr: true,
		},
		{
			name: "add to array",
			obj: map[string]any{
				"foo": []any{
					"bar",
				},
			},
			patch: &Add{path: []string{"foo", "1"}, value: "qux"},
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
			patch: &Add{path: []string{"foo", "-"}, value: "qux"},
			want: map[string]any{
				"foo": []any{
					"bar",
					"qux",
				},
			},
			wantErr: false,
		},
		{
			name: "add to array doesn't overwrite first element",
			obj: map[string]any{
				"foo": []any{
					"bar",
				},
			},
			patch: &Add{path: []string{"foo", "0"}, value: "qux"},
			want: map[string]any{
				"foo": []any{
					"qux",
					"bar",
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
			patch:   &Add{path: []string{"foo", "2"}, value: "qux"},
			want:    nil,
			wantErr: true,
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
			patch: &Add{path: []string{"foo", "1", "quux", "1"}, value: "qux"},
			want: map[string]any{
				"foo": []any{
					"bar",
					map[string]any{
						"baz": "qux",
						"quux": []any{
							"corge",
							"qux",
							"grault",
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
			patch:   &Add{path: []string{"foo", "bar"}, value: "baz"},
			want:    nil,
			wantErr: true,
		},
		{
			name: "add to array fails with invalid path with array",
			obj: map[string]any{
				"foo": []any{
					"bar",
				},
			},
			patch:   &Add{path: []string{"foo", "bar"}, value: "baz"},
			want:    nil,
			wantErr: true,
		},
		{
			name: "add allows integer keys",
			obj: map[string]any{
				"foo": map[string]any{
					"bar": "baz",
				},
			},
			patch: &Add{path: []string{"foo", "1"}, value: "baz"},
			want: map[string]any{
				"foo": map[string]any{
					"bar": "baz",
					"1":   "baz",
				},
			},
			wantErr: false,
		},
		{
			name: "add fails if key not in map",
			obj: map[string]any{
				"foo": map[string]any{
					"bar": "baz",
				},
			},
			patch:   &Add{path: []string{"foo", "baz", "qux"}, value: "qux"},
			want:    nil,
			wantErr: true,
		},
	})
}
