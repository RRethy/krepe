package jsonpatch

import (
	"testing"
)

func TestMove(t *testing.T) {
	runPatchTests(t, []*patchTest{
		{
			name: "move valid from to valid path",
			obj: map[string]any{
				"key1": "value1",
				"key2": "value2",
			},
			patch: &Move{from: []string{"key1"}, path: []string{"key3"}},
			want: map[string]any{
				"key2": "value2",
				"key3": "value1",
			},
			wantErr: false,
		},
		{
			name: "move from valid from to invalid path",
			obj: map[string]any{
				"key1": "value1",
				"key2": "value2",
			},
			patch:   &Move{from: []string{"key1"}, path: []string{"key3", "key4"}},
			want:    nil,
			wantErr: true,
		},
		{
			name: "move from invalid from to valid path",
			obj: map[string]any{
				"key1": "value1",
				"key2": "value2",
			},
			patch:   &Move{from: []string{"key3"}, path: []string{"key4"}},
			want:    nil,
			wantErr: true,
		},
		{
			name: "move from invalid from to invalid path",
			obj: map[string]any{
				"key1": "value1",
				"key2": "value2",
			},
			patch:   &Move{from: []string{"key3"}, path: []string{"key4", "key5"}},
			want:    nil,
			wantErr: true,
		},
	})
}
