package krmmerge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type twoWayMergeTest[T mergeable] struct {
	name     string
	local    T
	upstream T
	want     any
}

func runTwoWayMergeTests[T mergeable](t *testing.T, mergeFunc func(T, T) any, tests []twoWayMergeTest[T]) {
	t.Helper()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := mergeFunc(test.local, test.upstream)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestTwoWayMerge(t *testing.T) {
	runTwoWayMergeTests(t, twoWayMerge, []twoWayMergeTest[any]{
		{
			name:     "scalars",
			local:    "foo",
			upstream: "bar",
			want:     "bar",
		},
		{
			name:     "maps",
			local:    map[string]any{"foo": "bar"},
			upstream: map[string]any{"baz": "qux"},
			want:     map[string]any{"foo": "bar", "baz": "qux"},
		},
		{
			name:     "associative slices",
			local:    []any{map[string]any{"name": "foo", "bar": "baz"}},
			upstream: []any{map[string]any{"name": "qux", "quux": "corge"}},
			want: []any{
				map[string]any{"name": "foo", "bar": "baz"},
				map[string]any{"name": "qux", "quux": "corge"},
			},
		},
		{
			name:     "non-associative slices",
			local:    []any{"foo", "bar"},
			upstream: []any{"baz", "qux"},
			want:     []any{"baz", "qux"},
		},
		{
			name:     "mismatching types",
			local:    "foo",
			upstream: []any{"baz", "qux"},
			want:     []any{"baz", "qux"},
		},
	})
}

func TestTwoWayMergeMap(t *testing.T) {
	runTwoWayMergeTests(t, twoWayMergeMap, []twoWayMergeTest[map[string]any]{
		{
			name: "disjoint maps",
			local: map[string]any{
				"a": 1,
				"b": 2,
				"c": 3,
			},
			upstream: map[string]any{
				"d": 1,
				"e": 2,
				"f": 3,
			},
			want: map[string]any{
				"a": 1,
				"b": 2,
				"c": 3,
				"d": 1,
				"e": 2,
				"f": 3,
			},
		},
		{
			name: "overlapping maps",
			local: map[string]any{
				"a": 1,
				"b": 2,
				"c": 3,
			},
			upstream: map[string]any{
				"b": 3,
				"c": 4,
				"d": 5,
			},
			want: map[string]any{
				"a": 1,
				"b": 3,
				"c": 4,
				"d": 5,
			},
		},
		{
			name: "equal maps",
			local: map[string]any{
				"a": 1,
				"b": 2,
				"c": 3,
			},
			upstream: map[string]any{
				"a": 1,
				"b": 2,
				"c": 3,
			},
			want: map[string]any{
				"a": 1,
				"b": 2,
				"c": 3,
			},
		},
		{
			name: "complex maps",
			local: map[string]any{
				"a": map[string]any{
					"a": 0,
					"b": 1,
					"c": 2,
					"d": []any{
						1,
						2,
					},
					"e": []any{
						map[string]any{"name": "foo", "bar": "baz"},
						map[string]any{"name": "qux", "quux": "corge"},
						map[string]any{"name": "quux", "quux": "corge"},
					},
				},
			},
			upstream: map[string]any{
				"a": map[string]any{
					"a": 3,
					"b": 1,
					"d": []any{
						3,
						4,
					},
					"e": []any{
						map[string]any{"name": "foo", "bar": "baz2"},
						map[string]any{"name": "qux", "quux": "corge"},
					},
				},
			},
			want: map[string]any{
				"a": map[string]any{
					"a": 3,
					"b": 1,
					"c": 2,
					"d": []any{
						3,
						4,
					},
					"e": []any{
						map[string]any{"name": "foo", "bar": "baz2"},
						map[string]any{"name": "qux", "quux": "corge"},
						map[string]any{"name": "quux", "quux": "corge"},
					},
				},
			},
		},
	})
}

func TestTwoWayMergeSlice(t *testing.T) {
	runTwoWayMergeTests(t, twoWayMergeSlice, []twoWayMergeTest[[]any]{
		{
			name: "no local slice associative key",
			local: []any{
				map[string]any{"a": 1, "b": 2, "c": 3},
				map[string]any{"d": 1, "e": 2, "f": 3},
			},
			upstream: []any{
				map[string]any{"name": "foo", "a": 1, "b": 2, "c": 3},
				map[string]any{"name": "bar", "d": 1, "e": 2, "f": 3},
			},
			want: []any{
				map[string]any{"name": "foo", "a": 1, "b": 2, "c": 3},
				map[string]any{"name": "bar", "d": 1, "e": 2, "f": 3},
			},
		},
		{
			name: "no upstream slice associative key",
			local: []any{
				map[string]any{"name": "foo", "a": 1, "b": 2, "c": 3},
				map[string]any{"name": "bar", "d": 1, "e": 2, "f": 3},
			},
			upstream: []any{
				map[string]any{"a": 1, "b": 2, "c": 3},
				map[string]any{"d": 1, "e": 2, "f": 3},
			},
			want: []any{
				map[string]any{"a": 1, "b": 2, "c": 3},
				map[string]any{"d": 1, "e": 2, "f": 3},
			},
		},
		{
			name: "both slices associative",
			local: []any{
				map[string]any{"name": "foo", "a": 1, "b": 2, "c": 3},
				map[string]any{"name": "bar", "d": 1, "e": 2, "f": 3},
			},
			upstream: []any{
				map[string]any{"name": "baz", "a": 1, "b": 2, "c": 3},
				map[string]any{"name": "qux", "d": 1, "e": 2, "f": 3},
			},
			want: []any{
				map[string]any{"name": "foo", "a": 1, "b": 2, "c": 3},
				map[string]any{"name": "bar", "d": 1, "e": 2, "f": 3},
				map[string]any{"name": "baz", "a": 1, "b": 2, "c": 3},
				map[string]any{"name": "qux", "d": 1, "e": 2, "f": 3},
			},
		},
		{
			name: "neither slice associative",
			local: []any{
				map[string]any{"a": "b"},
				map[string]any{"c": "d"},
			},
			upstream: []any{
				map[string]any{"e": "f"},
				map[string]any{"g": "h"},
			},
			want: []any{
				map[string]any{"e": "f"},
				map[string]any{"g": "h"},
			},
		},
	})
}

func TestTwoWayMergeSliceAssociative(t *testing.T) {
	runTwoWayMergeTests(t, twoWayMergeSliceAssociative, []twoWayMergeTest[[]any]{
		{
			name: "no local slice associative key",
			local: []any{
				map[string]any{"a": 1, "b": 2, "c": 3},
				map[string]any{"d": 1, "e": 2, "f": 3},
			},
			upstream: []any{
				map[string]any{"name": "foo", "a": 1, "b": 2, "c": 3},
				map[string]any{"name": "bar", "d": 1, "e": 2, "f": 3},
			},
			want: []any{
				map[string]any{"name": "foo", "a": 1, "b": 2, "c": 3},
				map[string]any{"name": "bar", "d": 1, "e": 2, "f": 3},
			},
		},
		{
			name: "no upstream slice associative key",
			local: []any{
				map[string]any{"name": "foo", "a": 1, "b": 2, "c": 3},
				map[string]any{"name": "bar", "d": 1, "e": 2, "f": 3},
			},
			upstream: []any{
				map[string]any{"a": 1, "b": 2, "c": 3},
				map[string]any{"d": 1, "e": 2, "f": 3},
			},
			want: []any{
				map[string]any{"a": 1, "b": 2, "c": 3},
				map[string]any{"d": 1, "e": 2, "f": 3},
			},
		},
		{
			name: "equal slices",
			local: []any{
				map[string]any{"name": "foo", "a": 1, "b": 2, "c": 3},
				map[string]any{"name": "bar", "d": 1, "e": 2, "f": 3},
			},
			upstream: []any{
				map[string]any{"name": "foo", "a": 1, "b": 2, "c": 3},
				map[string]any{"name": "bar", "d": 1, "e": 2, "f": 3},
			},
			want: []any{
				map[string]any{"name": "foo", "a": 1, "b": 2, "c": 3},
				map[string]any{"name": "bar", "d": 1, "e": 2, "f": 3},
			},
		},
		{
			name: "no overlap",
			local: []any{
				map[string]any{"name": "foo", "a": 1, "b": 2, "c": 3},
				map[string]any{"name": "bar", "d": 1, "e": 2, "f": 3},
			},
			upstream: []any{
				map[string]any{"name": "baz", "a": 1, "b": 2, "c": 3},
				map[string]any{"name": "qux", "d": 1, "e": 2, "f": 3},
			},
			want: []any{
				map[string]any{"name": "foo", "a": 1, "b": 2, "c": 3},
				map[string]any{"name": "bar", "d": 1, "e": 2, "f": 3},
				map[string]any{"name": "baz", "a": 1, "b": 2, "c": 3},
				map[string]any{"name": "qux", "d": 1, "e": 2, "f": 3},
			},
		},
		{
			name: "overlap",
			local: []any{
				map[string]any{"name": "foo", "a": 1, "b": 2, "c": 3},
				map[string]any{"name": "bar", "d": 1, "e": 2, "f": 3},
			},
			upstream: []any{
				map[string]any{"name": "foo", "a": 2, "b": 3, "c": 4},
				map[string]any{"name": "qux", "d": 1, "e": 2, "f": 3},
			},
			want: []any{
				map[string]any{"name": "foo", "a": 2, "b": 3, "c": 4},
				map[string]any{"name": "bar", "d": 1, "e": 2, "f": 3},
				map[string]any{"name": "qux", "d": 1, "e": 2, "f": 3},
			},
		},
	})
}

func TestTwoWayMergeSliceNonAssociative(t *testing.T) {
	runTwoWayMergeTests(t, twoWayMergeSliceNonAssociative, []twoWayMergeTest[[]any]{
		{
			name:     "matching slices",
			local:    []any{"foo", "bar"},
			upstream: []any{"foo", "bar"},
			want:     []any{"foo", "bar"},
		},
		{
			name:     "mismatching slices",
			local:    []any{"foo", "bar"},
			upstream: []any{"baz", "qux"},
			want:     []any{"baz", "qux"},
		},
	})
}

func TestTwoWayMergeScalar(t *testing.T) {
	runTwoWayMergeTests(t, twoWayMergeScalar, []twoWayMergeTest[any]{
		{
			name:     "matching scalars",
			local:    "foo",
			upstream: "foo",
			want:     "foo",
		},
		{
			name:     "mismatching scalars",
			local:    "foo",
			upstream: "bar",
			want:     "bar",
		},
	})
}
