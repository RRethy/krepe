package merger

import (
	"testing"

	"github.com/RRethy/krepe/deepishcopy"
	"github.com/stretchr/testify/assert"
)

type threeWayMergeTest[T any] struct {
	name     string
	origin   T
	local    T
	upstream T
	want     any
}

func runThreeWayMergeTests[T any](t *testing.T, threeWayMergeFunc func(T, T, T) any, tests []threeWayMergeTest[T]) {
	t.Helper()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			originCopy := deepishcopy.Copy(test.origin)
			localCopy := deepishcopy.Copy(test.local)
			upstreamCopy := deepishcopy.Copy(test.upstream)
			got := threeWayMergeFunc(test.origin, test.local, test.upstream)
			assert.Equal(t, test.want, got)
			assert.Equal(t, originCopy, test.origin, "origin should not be mutated")
			assert.Equal(t, localCopy, test.local, "local should not be mutated")
			assert.Equal(t, upstreamCopy, test.upstream, "upstream should not be mutated")
		})
	}
}

func TestThreeWayMergeAny(t *testing.T) {
	runThreeWayMergeTests(t, threeWayMerge, []threeWayMergeTest[any]{
		{
			name:     "all maps",
			origin:   map[string]any{"a": 1, "b": 2},
			local:    map[string]any{"a": 1, "b": 3},
			upstream: map[string]any{"a": 1, "b": 2},
			want:     map[string]any{"a": 1, "b": 3},
		},
		{
			name: "all slices",
			origin: []any{
				map[string]any{"name": "a", "value": 1},
				map[string]any{"name": "b", "value": 1},
			},
			local: []any{
				map[string]any{"name": "a", "value": 1},
			},
			upstream: []any{
				map[string]any{"name": "a", "value": 1},
				map[string]any{"name": "b", "value": 2},
			},
			want: []any{
				map[string]any{"name": "a", "value": 1},
				map[string]any{"name": "b", "value": 2},
			},
		},
		{
			name:     "all scalars",
			origin:   1,
			local:    2,
			upstream: 3,
			want:     3,
		},
		{
			name:     "local and upstream maps but origin scalar",
			origin:   1,
			local:    map[string]any{"a": 1, "b": 2},
			upstream: map[string]any{"a": 1, "b": 3},
			want:     map[string]any{"a": 1, "b": 3},
		},
		{
			name:     "local origin maps but upstream scalar",
			origin:   map[string]any{"a": 1, "b": 2, "d": 5},
			local:    map[string]any{"a": 1, "b": 3, "c": 4},
			upstream: 1,
			want:     1,
		},
		{
			name:     "origin and upstream slices but local scalar",
			origin:   []any{1, 2, 3},
			local:    1,
			upstream: []any{4, 5, 6},
			want:     []any{4, 5, 6},
		},
		{
			name:     "local and upstream slices but origin scalar",
			origin:   1,
			local:    []any{1, 2, 3},
			upstream: []any{4, 5, 6},
			want:     []any{4, 5, 6},
		},
		{
			name:     "local and upstream scalar but origin map",
			origin:   map[string]any{"a": 1, "b": 2},
			local:    1,
			upstream: 2,
			want:     2,
		},
		{
			name:     "origin and local scalar but upstream map",
			origin:   1,
			local:    2,
			upstream: map[string]any{"a": 1, "b": 2},
			want:     map[string]any{"a": 1, "b": 2},
		},
		{
			name:     "origin upstream nil but local map",
			origin:   nil,
			local:    map[string]any{"a": 1, "b": 2},
			upstream: nil,
			want:     map[string]any{"a": 1, "b": 2},
		},
		{
			name:     "upstream nil but origin local slice",
			origin:   []any{1, 2, 3},
			local:    []any{4, 5, 6},
			upstream: nil,
			want:     nil,
		},
	})
}

