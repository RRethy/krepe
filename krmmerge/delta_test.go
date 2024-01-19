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
			name: "associative slice with no remove slice associative key",
			source: []any{
				map[string]any{"name": "foo", "a": 2, "b": 3, "c": 4},
			},
			remove: []any{
				map[string]any{"a": 1, "b": 2, "c": 3},
				map[string]any{"name": "foo", "a": 4, "b": 5, "c": 6},
			},
			want: []any{
				map[string]any{"name": "foo", "a": 2, "b": 3, "c": 4},
			},
		},
		{
			name: "associative slice with no source slice associative key",
			source: []any{
				map[string]any{"a": 1, "b": 2, "c": 3},
				map[string]any{"name": "foo", "a": 4, "b": 5, "c": 6},
			},
			remove: []any{
				map[string]any{"name": "foo", "a": 2, "b": 3, "c": 4},
			},
			want: []any{
				map[string]any{"a": 1, "b": 2, "c": 3},
				map[string]any{"name": "foo", "a": 4, "b": 5, "c": 6},
			},
		},
		{
			name: "associative slice with overlap",
			source: []any{
				map[string]any{"name": "foo", "a": 1, "b": 2, "c": 3},
				map[string]any{"name": "bar", "a": 4, "b": 5, "c": 6},
				map[string]any{"name": "baz", "a": 7, "b": 8, "c": 9},
			},
			remove: []any{
				map[string]any{"name": "foo", "a": 2, "b": 3, "c": 4},
				map[string]any{"name": "foo2", "a": 2, "b": 3, "c": 4},
			},
			want: []any{
				map[string]any{"name": "bar", "a": 4, "b": 5, "c": 6},
				map[string]any{"name": "baz", "a": 7, "b": 8, "c": 9},
			},
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
		{
			name: "mismatched map type",
			source: map[string]any{
				"a": 1,
				"b": 2,
			},
			remove: 5,
			want: map[string]any{
				"a": 1,
				"b": 2,
			},
		},
		{
			name: "mismatched slice type",
			source: []any{
				1,
				2,
			},
			remove: "foo",
			want: []any{
				1,
				2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := delta(tt.source, tt.remove)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDeltaMap(t *testing.T) {
	tests := []struct {
		name   string
		source map[string]any
		remove map[string]any
		want   map[string]any
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
			name:   "map with overlapping keys",
			source: map[string]any{"a": 1, "b": 2, "c": 3},
			remove: map[string]any{"b": 3, "c": 4},
			want:   map[string]any{"a": 1, "b": 2, "c": 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := deltaMap(tt.source, tt.remove)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDeltaSlice(t *testing.T) {
	tests := []struct {
		name   string
		source []any
		remove []any
		want   []any
	}{
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
			name: "associative slice with no remove slice associative key",
			source: []any{
				map[string]any{"name": "foo", "a": 2, "b": 3, "c": 4},
			},
			remove: []any{
				map[string]any{"a": 1, "b": 2, "c": 3},
				map[string]any{"name": "foo", "a": 4, "b": 5, "c": 6},
			},
			want: []any{
				map[string]any{"name": "foo", "a": 2, "b": 3, "c": 4},
			},
		},
		{
			name: "associative slice with no source slice associative key",
			source: []any{
				map[string]any{"a": 1, "b": 2, "c": 3},
				map[string]any{"name": "foo", "a": 4, "b": 5, "c": 6},
			},
			remove: []any{
				map[string]any{"name": "foo", "a": 2, "b": 3, "c": 4},
			},
			want: []any{
				map[string]any{"a": 1, "b": 2, "c": 3},
				map[string]any{"name": "foo", "a": 4, "b": 5, "c": 6},
			},
		},
		{
			name: "associative slice with overlap",
			source: []any{
				map[string]any{"name": "foo", "a": 1, "b": 2, "c": 3},
				map[string]any{"name": "bar", "a": 4, "b": 5, "c": 6},
				map[string]any{"name": "baz", "a": 7, "b": 8, "c": 9},
			},
			remove: []any{
				map[string]any{"name": "foo", "a": 2, "b": 3, "c": 4},
				map[string]any{"name": "foo2", "a": 2, "b": 3, "c": 4},
			},
			want: []any{
				map[string]any{"name": "bar", "a": 4, "b": 5, "c": 6},
				map[string]any{"name": "baz", "a": 7, "b": 8, "c": 9},
			},
		},
		{
			name: "associative slice with no overlap",
			source: []any{
				map[string]any{"name": "foo", "a": 1, "b": 2, "c": 3},
				map[string]any{"name": "bar", "a": 4, "b": 5, "c": 6},
			},
			remove: []any{
				map[string]any{"name": "baz", "a": 2, "b": 3, "c": 4},
				map[string]any{"name": "baz2", "a": 2, "b": 3, "c": 4},
			},
			want: []any{
				map[string]any{"name": "foo", "a": 1, "b": 2, "c": 3},
				map[string]any{"name": "bar", "a": 4, "b": 5, "c": 6},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := deltaSlice(tt.source, tt.remove)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDeltaSliceNonAssociative(t *testing.T) {
	tests := []struct {
		name   string
		source []any
		remove []any
		want   []any
	}{
		{
			name:   "doesn't match",
			source: []any{1, 2, 3},
			remove: []any{1, 2},
			want:   []any{1, 2, 3},
		},
		{
			name:   "doesn't match but same length",
			source: []any{1, 2, 3},
			remove: []any{4, 5, 6},
			want:   []any{1, 2, 3},
		},
		{
			name:   "matches",
			source: []any{1, 2, 3},
			remove: []any{1, 2, 3},
			want:   []any(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := deltaSliceNonAssociative(tt.source, tt.remove)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDeltaSliceAssociative(t *testing.T) {
	tests := []struct {
		name   string
		source []any
		remove []any
		want   []any
	}{
		{
			name: "no remove slice associative key",
			source: []any{
				map[string]any{"name": "foo", "a": 2, "b": 3, "c": 4},
			},
			remove: []any{
				map[string]any{"a": 1, "b": 2, "c": 3},
				map[string]any{"name": "foo", "a": 4, "b": 5, "c": 6},
			},
			want: []any{
				map[string]any{"name": "foo", "a": 2, "b": 3, "c": 4},
			},
		},
		{
			name: "no source slice associative key",
			source: []any{
				map[string]any{"a": 1, "b": 2, "c": 3},
				map[string]any{"name": "foo", "a": 4, "b": 5, "c": 6},
			},
			remove: []any{
				map[string]any{"name": "foo", "a": 2, "b": 3, "c": 4},
			},
			want: []any{
				map[string]any{"a": 1, "b": 2, "c": 3},
				map[string]any{"name": "foo", "a": 4, "b": 5, "c": 6},
			},
		},
		{
			name: "overlap",
			source: []any{
				map[string]any{"name": "foo", "a": 1, "b": 2, "c": 3},
				map[string]any{"name": "bar", "a": 4, "b": 5, "c": 6},
				map[string]any{"name": "baz", "a": 7, "b": 8, "c": 9},
			},
			remove: []any{
				map[string]any{"name": "foo", "a": 2, "b": 3, "c": 4},
				map[string]any{"name": "foo2", "a": 2, "b": 3, "c": 4},
			},
			want: []any{
				map[string]any{"name": "bar", "a": 4, "b": 5, "c": 6},
				map[string]any{"name": "baz", "a": 7, "b": 8, "c": 9},
			},
		},
		{
			name: "no overlap",
			source: []any{
				map[string]any{"name": "foo", "a": 1, "b": 2, "c": 3},
				map[string]any{"name": "bar", "a": 4, "b": 5, "c": 6},
			},
			remove: []any{
				map[string]any{"name": "baz", "a": 2, "b": 3, "c": 4},
				map[string]any{"name": "baz2", "a": 2, "b": 3, "c": 4},
			},
			want: []any{
				map[string]any{"name": "foo", "a": 1, "b": 2, "c": 3},
				map[string]any{"name": "bar", "a": 4, "b": 5, "c": 6},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := deltaSliceAssociative(tt.source, tt.remove)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDeltaScalar(t *testing.T) {
	tests := []struct {
		name   string
		source any
		remove any
		want   any
	}{
		{
			name:   "match",
			source: 5,
			remove: 5,
			want:   nil,
		},
		{
			name:   "no match",
			source: 5,
			remove: 10,
			want:   5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := deltaScalar(tt.source, tt.remove)
			assert.Equal(t, tt.want, got)
		})
	}
}
