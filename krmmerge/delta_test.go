package krmmerge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDelta(t *testing.T) {
	tests := []struct {
		name   string
		source any
		remove any
		want   any
	}{
		{
			name:   "map with matching keys but different values",
			source: map[string]any{"a": 1, "b": 2, "c": 3},
			remove: map[string]any{"b": 3},
			want:   map[string]any{"a": 1, "b": 2, "c": 3},
		},
		{
			name:   "map with matching keys and value",
			source: map[string]any{"a": 1, "b": 2, "c": 3},
			remove: map[string]any{"b": 2, "d": 4},
			want:   map[string]any{"a": 1, "c": 3},
		},
		{
			name:   "map with different keys",
			source: map[string]any{"a": 1, "b": 2, "c": 3},
			remove: map[string]any{"x": 10, "y": 20},
			want:   map[string]any{"a": 1, "b": 2, "c": 3},
		},
		{
			name:   "non-associative slice that doesn't match",
			source: []any{1, 2, 3},
			remove: []any{1, 2},
			want:   []any{1, 2, 3},
		},
		{
			name:   "non-associative slice that doesn't match but same length",
			source: []any{1, 2, 3},
			remove: []any{4, 5, 6},
			want:   []any{1, 2, 3},
		},
		{
			name:   "non-associative slice that matches",
			source: []any{1, 2, 3},
			remove: []any{1, 2, 3},
			want:   []any(nil),
		},
		{
			name:   "scalar that doesn't match",
			source: 42,
			remove: "remove me",
			want:   42,
		},
		{
			name:   "scalar that matches",
			source: 42,
			remove: "remove me",
			want:   42,
		},
		{
			name:   "complex map",
			source: map[string]any{"a": 1, "b": 2, "c": 3, "d": map[string]any{"e": 4, "f": []any{1, 2, 3}}},
			remove: map[string]any{"b": 2, "d": map[string]any{"e": 4, "f": []any{1, 2, 3, 4, 5, 6}}, "g": 7},
			want:   map[string]any{"a": 1, "c": 3, "d": map[string]any{"f": []any{1, 2, 3}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := delta(tt.source, tt.remove)
			assert.Equal(t, tt.want, got)
		})
	}
}
