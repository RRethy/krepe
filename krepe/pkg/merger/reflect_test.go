package merger

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct2 struct {
	D string
	E int
}

type testStruct struct {
	A string
	B int
	C testStruct2
}

func TestStructToMap(t *testing.T) {
	tests := []struct {
		name string
		obj  any
		want map[string]any
	}{
		{
			name: "struct with a bunch of scalar fields",
			obj: struct {
				A string
				B int
				C bool
			}{A: "a", B: 1, C: true},
			want: map[string]any{"A": "a", "B": 1, "C": true},
		},
		{
			name: "non-struct",
			obj:  1,
			want: map[string]any{},
		},
		{
			name: "struct with a bunch of non-scalar fields",
			obj: struct {
				A []string
				B map[string]string
				C struct {
					D string
					E int
				}
			}{
				A: []string{"a", "b"},
				B: map[string]string{"a": "b"},
				C: struct {
					D string
					E int
				}{D: "d", E: 1},
			},
			want: map[string]any{
				"A": []string{"a", "b"},
				"B": map[string]string{"a": "b"},
				"C": struct {
					D string
					E int
				}{D: "d", E: 1},
			},
		},
		{
			name: "struct with empty fields",
			obj: struct {
				A string
				B int
			}{},
			want: map[string]any{"A": "", "B": 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, structToMap(tt.obj))
		})
	}
}

func TestPtrStructToMap(t *testing.T) {
	tests := []struct {
		name string
		obj  any
		want map[string]any
	}{
		{
			name: "pointer to struct with a bunch of scalar fields",
			obj: &struct {
				A string
				B int
				C bool
			}{A: "a", B: 1, C: true},
			want: map[string]any{"A": "a", "B": 1, "C": true},
		},
		{
			name: "pointer to non-struct",
			obj:  new(int),
			want: map[string]any{},
		},
		{
			name: "pointer to struct with a bunch of non-scalar fields",
			obj: &struct {
				A []string
				B map[string]string
				C struct {
					D string
					E int
				}
			}{
				A: []string{"a", "b"},
				B: map[string]string{"a": "b"},
				C: struct {
					D string
					E int
				}{D: "d", E: 1},
			},
			want: map[string]any{
				"A": []string{"a", "b"},
				"B": map[string]string{"a": "b"},
				"C": struct {
					D string
					E int
				}{D: "d", E: 1},
			},
		},
		{
			name: "non-pointer",
			obj:  testStruct{},
			want: map[string]any{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ptrStructToMap(tt.obj))
		})
	}
}

func TestMapStringAnyToStruct(t *testing.T) {
	m := map[string]any{
		"A": "a",
		"B": 1,
		"C": struct {
			D string
			E int
		}{D: "d", E: 1},
	}
	want := testStruct{A: "a", B: 1, C: testStruct2{D: "d", E: 1}}
	assert.Equal(t, want, mapStringAnyToStruct(m, reflect.TypeOf(want)))
}

func TestMapStringAnyToPtrStruct(t *testing.T) {
	m := map[string]any{
		"A": "a",
		"B": 1,
		"C": struct {
			D string
			E int
		}{D: "d", E: 1},
	}
	want := &testStruct{A: "a", B: 1, C: testStruct2{D: "d", E: 1}}
	assert.Equal(t, want, mapStringAnyToPtrStruct(m, reflect.TypeOf(want)))
	assert.Equal(t, want, mapStringAnyToPtrStruct(m, reflect.TypeOf(*want)))
}

func TestIsUniformStructSlice(t *testing.T) {
	tests := []struct {
		name  string
		slice []any
		want  bool
	}{
		{
			name:  "empty slice",
			slice: []any{},
			want:  false,
		},
		{
			name:  "slice of same type structs",
			slice: []any{testStruct{A: "a", B: 1}, testStruct{A: "b", B: 2}},
			want:  true,
		},
		{
			name:  "slice of different type structs",
			slice: []any{testStruct{A: "a", B: 1}, struct{ A string }{A: "b"}},
			want:  false,
		},
		{
			name:  "slice of non-structs",
			slice: []any{1, 2},
			want:  false,
		},
		{
			name:  "slice of ptr to same type structs",
			slice: []any{&testStruct{A: "a", B: 1}, &testStruct{A: "b", B: 2}},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, isUniformStructSlice(tt.slice))
		})
	}
}

