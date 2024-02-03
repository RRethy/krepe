package merger

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			obj: struct {
				A string
				B int
			}{},
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
	type s struct {
		A string
		B int
		C struct {
			D string
			E int
		}
	}
	m := map[string]any{
		"A": "a",
		"B": 1,
		"C": struct {
			D string
			E int
		}{D: "d", E: 1},
	}
	want := s{A: "a", B: 1, C: struct {
		D string
		E int
	}{D: "d", E: 1}}
	assert.Equal(t, want, mapStringAnyToStruct(m, reflect.TypeOf(want)))
}

func TestMapStringAnyToPtrStruct(t *testing.T) {
	type s struct {
		A string
		B int
		C struct {
			D string
			E int
		}
	}
	m := map[string]any{
		"A": "a",
		"B": 1,
		"C": struct {
			D string
			E int
		}{D: "d", E: 1},
	}
	want := &s{A: "a", B: 1, C: struct {
		D string
		E int
	}{D: "d", E: 1}}
	assert.Equal(t, want, mapStringAnyToPtrStruct(m, reflect.TypeOf(want)))
	assert.Equal(t, want, mapStringAnyToPtrStruct(m, reflect.TypeOf(*want)))
}

func TestIsUniformStructSlice(t *testing.T) {
	type s struct {
		A string
		B int
	}

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
			slice: []any{s{A: "a", B: 1}, s{A: "b", B: 2}},
			want:  true,
		},
		{
			name:  "slice of different type structs",
			slice: []any{s{A: "a", B: 1}, struct{ A string }{A: "b"}},
			want:  false,
		},
		{
			name:  "slice of non-structs",
			slice: []any{1, 2},
			want:  false,
		},
		{
			name:  "slice of ptr to same type structs",
			slice: []any{&s{A: "a", B: 1}, &s{A: "b", B: 2}},
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
	type s struct {
		A string
		B int
	}

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
			slice: []any{&s{A: "a", B: 1}, &s{A: "b", B: 2}},
			want:  true,
		},
		{
			name:  "slice of different type struct ptr",
			slice: []any{&s{A: "a", B: 1}, &struct{ A string }{A: "b"}},
			want:  false,
		},
		{
			name:  "slice of non-struct ptrs",
			slice: []any{new(int), new(int)},
			want:  false,
		},
		{
			name:  "slice of struct",
			slice: []any{s{A: "a", B: 1}, s{A: "b", B: 2}},
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
	type s2 struct {
		D string
		E int
	}
	type s struct {
		A string
		B int
		C s2
	}

	tests := []struct {
		name  string
		slice []any
		want  []map[string]any
	}{
		{
			name:  "empty slice",
			slice: []any{},
			want:  []map[string]any{},
		},
		{
			name:  "slice of non-structs",
			slice: []any{1, 2},
			want:  []map[string]any{{}, {}},
		},
		{
			name:  "slice of structs",
			slice: []any{s{A: "a", B: 1}, s{A: "b", B: 2}},
			want:  []map[string]any{{"A": "a", "B": 1, "C": s2{}}, {"A": "b", "B": 2, "C": s2{}}},
		},
		{
			name:  "slice of mixed structs and non-structs",
			slice: []any{s{A: "a", B: 1}, 2},
			want:  []map[string]any{{"A": "a", "B": 1, "C": s2{}}, {}},
		},
		{
			name:  "slice of complex structs",
			slice: []any{s{A: "a", B: 1, C: s2{D: "d", E: 1}}, s{A: "b", B: 2, C: s2{D: "e", E: 2}}},
			want: []map[string]any{
				{"A": "a", "B": 1, "C": s2{D: "d", E: 1}},
				{"A": "b", "B": 2, "C": s2{D: "e", E: 2}},
			},
		},
		{
			name:  "slice of ptr to structs",
			slice: []any{&s{A: "a", B: 1}, &s{A: "b", B: 2}},
			want:  []map[string]any{{}, {}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, sliceStructToSliceMap(tt.slice))
		})
	}
}

func TestSlicePtrStructToSliceMap(t *testing.T) {
	type s2 struct {
		D string
		E int
	}
	type s struct {
		A string
		B int
		C s2
	}

	tests := []struct {
		name  string
		slice []any
		want  []map[string]any
	}{
		{
			name:  "empty slice",
			slice: []any{},
			want:  []map[string]any{},
		},
		{
			name:  "slice of non-struct ptrs",
			slice: []any{new(int), new(int)},
			want:  []map[string]any{{}, {}},
		},
		{
			name:  "slice of struct ptrs",
			slice: []any{&s{A: "a", B: 1}, &s{A: "b", B: 2}},
			want:  []map[string]any{{"A": "a", "B": 1, "C": s2{}}, {"A": "b", "B": 2, "C": s2{}}},
		},
		{
			name:  "slice of mixed struct ptrs and non-struct ptrs",
			slice: []any{&s{A: "a", B: 1}, new(int)},
			want:  []map[string]any{{"A": "a", "B": 1, "C": s2{}}, {}},
		},
		{
			name:  "slice of complex struct ptrs",
			slice: []any{&s{A: "a", B: 1, C: s2{D: "d", E: 1}}, &s{A: "b", B: 2, C: s2{D: "e", E: 2}}},
			want: []map[string]any{
				{"A": "a", "B": 1, "C": s2{D: "d", E: 1}},
				{"A": "b", "B": 2, "C": s2{D: "e", E: 2}},
			},
		},
		{
			name:  "slice of structs",
			slice: []any{s{A: "a", B: 1}, s{A: "b", B: 2}},
			want:  []map[string]any{{}, {}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, slicePtrStructToSliceMap(tt.slice))
		})
	}
}
