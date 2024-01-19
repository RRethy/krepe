package jsonpatch

import (
	"testing"
)

func TestReplace(t *testing.T) {
	runPatchTests(t, []*patchTest{
		{
			name: "replace from map",
			obj: map[string]any{
				"key1": "value1",
				"key2": "value2",
			},
			patch: &Replace{path: []string{"key1"}, value: "value3"},
			want: map[string]any{
				"key1": "value3",
				"key2": "value2",
			},
			wantErr: false,
		},
		{
			name: "replace from array",
			obj: map[string]any{
				"key1": []any{"value1", "value2", "value3"},
			},
			patch: &Replace{path: []string{"key1", "1"}, value: "value4"},
			want: map[string]any{
				"key1": []any{"value1", "value4", "value3"},
			},
			wantErr: false,
		},
		{
			name: "replace with complex path",
			obj: map[string]any{
				"key1": []any{
					map[string]any{
						"key2": "value2",
						"key3": "value3",
					},
					map[string]any{
						"key4": "value4",
					},
				},
			},
			patch: &Replace{path: []string{"key1", "0", "key2"}, value: "value5"},
			want: map[string]any{
				"key1": []any{
					map[string]any{
						"key2": "value5",
						"key3": "value3",
					},
					map[string]any{
						"key4": "value4",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "replace non-existing key from map",
			obj: map[string]any{
				"key1": "value1",
				"key2": "value2",
			},
			patch:   &Replace{path: []string{"key3"}, value: "value3"},
			want:    nil,
			wantErr: true,
		},
		{
			name: "replace index out of bounds from array",
			obj: map[string]any{
				"key1": []any{"value1", "value2", "value3"},
			},
			patch:   &Replace{path: []string{"key1", "3"}, value: "value4"},
			want:    nil,
			wantErr: true,
		},
	})
}