func TestIsUniformPtrStructSlice(t *testing.T) {
	tests := []struct {
		name  string
		slice []any
		want  bool
	}{
		{
			name:  "empty slice",
			slice: []any{},
			want:  false,
		},
		{
			name:  "slice of same type struct ptr",
			slice: []any{&testStruct{A: "a", B: 1}, &testStruct{A: "b", B: 2}},
			want:  true,
		},
		{
			name:  "slice of different type struct ptr",
			slice: []any{&testStruct{A: "a", B: 1}, &struct{ A string }{A: "b"}},
			want:  false,
		},
		{
			name:  "slice of non-struct ptrs",
			slice: []any{new(int), new(int)},
			want:  false,
		},
		{
			name:  "slice of struct",
			slice: []any{testStruct{A: "a", B: 1}, testStruct{A: "b", B: 2}},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, isUniformPtrStructSlice(tt.slice))
		})
	}
}

func TestSliceStructToSliceMap(t *testing.T) {
	tests := []struct {
		name  string
		slice []any
		want  []any
	}{
		{
			name:  "empty slice",
			slice: []any{},
			want:  []any{},
		},
		{
			name:  "slice of non-structs",
			slice: []any{1, 2},
			want:  []any{map[string]any{}, map[string]any{}},
		},
		{
			name:  "slice of structs",
			slice: []any{testStruct{A: "a", B: 1}, testStruct{A: "b", B: 2}},
			want:  []any{map[string]any{"A": "a", "B": 1, "C": testStruct2{}}, map[string]any{"A": "b", "B": 2, "C": testStruct2{}}},
		},
		{
			name:  "slice of mixed structs and non-structs",
			slice: []any{testStruct{A: "a", B: 1}, 2},
			want:  []any{map[string]any{"A": "a", "B": 1, "C": testStruct2{}}, map[string]any{}},
		},
		{
			name:  "slice of complex structs",
			slice: []any{testStruct{A: "a", B: 1, C: testStruct2{D: "d", E: 1}}, testStruct{A: "b", B: 2, C: testStruct2{D: "e", E: 2}}},
			want: []any{
				map[string]any{"A": "a", "B": 1, "C": testStruct2{D: "d", E: 1}},
				map[string]any{"A": "b", "B": 2, "C": testStruct2{D: "e", E: 2}},
			},
		},
		{
			name:  "slice of ptr to structs",
			slice: []any{&testStruct{A: "a", B: 1}, &testStruct{A: "b", B: 2}},
			want:  []any{map[string]any{}, map[string]any{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, sliceStructToSliceMap(tt.slice))
		})
	}
}

func TestSlicePtrStructToSliceMap(t *testing.T) {
	tests := []struct {
		name  string
		slice []any
		want  []any
	}{
		{
			name:  "empty slice",
			slice: []any{},
			want:  []any{},
		},
		{
			name:  "slice of non-struct ptrs",
			slice: []any{new(int), new(int)},
			want:  []any{map[string]any{}, map[string]any{}},
		},
		{
			name:  "slice of struct ptrs",
			slice: []any{&testStruct{A: "a", B: 1}, &testStruct{A: "b", B: 2}},
			want:  []any{map[string]any{"A": "a", "B": 1, "C": testStruct2{}}, map[string]any{"A": "b", "B": 2, "C": testStruct2{}}},
		},
		{
			name:  "slice of mixed struct ptrs and non-struct ptrs",
			slice: []any{&testStruct{A: "a", B: 1}, new(int)},
			want:  []any{map[string]any{"A": "a", "B": 1, "C": testStruct2{}}, map[string]any{}},
		},
		{
			name:  "slice of complex struct ptrs",
			slice: []any{&testStruct{A: "a", B: 1, C: testStruct2{D: "d", E: 1}}, &testStruct{A: "b", B: 2, C: testStruct2{D: "e", E: 2}}},
			want: []any{
				map[string]any{"A": "a", "B": 1, "C": testStruct2{D: "d", E: 1}},
				map[string]any{"A": "b", "B": 2, "C": testStruct2{D: "e", E: 2}},
			},
		},
		{
			name:  "slice of structs",
			slice: []any{testStruct{A: "a", B: 1}, testStruct{A: "b", B: 2}},
			want:  []any{map[string]any{}, map[string]any{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, slicePtrStructToSliceMap(tt.slice))
		})
	}
}

func TestIsUniformStructSlices(t *testing.T) {
	tests := []struct {
		name   string
		slices [][]any
		want   bool
	}{
		{
			name:   "empty slices",
			slices: [][]any{},
			want:   false,
		},
		{
			name:   "slices of same type structs",
			slices: [][]any{{testStruct{A: "a", B: 1}, testStruct{A: "b", B: 2}}, {testStruct{A: "c", B: 3}, testStruct{A: "d", B: 4}}},
			want:   true,
		},
		{
			name:   "slices of different type structs",
			slices: [][]any{{testStruct{A: "a", B: 1}, testStruct{A: "b", B: 2}}, {struct{ A string }{A: "c"}, struct{ A string }{A: "d"}}},
			want:   false,
		},
		{
			name:   "slices of non-structs",
			slices: [][]any{{1, 2}, {3, 4}},
			want:   false,
		},
		{
			name:   "slices of mixed structs and non-structs",
			slices: [][]any{{testStruct{A: "a", B: 1}, 2}, {testStruct{A: "c", B: 3}, 4}},
			want:   false,
		},
		{
			name:   "slices of ptr to same type structs",
			slices: [][]any{{&testStruct{A: "a", B: 1}, &testStruct{A: "b", B: 2}}, {&testStruct{A: "c", B: 3}, &testStruct{A: "d", B: 4}}},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, isUniformStructSlices(tt.slices...))
		})
	}
}