func TestThreeWayMergeMap(t *testing.T) {
	runThreeWayMergeTests(t, threeWayMergeMap, []threeWayMergeTest[map[string]any]{
		{
			name:     "matching keys but different values",
			origin:   map[string]any{"a": 1, "b": 2, "c": 3},
			local:    map[string]any{"a": 1, "b": 3, "c": 3},
			upstream: map[string]any{"a": 1, "b": 4, "c": 4},
			want:     map[string]any{"a": 1, "b": 4, "c": 4},
		},
		{
			name:     "matching keys and value but changed in upstream",
			origin:   map[string]any{"a": 1},
			local:    map[string]any{"a": 1},
			upstream: map[string]any{"a": 2},
			want:     map[string]any{"a": 2},
		},
		{
			name:     "matching keys and value",
			origin:   map[string]any{"a": 1},
			local:    map[string]any{"a": 1},
			upstream: map[string]any{"a": 1},
			want:     map[string]any{"a": 1},
		},
		{
			name:     "different keys",
			origin:   map[string]any{"a": 1, "b": 2, "c": 3},
			local:    map[string]any{"x": 10, "y": 20},
			upstream: map[string]any{"a": 1, "b": 2, "c": 3},
			want:     map[string]any{"x": 10, "y": 20},
		},
		{
			name:     "overlapping keys",
			origin:   map[string]any{"a": 1, "b": 2, "c": 3},
			local:    map[string]any{"a": 1, "b": 3, "d": 5},
			upstream: map[string]any{"a": 1, "b": 4, "c": 4},
			want:     map[string]any{"a": 1, "b": 4, "c": 4, "d": 5},
		},
		{
			name:     "key removed in upstream",
			origin:   map[string]any{"a": 1, "b": 2},
			local:    map[string]any{"a": 1, "b": 2},
			upstream: map[string]any{"b": 2},
			want:     map[string]any{"b": 2},
		},
		{
			name:     "key removed in upstream and modified local",
			origin:   map[string]any{"a": 1, "b": 2},
			local:    map[string]any{"a": 2, "b": 2},
			upstream: map[string]any{"b": 2},
			want:     map[string]any{"b": 2},
		},
		{
			name:     "key removed from both local and upstream",
			origin:   map[string]any{"a": 1, "b": 2},
			local:    map[string]any{"b": 2},
			upstream: map[string]any{"b": 2},
			want:     map[string]any{"b": 2},
		},
		{
			name:     "new key in upstream",
			origin:   map[string]any{"a": 1, "b": 2},
			local:    map[string]any{"a": 1, "b": 2},
			upstream: map[string]any{"a": 1, "b": 2, "c": 3},
			want:     map[string]any{"a": 1, "b": 2, "c": 3},
		},
		{
			name:     "complex key removed from upstream with no local modification",
			origin:   map[string]any{"foo": map[string]any{"a": 1, "b": 2}, "bar": 3},
			local:    map[string]any{"foo": map[string]any{"a": 1, "b": 2}, "bar": 3},
			upstream: map[string]any{"bar": 3},
			want:     map[string]any{"bar": 3},
		},
		{
			name:     "complex key removed from upstream with local modification",
			origin:   map[string]any{"foo": map[string]any{"a": 1, "b": 2}, "bar": 3},
			local:    map[string]any{"foo": map[string]any{"a": 1, "b": 3}, "bar": 3},
			upstream: map[string]any{"bar": 3},
			want:     map[string]any{"bar": 3},
		},
		{
			name:     "key removed from local but present in origin and upstream",
			origin:   map[string]any{"foo": 1, "bar": 2},
			local:    map[string]any{"foo": 1},
			upstream: map[string]any{"foo": 1, "bar": 2},
			want:     map[string]any{"foo": 1},
		},
		{
			name:     "key removed from local but present in origin and modified in upstream",
			origin:   map[string]any{"foo": 1, "bar": 2},
			local:    map[string]any{"foo": 1},
			upstream: map[string]any{"foo": 1, "bar": 3},
			want:     map[string]any{"foo": 1, "bar": 3},
		},
		{
			name:     "key added in upstream but same as present in local",
			origin:   map[string]any{"foo": 1},
			local:    map[string]any{"foo": 1, "bar": 3},
			upstream: map[string]any{"foo": 1, "bar": 3},
			want:     map[string]any{"foo": 1, "bar": 3},
		},
		{
			name:     "key added in upstream but different than present in local",
			origin:   map[string]any{"foo": 1},
			local:    map[string]any{"foo": 1, "bar": 4},
			upstream: map[string]any{"foo": 1, "bar": 3},
			want:     map[string]any{"foo": 1, "bar": 3},
		},
		{
			name:     "key present and same in all",
			origin:   map[string]any{"foo": 1, "bar": 2},
			local:    map[string]any{"foo": 1, "bar": 2},
			upstream: map[string]any{"foo": 1, "bar": 2},
			want:     map[string]any{"foo": 1, "bar": 2},
		},
		{
			name:     "key present in all but changed in upstream",
			origin:   map[string]any{"foo": 1, "bar": 2},
			local:    map[string]any{"foo": 1, "bar": 2},
			upstream: map[string]any{"foo": 1, "bar": 3},
			want:     map[string]any{"foo": 1, "bar": 3},
		},
		{
			name:     "key present in all but changed in local",
			origin:   map[string]any{"foo": 1, "bar": 2},
			local:    map[string]any{"foo": 1, "bar": 3},
			upstream: map[string]any{"foo": 1, "bar": 2},
			want:     map[string]any{"foo": 1, "bar": 3},
		},
		{
			name:     "key present in all but changed in local and upstream",
			origin:   map[string]any{"foo": 1, "bar": 2},
			local:    map[string]any{"foo": 1, "bar": 3},
			upstream: map[string]any{"foo": 1, "bar": 4},
			want:     map[string]any{"foo": 1, "bar": 4},
		},
		{
			name:     "key changed in local with nil",
			origin:   map[string]any{"foo": 1},
			local:    map[string]any{"foo": nil},
			upstream: map[string]any{"foo": 1},
			want:     nil,
		},
		{
			name:     "key changed in upstream with nil",
			origin:   map[string]any{"foo": 1},
			local:    map[string]any{"foo": 1},
			upstream: map[string]any{"foo": nil},
			want:     nil,
		},
		{
			name:     "key changed in local and upstream with nil",
			origin:   map[string]any{"foo": 1},
			local:    map[string]any{"foo": nil},
			upstream: map[string]any{"foo": nil},
			want:     nil,
		},
		{
			name:     "key changed in local with scalar and upstream with nil",
			origin:   map[string]any{"foo": 1},
			local:    map[string]any{"foo": 2},
			upstream: map[string]any{"foo": nil},
			want:     nil,
		},
	})
}

