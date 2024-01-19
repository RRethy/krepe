package jsonpatch

import (
	"testing"
)

func TestCopy(t *testing.T) {
	runPatchTests(t, []*patchTest{
		{
			name: "copy from map",
			obj: map[string]any{
				"key1": "value1",
				"key2": "value2",
			},
			patch: &Copy{from: []string{"key1"}, path: []string{"key3"}},
			want: map[string]any{
				"key1": "value1",
				"key2": "value2",
				"key3": "value1",
			},
			wantErr: false,
		},
		{
			name: "copy from array to map element",
			obj: map[string]any{
				"key1": []any{"value1", "value2", "value3"},
			},
			patch: &Copy{from: []string{"key1", "1"}, path: []string{"key2"}},
			want: map[string]any{
				"key1": []any{"value1", "value2", "value3"},
				"key2": "value2",
			},
			wantErr: false,
		},
		{
			name: "copy with complex path",
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
			patch: &Copy{from: []string{"key1", "0", "key2"}, path: []string{"key1", "1", "key5"}},
			want: map[string]any{
				"key1": []any{
					map[string]any{
						"key2": "value2",
						"key3": "value3",
					},
					map[string]any{
						"key4": "value4",
						"key5": "value2",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "copy complex object",
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
			patch: &Copy{from: []string{"key1"}, path: []string{"key5"}},
			want: map[string]any{
				"key1": []any{
					map[string]any{
						"key2": "value2",
						"key3": "value3",
					},
					map[string]any{
						"key4": "value4",
					},
				},
				"key5": []any{
					map[string]any{
						"key2": "value2",
						"key3": "value3",
					},
					map[string]any{
						"key4": "value4",
					},
				},
			},
		},
		{
			name: "copy from non-existent path in map",
			obj: map[string]any{
				"key1": "value1",
				"key2": "vobjalue2",
			},
			patch: &Copy{from: []string{"key3"}, path: []string{"key4"}},
			want: map[string]any{
				"key1": "value1",
				"key2": "value2",
			},
			wantErr: true,
		},
		{
			name: "copy from non-existent path in array",
			obj: map[string]any{
				"key1": []any{"value1", "value2", "value3"},
			},
			patch: &Copy{from: []string{"key1", "3"}, path: []string{"key2"}},
			want: map[string]any{
				"key1": []any{"value1", "value2", "value3"},
			},
			wantErr: true,
		},
		{
			name: "copy to non-existent path",
			obj: map[string]any{
				"key1": "value1",
				"key2": "value2",
			},
			patch: &Copy{from: []string{"key1"}, path: []string{"key3", "key4"}},
			want: map[string]any{
				"key1": "value1",
				"key2": "value2",
			},
			wantErr: true,
		},
	})
}
