package jsonpatch

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReplace(t *testing.T) {
	tests := []struct {
		name    string
		obj     map[string]any
		op      string
		path    string
		value   any
		want    any
		wantErr bool
	}{
		{
			name: "replace from map",
			obj: map[string]any{
				"key1": "value1",
				"key2": "value2",
			},
			path:  "/key1",
			value: "value3",
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
			path:  "/key1/1",
			value: "value4",
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
			path:  "/key1/0/key2",
			value: "value5",
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
			path:  "/key3",
			value: "value3",
			want: map[string]any{
				"key1": "value1",
				"key2": "value2",
			},
			wantErr: true,
		},
		{
			name: "replace index out of bounds from array",
			obj: map[string]any{
				"key1": []any{"value1", "value2", "value3"},
			},
			path:  "/key1/3",
			value: "value4",
			want: map[string]any{
				"key1": []any{"value1", "value2", "value3"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ptrs, _ := pathToJsonPtrs(tt.path)
			newObj, err := replace(tt.obj, ptrs, tt.value)
			assert.Equal(t, tt.want, newObj)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
