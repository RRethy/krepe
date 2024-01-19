package jsonpatch

import (
	"testing"
)

func TestRemove(t *testing.T) {
	runPatchTests(t, []*patchTest{
		{
			name: "remove from map",
			obj: map[string]any{
				"key1": "value1",
				"key2": "value2",
			},
			patch: &Remove{path: []string{"key1"}},
			want: map[string]any{
				"key2": "value2",
			},
			wantErr: false,
		},
		{
			name: "remove from array",
			obj: map[string]any{
				"key1": []any{"value1", "value2", "value3"},
			},
			patch: &Remove{path: []string{"key1", "1"}},
			want: map[string]any{
				"key1": []any{"value1", "value3"},
			},
			wantErr: false,
		},
		{
			name: "remove with complex path",
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
			patch: &Remove{path: []string{"key1", "0", "key2"}},
			want: map[string]any{
				"key1": []any{
					map[string]any{
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
			name: "remove non-existing key from map",
			obj: map[string]any{
				"key1": "value1",
				"key2": "value2",
			},
			patch:   &Remove{path: []string{"key3"}},
			want:    nil,
			wantErr: true,
		},
		{
			name: "remove index out of bounds from array",
			obj: map[string]any{
				"key1": []any{"value1", "value2", "value3"},
			},
			patch:   &Remove{path: []string{"key1", "3"}},
			want:    nil,
			wantErr: true,
		},
	})
}
