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
}

func TestIsUniformPtrStructSlice(t *testing.T) {
}

func TestSliceStructToSliceMap(t *testing.T) {
}

func TestSlicePtrStructToSliceMap(t *testing.T) {
}
