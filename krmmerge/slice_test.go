package krmmerge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAssociativeKey(t *testing.T) {
	tests := []struct {
		name  string
		slice []any
		want  string
	}{
		{
			name:  "nil slice",
			slice: nil,
			want:  "",
		},
		{
			name:  "empty slice",
			slice: []any{},
			want:  "",
		},
		{
			name:  "non-associative slice",
			slice: []any{1, 2, 3},
			want:  "",
		},
		{
			name:  "associative slice with no associative key",
			slice: []any{map[string]any{"a": 1, "b": 2, "c": 3}},
			want:  "",
		},
		{
			name: "associative slice with associative key",
			slice: []any{
				map[string]any{"name": "foo", "a": 2, "b": 3, "c": 4},
				map[string]any{"name": "bar", "a": 2, "b": 3, "c": 4},
				map[string]any{"name": "baz", "a": 2, "b": 3, "c": 4},
			},
			want: "name",
		},
		{
			name: "associative slice with multiple associative keys with duplicate values",
			slice: []any{
				map[string]any{"name": "foo", "a": 2, "b": 3, "c": 4},
				map[string]any{"name": "bar", "a": 2, "b": 3, "c": 4, "d": 5},
				map[string]any{"name": "bar", "a": 2, "b": 3, "c": 4, "d": 5},
			},
			want: "",
		},
		{
			name: "associative slice with multiple associative keys and one non-associative key",
			slice: []any{
				map[string]any{"name": "foo", "a": 2, "b": 3, "c": 4},
				map[string]any{"name": "bar", "a": 2, "b": 3, "c": 4, "d": 5},
				3,
			},
			want: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := getAssociativeKey(test.slice)
			assert.Equal(t, test.want, got)
		})
	}
}