func TestThreeWayMergeSlice(t *testing.T) {
	runThreeWayMergeTests(t, threeWayMergeSlice, []threeWayMergeTest[[]any]{
		{
			name:     "non-associative slice",
			origin:   []any{1, 2, 3},
			local:    []any{3, 4},
			upstream: []any{6, 7},
			want:     []any{6, 7},
		},
		{
			name: "associative slice",
			origin: []any{
				map[string]any{"name": "a", "value": 1, "other": 2},
				map[string]any{"name": "b", "value": 2, "other": 3},
				map[string]any{"name": "c", "value": 3, "other": 4},
			},
			local: []any{
				map[string]any{"name": "a", "value": 1, "other": 2},
				map[string]any{"name": "b", "value": 2, "other": 3},
				map[string]any{"name": "d", "value": 4, "other": 5},
			},
			upstream: []any{
				map[string]any{"name": "a", "value": 1, "other": 3},
				map[string]any{"name": "c", "value": 3, "other": 4},
			},
			want: []any{
				map[string]any{"name": "a", "value": 1, "other": 3},
				map[string]any{"name": "d", "value": 4, "other": 5},
			},
		},
	})
}

func TestThreeWayMergeSliceAssociative(t *testing.T) {
	runThreeWayMergeTests(t, threeWayMergeSliceAssociative, []threeWayMergeTest[[]any]{
		{
			name: "elem in origin local and upstream with same value",
			origin: []any{
				map[string]any{"name": "foo", "a": 1},
				map[string]any{"name": "bar", "b": 2},
			},
			local: []any{
				map[string]any{"name": "foo", "a": 1},
				map[string]any{"name": "bar", "b": 2},
			},
			upstream: []any{
				map[string]any{"name": "foo", "a": 1},
				map[string]any{"name": "bar", "b": 2},
			},
			want: []any{
				map[string]any{"name": "foo", "a": 1},
				map[string]any{"name": "bar", "b": 2},
			},
		},
		{
			name: "elem in origin local but not upstream with same values",
			origin: []any{
				map[string]any{"name": "foo", "a": 1},
				map[string]any{"name": "bar", "b": 2},
			},
			local: []any{
				map[string]any{"name": "foo", "a": 1},
				map[string]any{"name": "bar", "b": 2},
			},
			upstream: []any{
				map[string]any{"name": "foo", "a": 1},
			},
			want: []any{
				map[string]any{"name": "foo", "a": 1},
			},
		},
		{
			name: "elem in origin local but not upstream with different values",
			origin: []any{
				map[string]any{"name": "foo", "a": 1},
				map[string]any{"name": "bar", "b": 2},
			},
			local: []any{
				map[string]any{"name": "foo", "a": 1},
				map[string]any{"name": "bar", "b": 3},
			},
			upstream: []any{
				map[string]any{"name": "foo", "a": 1},
			},
			want: []any{
				map[string]any{"name": "foo", "a": 1},
			},
		},
		{
			name: "elem in origin upstream but not local with same values",
			origin: []any{
				map[string]any{"name": "foo", "a": 1},
				map[string]any{"name": "bar", "b": 2},
			},
			local: []any{
				map[string]any{"name": "foo", "a": 1},
			},
			upstream: []any{
				map[string]any{"name": "foo", "a": 1},
				map[string]any{"name": "bar", "b": 2},
			},
			want: []any{
				map[string]any{"name": "foo", "a": 1},
			},
		},
		{
			name: "elem in origin upstream but not local with different values",
			origin: []any{
				map[string]any{"name": "foo", "a": 1},
				map[string]any{"name": "bar", "b": 2},
			},
			local: []any{
				map[string]any{"name": "foo", "a": 1},
			},
			upstream: []any{
				map[string]any{"name": "foo", "a": 1},
				map[string]any{"name": "bar", "b": 3},
			},
			want: []any{
				map[string]any{"name": "foo", "a": 1},
				map[string]any{"name": "bar", "b": 3},
			},
		},
		{
			name: "elem in local upstream with diff values but not origin",
			origin: []any{
				map[string]any{"name": "foo", "a": 1},
			},
			local: []any{
				map[string]any{"name": "foo", "a": 1},
				map[string]any{"name": "bar", "b": 2},
			},
			upstream: []any{
				map[string]any{"name": "foo", "a": 1},
				map[string]any{"name": "bar", "b": 3},
			},
			want: []any{
				map[string]any{"name": "foo", "a": 1},
				map[string]any{"name": "bar", "b": 3},
			},
		},
		{
			name: "elem in local upstream with same values but not origin",
			origin: []any{
				map[string]any{"name": "foo", "a": 1},
			},
			local: []any{
				map[string]any{"name": "foo", "a": 1},
				map[string]any{"name": "bar", "b": 2},
			},
			upstream: []any{
				map[string]any{"name": "foo", "a": 1},
				map[string]any{"name": "bar", "b": 2},
			},
			want: []any{
				map[string]any{"name": "foo", "a": 1},
				map[string]any{"name": "bar", "b": 2},
			},
		},
		{
			name: "elem in upstream only",
			origin: []any{
				map[string]any{"name": "foo", "a": 1},
			},
			local: []any{
				map[string]any{"name": "foo", "a": 1},
			},
			upstream: []any{
				map[string]any{"name": "bar", "b": 2},
				map[string]any{"name": "foo", "a": 1},
			},
			want: []any{
				map[string]any{"name": "foo", "a": 1},
				map[string]any{"name": "bar", "b": 2},
			},
		},
		{
			name: "elem in local only",
			origin: []any{
				map[string]any{"name": "foo", "a": 1},
			},
			local: []any{
				map[string]any{"name": "foo", "a": 1},
				map[string]any{"name": "bar", "b": 2},
			},
			upstream: []any{
				map[string]any{"name": "foo", "a": 1},
			},
			want: []any{
				map[string]any{"name": "foo", "a": 1},
				map[string]any{"name": "bar", "b": 2},
			},
		},
		{
			name: "non associative slice with no matching keys",
			origin: []any{
				map[string]any{"name": "foo", "a": 1},
			},
			local: []any{
				map[string]any{"name": "foo", "a": 1},
				2,
			},
			upstream: []any{
				map[string]any{"name": "foo", "a": 1},
			},
			want: []any{
				map[string]any{"name": "foo", "a": 1},
				2,
			},
		},
		{
			name: "associative but no matching keys",
			origin: []any{
				map[string]any{"name": "foo", "a": 1},
			},
			local: []any{
				map[string]any{"type": "foo", "a": 1},
			},
			upstream: []any{
				map[string]any{"name": "foo", "a": 1},
			},
			want: []any{
				map[string]any{"type": "foo", "a": 1},
			},
		},
	})
}

