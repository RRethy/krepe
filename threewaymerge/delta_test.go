package threewaymerge

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.design/x/reflect"
)

type deltaTest[T mergeable] struct {
	name   string
	source T
	remove T
	want   any
}

func runDeltaTests[T mergeable](t *testing.T, deltaFunc func(T, T) any, tests []deltaTest[T]) {
	t.Helper()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sourceCopy := reflect.DeepCopy(test.source)
			removeCopy := reflect.DeepCopy(test.remove)
			got := deltaFunc(test.source, test.remove)
			assert.Equal(t, test.want, got)
			assert.Equal(t, sourceCopy, test.source, "source should not be mutated")
			assert.Equal(t, removeCopy, test.remove, "remove should not be mutated")
		})
	}
}

func TestDelta(t *testing.T) {
	runDeltaTests(t, delta, []deltaTest[any]{
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
			want:   nil,
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
	})
}

func TestDeltaMap(t *testing.T) {
	runDeltaTests(t, deltaMap, []deltaTest[map[string]any]{
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
	})
}

func TestDeltaSlice(t *testing.T) {
	runDeltaTests(t, deltaSlice, []deltaTest[[]any]{
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
			want:   nil,
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
	})
}

func TestDeltaSliceNonAssociative(t *testing.T) {
	runDeltaTests(t, deltaSliceNonAssociative, []deltaTest[[]any]{
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
			want:   nil,
		},
	})
}

func TestDeltaSliceAssociative(t *testing.T) {
	runDeltaTests(t, deltaSliceAssociative, []deltaTest[[]any]{
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
	})
}

func TestDeltaScalar(t *testing.T) {
	runDeltaTests(t, deltaScalar, []deltaTest[any]{
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
	})
}
