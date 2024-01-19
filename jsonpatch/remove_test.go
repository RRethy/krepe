package jsonpatch

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Testremove(t *testing.T) {
	tests := []struct {
		name    string
		obj     map[string]any
		op      string
		path    string
		want    any
		wantErr bool
	}{
		{
			name: "remove from map",
			obj: map[string]any{
				"key1": "value1",
				"key2": "value2",
			},
			path: "/key1",
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
			path: "/key1/1",
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
			path: "/key1/0/key2",
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
			path: "/key3",
			want: map[string]any{
				"key1": "value1",
				"key2": "value2",
			},
			wantErr: true,
		},
		{
			name: "remove index out of bounds from array",
			obj: map[string]any{
				"key1": []any{"value1", "value2", "value3"},
			},
			path: "/key1/3",
			want: map[string]any{
				"key1": []any{"value1", "value2", "value3"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ptrs, _ := pathToJsonPtrs(tt.path)
			newObj, err := remove(tt.obj, ptrs)
			assert.Equal(t, tt.want, newObj)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