func TestThreeWayMergeSliceNonAssociative(t *testing.T) {
	runThreeWayMergeTests(t, threeWayMergeSliceNonAssociative, []threeWayMergeTest[[]any]{
		{
			name:     "origin and upstream match",
			origin:   []any{1, 2, 3},
			local:    []any{3, 4},
			upstream: []any{1, 2, 3},
			want:     []any{3, 4},
		},
		{
			name:     "origin and upstream don't match",
			origin:   []any{1, 2, 3},
			local:    []any{1, 2},
			upstream: []any{4, 5, 6},
			want:     []any{4, 5, 6},
		},
	})
}

func TestThreeWayMergeScalar(t *testing.T) {
	runThreeWayMergeTests[any](t, threeWayMergeScalar, []threeWayMergeTest[any]{
		{
			name:     "origin and upstream match",
			origin:   1,
			local:    2,
			upstream: 1,
			want:     2,
		},
		{
			name:     "origin and upstream don't match",
			origin:   1,
			local:    2,
			upstream: 3,
			want:     3,
		},
	})
}

func TestThreeWayMergePtrStruct(t *testing.T) {
	runThreeWayMergeTests[any](t, threeWayMergePtrStruct, []threeWayMergeTest[any]{
		{
			name:     "nil origin",
			origin:   nil,
			local:    &struct{ A int }{A: 1},
			upstream: &struct{ A int }{A: 2},
			want:     &struct{ A int }{A: 2},
		},
		{
			name:     "nil local with different upstream",
			origin:   &struct{ A int }{A: 1},
			local:    nil,
			upstream: &struct{ A int }{A: 2},
			want:     &struct{ A int }{A: 2},
		},
		{
			name:     "nil local with same upstream",
			origin:   &struct{ A int }{A: 1},
			local:    nil,
			upstream: &struct{ A int }{A: 1},
			want:     nil,
		},
		{
			name:     "nil upstream with different local",
			origin:   &struct{ A int }{A: 1},
			local:    &struct{ A int }{A: 2},
			upstream: nil,
			want:     nil,
		},
		{
			name:     "nil upstream with same local",
			origin:   &struct{ A int }{A: 1},
			local:    &struct{ A int }{A: 1},
			upstream: nil,
			want:     nil,
		},
		{
			name:     "nil origin and local",
			origin:   nil,
			local:    nil,
			upstream: &struct{ A int }{A: 1},
			want:     &struct{ A int }{A: 1},
		},
		{
			name:     "nil origin and upstream",
			origin:   nil,
			local:    &struct{ A int }{A: 1},
			upstream: nil,
			want:     &struct{ A int }{A: 1},
		},
		{
			name:     "nil local and upstream",
			origin:   &struct{ A int }{A: 1},
			local:    nil,
			upstream: nil,
			want:     nil,
		},
		{
			name:     "nil origin, local and upstream",
			origin:   nil,
			local:    nil,
			upstream: nil,
			want:     nil,
		},
		{
			name:     "ptr to empty structs",
			origin:   &struct{}{},
			local:    &struct{}{},
			upstream: &struct{}{},
			want:     &struct{}{},
		},
		{
			name:     "ptr to non-empty structs",
			origin:   &struct{ A int }{A: 1},
			local:    &struct{ A int }{A: 2},
			upstream: &struct{ A int }{A: 3},
			want:     &struct{ A int }{A: 3},
		},
		{
			name:   "ptr to different types of structs",
			origin: &struct{ A int }{A: 1},
			local:  &struct{ A int }{A: 2},
			upstream: &struct {
				A int
				C int
			}{A: 1, C: 3},
			want: &struct {
				A int
				C int
			}{A: 2, C: 3},
		},
	})
}

