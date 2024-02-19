package merger

import (
	"testing"

	"github.com/RRethy/krepe/deepishcopy"
	"github.com/stretchr/testify/assert"
)

type twoWayMergeTest[T any] struct {
	name     string
	local    T
	upstream T
	want     any
}

func runTwoWayMergeTests[T any](t *testing.T, mergeFunc func(T, T) any, tests []twoWayMergeTest[T]) {
	t.Helper()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			localCopy := deepishcopy.Copy(test.local)
			upstreamCopy := deepishcopy.Copy(test.upstream)
			got := mergeFunc(test.local, test.upstream)
			assert.Equal(t, test.want, got)
			assert.Equal(t, localCopy, test.local, "local should not be mutated")
			assert.Equal(t, upstreamCopy, test.upstream, "upstream should not be mutated")
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
		{
			name:     "nil upstream",
			local:    "foo",
			upstream: nil,
			want:     "foo",
		},
		{
			name:     "nil upstream but local map",
			local:    map[string]any{"foo": "bar"},
			upstream: nil,
			want:     map[string]any{"foo": "bar"},
		},
	})
}

func TestTwoWayMergeMap(t *testing.T) {
	runTwoWayMergeTests(t, twoWayMergeMap, []twoWayMergeTest[map[string]any]{
		{
			name:     "disjoint maps",
			local:    map[string]any{"a": 1, "b": 2, "c": 3},
			upstream: map[string]any{"d": 1, "e": 2, "f": 3},
			want:     map[string]any{"a": 1, "b": 2, "c": 3, "d": 1, "e": 2, "f": 3},
		},
		{
			name:     "overlapping maps",
			local:    map[string]any{"a": 1, "b": 2, "c": 3},
			upstream: map[string]any{"b": 3, "c": 4, "d": 5},
			want:     map[string]any{"a": 1, "b": 3, "c": 4, "d": 5},
		},
		{
			name:     "equal maps",
			local:    map[string]any{"a": 1, "b": 2, "c": 3},
			upstream: map[string]any{"a": 1, "b": 2, "c": 3},
			want:     map[string]any{"a": 1, "b": 2, "c": 3},
		},
		{
			name: "complex maps",
			local: map[string]any{
				"a": map[string]any{
					"a": 0,
					"b": 1,
					"c": 2,
					"d": []any{1, 2},
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
					"d": []any{3, 4},
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
					"d": []any{3, 4},
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
		{
			name:     "nil upstream",
			local:    []any{"foo", "bar"},
			upstream: nil,
			want:     []any{"foo", "bar"},
		},
		{
			name:     "nil local",
			local:    nil,
			upstream: []any{"foo", "bar"},
			want:     []any{"foo", "bar"},
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
		{
			name:     "nil upstream",
			local:    "foo",
			upstream: nil,
			want:     "foo",
		},
		{
			name:     "nil local",
			local:    nil,
			upstream: "foo",
			want:     "foo",
		},
	})
}

func TestTwoWayMergePtrStruct(t *testing.T) {
	runTwoWayMergeTests(t, twoWayMergePtrStruct, []twoWayMergeTest[any]{
		{
			name: "nil upstream",
			local: &testStruct{
				A: "foo",
				B: 1,
			},
			upstream: nil,
			want: &testStruct{
				A: "foo",
				B: 1,
			},
		},
		{
			name:     "nil upstream and local",
			local:    nil,
			upstream: nil,
			want:     nil,
		},
		{
			name:  "nil local",
			local: nil,
			upstream: &testStruct{
				A: "foo",
				B: 1,
			},
			want: &testStruct{
				A: "foo",
				B: 1,
			},
		},
		{
			name:     "ptr to empty structs",
			local:    &testStruct{},
			upstream: &testStruct{},
			want:     &testStruct{},
		},
		{
			name: "ptr to non-empty structs",
			local: &testStruct{
				A: "foo",
				B: 1,
			},
			upstream: &testStruct{
				A: "bar",
				B: 2,
			},
			want: &testStruct{
				A: "bar",
				B: 2,
			},
		},
		{
			name: "ptr to non-empty structs with nil upstream",
			local: &testStruct{
				A: "foo",
				B: 1,
				C: testStruct2{
					D: "baz",
					E: 3,
				},
			},
			upstream: nil,
			want: &testStruct{
				A: "foo",
				B: 1,
				C: testStruct2{
					D: "baz",
					E: 3,
				},
			},
		},
		{
			name:  "ptr to non-empty structs with nil local",
			local: nil,
			upstream: &testStruct{
				A: "foo",
				B: 1,
				C: testStruct2{
					D: "baz",
					E: 3,
				},
			},
			want: &testStruct{
				A: "foo",
				B: 1,
				C: testStruct2{
					D: "baz",
					E: 3,
				},
			},
		},
	})
}

func TestTwoWayMergeStruct(t *testing.T) {
	tests := []struct {
		name     string
		local    any
		upstream any
		want     any
	}{
		{
			name:     "empty structs",
			local:    testStruct{},
			upstream: testStruct{},
			want:     testStruct{},
		},
		{
			name: "non-empty structs",
			local: testStruct{
				A: "foo",
				B: 1,
				C: testStruct2{D: "baz", E: 3},
				F: &testStruct{
					A: "qux",
					B: 4,
					C: testStruct2{D: "corge", E: 5},
					F: &testStruct{A: "grault", B: 6},
				},
			},
			upstream: testStruct{
				A: "bar",
				B: 2,
				C: testStruct2{D: "quux", E: 4},
				F: &testStruct{
					A: "garply",
					B: 5,
					C: testStruct2{D: "waldo", E: 6},
					F: nil,
				},
			},
			want: testStruct{
				A: "bar",
				B: 2,
				C: testStruct2{D: "quux", E: 4},
				F: &testStruct{
					A: "garply",
					B: 5,
					C: testStruct2{D: "waldo", E: 6},
					F: &testStruct{A: "grault", B: 6},
				},
			},
		},
		{
			name:     "non-empty local with empty upstream",
			local:    testStruct{A: "foo", B: 1, C: testStruct2{}},
			upstream: testStruct{},
			want:     testStruct{A: "", B: 0, C: testStruct2{}},
		},
		{
			name:     "empty local with non-empty upstream",
			local:    testStruct{},
			upstream: testStruct{A: "foo", B: 1, C: testStruct2{D: "baz", E: 3}},
			want:     testStruct{A: "foo", B: 1, C: testStruct2{D: "baz", E: 3}},
		},
		{
			name:     "different types",
			local:    testStruct{A: "foo", B: 1},
			upstream: testStruct2{D: "baz", E: 3},
			want:     testStruct2{D: "baz", E: 3},
		},
		{
			name:     "equal structs",
			local:    testStruct{A: "foo", B: 1},
			upstream: testStruct{A: "foo", B: 1},
			want:     testStruct{A: "foo", B: 1},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, twoWayMergeStruct(test.local, test.upstream))
		})
	}
}

func TestTwoWayMergeStructSlice(t *testing.T) {
	tests := []struct {
		name     string
		local    []any
		upstream []any
		want     []any
	}{
		{
			name:     "empty slices",
			local:    []any{},
			upstream: []any{},
			want:     []any{},
		},
		{
			name:     "nil slices",
			local:    nil,
			upstream: nil,
			want:     nil,
		},
		{
			name:     "non-empty slices",
			local:    []any{testStruct{A: "foo", B: 1}, testStruct{A: "bar", B: 2}},
			upstream: []any{testStruct{A: "baz", B: 3}, testStruct{A: "qux", B: 4}},
			want:     []any{testStruct{A: "baz", B: 3}, testStruct{A: "qux", B: 4}},
		},
		{
			name:     "non-empty local with nil upstream",
			local:    []any{testStruct{A: "foo", B: 1}, testStruct{A: "bar", B: 2}},
			upstream: nil,
			want:     []any{testStruct{A: "foo", B: 1}, testStruct{A: "bar", B: 2}},
		},
		{
			name:     "nil local with non-empty upstream",
			local:    nil,
			upstream: []any{testStruct{A: "foo", B: 1}, testStruct{A: "bar", B: 2}},
			want:     []any{testStruct{A: "foo", B: 1}, testStruct{A: "bar", B: 2}},
		},
		{
			name:     "different types",
			local:    []any{testStruct{A: "foo", B: 1}, testStruct{A: "bar", B: 2}},
			upstream: []any{testStruct2{D: "baz", E: 3}, testStruct2{D: "qux", E: 4}},
			want:     []any{testStruct2{D: "baz", E: 3}, testStruct2{D: "qux", E: 4}},
		},
		{
			name:     "with associative keys",
			local:    []any{testStruct3{Name: "foo", Value: 1}, testStruct3{Name: "bar", Value: 2}, testStruct3{Name: "baz", Value: 3}},
			upstream: []any{testStruct3{Name: "foo", Value: 5}, testStruct3{Name: "grault", Value: 6}},
			want:     []any{testStruct3{Name: "foo", Value: 5}, testStruct3{Name: "bar", Value: 2}, testStruct3{Name: "baz", Value: 3}, testStruct3{Name: "grault", Value: 6}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, twoWayMergeStructSlice(test.local, test.upstream))
		})
	}
}

func TestTwoWayMergePtrStructSlice(t *testing.T) {
	tests := []struct {
		name     string
		local    []any
		upstream []any
		want     []any
	}{
		{
			name:     "empty slices",
			local:    []any{},
			upstream: []any{},
			want:     []any{},
		},
		{
			name:     "nil slices",
			local:    nil,
			upstream: nil,
			want:     nil,
		},
		{
			name:     "non-empty slices",
			local:    []any{&testStruct{A: "foo", B: 1}, &testStruct{A: "bar", B: 2}},
			upstream: []any{&testStruct{A: "baz", B: 3}, &testStruct{A: "qux", B: 4}},
			want:     []any{&testStruct{A: "baz", B: 3}, &testStruct{A: "qux", B: 4}},
		},
		{
			name:     "non-empty local with nil upstream",
			local:    []any{&testStruct{A: "foo", B: 1}, &testStruct{A: "bar", B: 2}},
			upstream: nil,
			want:     []any{&testStruct{A: "foo", B: 1}, &testStruct{A: "bar", B: 2}},
		},
		{
			name:     "nil local with non-empty upstream",
			local:    nil,
			upstream: []any{&testStruct{A: "foo", B: 1}, &testStruct{A: "bar", B: 2}},
			want:     []any{&testStruct{A: "foo", B: 1}, &testStruct{A: "bar", B: 2}},
		},
		{
			name:     "different types",
			local:    []any{&testStruct{A: "foo", B: 1}, &testStruct{A: "bar", B: 2}},
			upstream: []any{&testStruct2{D: "baz", E: 3}, &testStruct2{D: "qux", E: 4}},
			want:     []any{&testStruct2{D: "baz", E: 3}, &testStruct2{D: "qux", E: 4}},
		},
		{
			name:     "with associative keys",
			local:    []any{&testStruct3{Name: "foo", Value: 1}, &testStruct3{Name: "bar", Value: 2}, &testStruct3{Name: "baz", Value: 3}},
			upstream: []any{&testStruct3{Name: "foo", Value: 5}, &testStruct3{Name: "grault", Value: 6}},
			want:     []any{&testStruct3{Name: "foo", Value: 5}, &testStruct3{Name: "bar", Value: 2}, &testStruct3{Name: "baz", Value: 3}, &testStruct3{Name: "grault", Value: 6}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, twoWayMergePtrStructSlice(test.local, test.upstream))
		})
	}
}