func TestIsUniformPtrStructSlices(t *testing.T) {
	tests := []struct {
		name   string
		slices [][]any
		want   bool
	}{
		{
			name:   "empty slices",
			slices: [][]any{},
			want:   false,
		},
		{
			name:   "slices of same type struct ptrs",
			slices: [][]any{{&testStruct{A: "a", B: 1}, &testStruct{A: "b", B: 2}}, {&testStruct{A: "c", B: 3}, &testStruct{A: "d", B: 4}}},
			want:   true,
		},
		{
			name:   "slices of different type struct ptrs",
			slices: [][]any{{&testStruct{A: "a", B: 1}, &testStruct{A: "b", B: 2}}, {&struct{ A string }{A: "c"}, &struct{ A string }{A: "d"}}},
			want:   false,
		},
		{
			name:   "slices of non-struct ptrs",
			slices: [][]any{{new(int), new(int)}, {new(int), new(int)}},
			want:   false,
		},
		{
			name:   "slices of mixed struct ptrs and non-struct ptrs",
			slices: [][]any{{&testStruct{A: "a", B: 1}, new(int)}, {&testStruct{A: "c", B: 3}, new(int)}},
			want:   false,
		},
		{
			name:   "slices of struct",
			slices: [][]any{{testStruct{A: "a", B: 1}, testStruct{A: "b", B: 2}}, {testStruct{A: "c", B: 3}, testStruct{A: "d", B: 4}}},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, isUniformPtrStructSlices(tt.slices...))
		})
	}
}

func TestSliceMapToSliceStruct(t *testing.T) {
	tests := []struct {
		name       string
		slice      []any
		structType reflect.Type
		want       []any
	}{
		{
			name:  "empty slice",
			slice: []any{},
			want:  []any{},
		},
		{
			name:       "slice of empty maps",
			slice:      []any{map[string]any{}, map[string]any{}},
			structType: reflect.TypeOf(testStruct{}),
			want:       []any{testStruct{}, testStruct{}},
		},
		{
			name:       "slice of non-empty maps",
			slice:      []any{map[string]any{"A": "a", "B": 1, "C": testStruct2{}}, map[string]any{"A": "b", "B": 2, "C": testStruct2{}}},
			structType: reflect.TypeOf(testStruct{}),
			want:       []any{testStruct{A: "a", B: 1, C: testStruct2{}}, testStruct{A: "b", B: 2, C: testStruct2{}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func TestSliceMapToSlicePtrStruct(t *testing.T) {
	tests := []struct {
		name       string
		slice      []any
		structType reflect.Type
		want       []any
	}{
		{
			name:       "empty slice",
			slice:      []any{},
			structType: reflect.TypeOf(&testStruct{}),
			want:       []any{},
		},
		{
			name:       "slice of empty maps with ptr struct type",
			slice:      []any{map[string]any{}, map[string]any{}},
			structType: reflect.TypeOf(&testStruct{}),
			want:       []any{&testStruct{}, &testStruct{}},
		},
		{
			name:       "slice of non-empty maps with ptr struct type",
			slice:      []any{map[string]any{"A": "a", "B": 1, "C": testStruct2{}}, map[string]any{"A": "b", "B": 2, "C": testStruct2{}}},
			structType: reflect.TypeOf(&testStruct{}),
			want:       []any{&testStruct{A: "a", B: 1, C: testStruct2{}}, &testStruct{A: "b", B: 2, C: testStruct2{}}},
		},
		{
			name:       "slice of non-empty maps with non-ptr struct type",
			slice:      []any{map[string]any{"A": "a", "B": 1, "C": testStruct2{}}, map[string]any{"A": "b", "B": 2, "C": testStruct2{}}},
			structType: reflect.TypeOf(testStruct{}),
			want:       []any{testStruct{A: "a", B: 1, C: testStruct2{}}, testStruct{A: "b", B: 2, C: testStruct2{}}},
		},
		{
			name:       "slice of empty maps with non-ptr struct type",
			slice:      []any{map[string]any{}, map[string]any{}},
			structType: reflect.TypeOf(testStruct{}),
			want:       []any{testStruct{}, testStruct{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}