func TestThreeWayMergeStruct(t *testing.T) {
	runThreeWayMergeTests[any](t, threeWayMergeStruct, []threeWayMergeTest[any]{
		{
			name:     "empty structs",
			origin:   struct{}{},
			local:    struct{}{},
			upstream: struct{}{},
			want:     struct{}{},
		},
		{
			name:     "non-empty structs",
			origin:   struct{ A int }{A: 1},
			local:    struct{ A int }{A: 2},
			upstream: struct{ A int }{A: 3},
			want:     struct{ A int }{A: 3},
		},
		{
			name:   "different types of structs",
			origin: struct{ A int }{A: 1},
			local:  struct{ A int }{A: 2},
			upstream: struct {
				A int
				C int
			}{A: 1, C: 3},
			want: struct {
				A int
				C int
			}{A: 2, C: 3},
		},
	})
}

func TestThreeWayMergeStructSlice(t *testing.T) {
	runThreeWayMergeTests[[]any](t, threeWayMergeStructSlice, []threeWayMergeTest[[]any]{
		{
			name:     "empty slices",
			origin:   []any{},
			local:    []any{},
			upstream: []any{},
			want:     []any{},
		},
		{
			name: "non-empty associate slices",
			origin: []any{
				testStruct3{Name: "a", Value: 1},
				testStruct3{Name: "b", Value: 2},
			},
			local: []any{
				testStruct3{Name: "a", Value: 1},
				testStruct3{Name: "b", Value: 3},
				testStruct3{Name: "c", Value: 4},
			},
			upstream: []any{
				testStruct3{Name: "a", Value: 1},
				testStruct3{Name: "b", Value: 4},
				testStruct3{Name: "d", Value: 5},
			},
			want: []any{
				testStruct3{Name: "a", Value: 1},
				testStruct3{Name: "b", Value: 4},
				testStruct3{Name: "c", Value: 4},
				testStruct3{Name: "d", Value: 5},
			},
		},
		{
			name:     "non-empty non-associate slices",
			origin:   []any{testStruct{A: "1"}, testStruct{A: "2"}},
			local:    []any{testStruct{A: "1"}, testStruct{A: "3"}, testStruct{A: "4"}},
			upstream: []any{testStruct{A: "1"}, testStruct{A: "4"}, testStruct{A: "5"}},
			want:     []any{testStruct{A: "1"}, testStruct{A: "4"}, testStruct{A: "5"}},
		},
		{
			name:   "empty origin uses two way merge",
			origin: []any{},
			local: []any{
				testStruct3{Name: "a", Value: 1},
				testStruct3{Name: "b", Value: 3},
				testStruct3{Name: "c", Value: 4},
			},
			upstream: []any{
				testStruct3{Name: "a", Value: 1},
				testStruct3{Name: "b", Value: 4},
				testStruct3{Name: "d", Value: 5},
			},
			want: []any{
				testStruct3{Name: "a", Value: 1},
				testStruct3{Name: "b", Value: 4},
				testStruct3{Name: "c", Value: 4},
				testStruct3{Name: "d", Value: 5},
			},
		},
		{
			name: "different list types",
			origin: []any{
				testStruct{A: "1"},
				testStruct{A: "2"},
			},
			local: []any{
				testStruct2{D: "1"},
				testStruct2{D: "2"},
			},
			upstream: []any{
				testStruct3{Name: "a", Value: 1},
				testStruct3{Name: "b", Value: 2},
			},
			want: []any{
				testStruct3{Name: "a", Value: 1},
				testStruct3{Name: "b", Value: 2},
			},
		},
	})
}

func TestThreeWayMergePtrStructSlice(t *testing.T) {
	runThreeWayMergeTests[[]any](t, threeWayMergePtrStructSlice, []threeWayMergeTest[[]any]{
		{
			name:     "empty slices",
			origin:   []any{},
			local:    []any{},
			upstream: []any{},
			want:     []any{},
		},
		{
			name: "non-empty associate slices",
			origin: []any{
				&testStruct3{Name: "a", Value: 1},
				&testStruct3{Name: "b", Value: 2},
			},
			local: []any{
				&testStruct3{Name: "a", Value: 1},
				&testStruct3{Name: "b", Value: 3},
				&testStruct3{Name: "c", Value: 4},
			},
			upstream: []any{
				&testStruct3{Name: "a", Value: 1},
				&testStruct3{Name: "b", Value: 4},
				&testStruct3{Name: "d", Value: 5},
			},
			want: []any{
				&testStruct3{Name: "a", Value: 1},
				&testStruct3{Name: "b", Value: 4},
				&testStruct3{Name: "c", Value: 4},
				&testStruct3{Name: "d", Value: 5},
			},
		},
		{
			name:     "non-empty non-associate slices",
			origin:   []any{&testStruct{A: "1"}, &testStruct{A: "2"}},
			local:    []any{&testStruct{A: "1"}, &testStruct{A: "3"}, &testStruct{A: "4"}},
			upstream: []any{&testStruct{A: "1"}, &testStruct{A: "4"}, &testStruct{A: "5"}},
			want:     []any{&testStruct{A: "1"}, &testStruct{A: "4"}, &testStruct{A: "5"}},
		},
		{
			name:   "empty origin uses two way merge",
			origin: []any{},
			local: []any{
				&testStruct3{Name: "a", Value: 1},
				&testStruct3{Name: "b", Value: 3},
				&testStruct3{Name: "c", Value: 4},
			},
			upstream: []any{
				&testStruct3{Name: "a", Value: 1},
				&testStruct3{Name: "b", Value: 4},
				&testStruct3{Name: "d", Value: 5},
			},
			want: []any{
				&testStruct3{Name: "a", Value: 1},
				&testStruct3{Name: "b", Value: 4},
				&testStruct3{Name: "c", Value: 4},
				&testStruct3{Name: "d", Value: 5},
			},
		},
		{
			name: "different list types",
			origin: []any{
				&testStruct{A: "1"},
				&testStruct{A: "2"},
			},
			local: []any{
				&testStruct2{D: "1"},
				&testStruct2{D: "2"},
			},
			upstream: []any{
				&testStruct3{Name: "a", Value: 1},
				&testStruct3{Name: "b", Value: 2},
			},
			want: []any{
				&testStruct3{Name: "a", Value: 1},
				&testStruct3{Name: "b", Value: 2},
			},
		},
	})
}
