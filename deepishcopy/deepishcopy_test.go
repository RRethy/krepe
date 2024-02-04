package deepishcopy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Foo    string
	Bar    int
	Nested *testStruct
}

func TestCopy(t *testing.T) {
	tests := []struct {
		name          string
		src           any
		want          any
		modifications func(any)
	}{
		{
			name: "nil",
			src:  nil,
			want: nil,
		},
		{
			name: "bool",
			src:  true,
			want: true,
		},
		{
			name: "int",
			src:  1,
			want: 1,
		},
		{
			name: "int8",
			src:  int8(2),
			want: int8(2),
		},
		{
			name: "int16",
			src:  int16(3),
			want: int16(3),
		},
		{
			name: "int32",
			src:  int32(4),
			want: int32(4),
		},
		{
			name: "int64",
			src:  int64(5),
			want: int64(5),
		},
		{
			name: "uint",
			src:  uint(6),
			want: uint(6),
		},
		{
			name: "uint8",
			src:  uint8(6),
			want: uint8(6),
		},
		{
			name: "uint16",
			src:  uint16(7),
			want: uint16(7),
		},
		{
			name: "uint32",
			src:  uint32(8),
			want: uint32(8),
		},
		{
			name: "uint64",
			src:  uint64(9),
			want: uint64(9),
		},
		{
			name: "float32",
			src:  float32(1.1),
			want: float32(1.1),
		},
		{
			name: "float64",
			src:  float64(2.1),
			want: float64(2.1),
		},
		{
			name: "complex64",
			src:  complex64(3.1),
			want: complex64(3.1),
		},
		{
			name: "complex128",
			src:  complex128(4.1),
			want: complex128(4.1),
		},
		{
			name: "string",
			src:  "foo",
			want: "foo",
		},
		{
			name: "slice",
			src:  []any{1, 2, 3},
			want: []any{1, 2, 3},
			modifications: func(obj any) {
				obj.([]any)[0] = 4
			},
		},
		{
			name: "map",
			src: map[string]any{
				"foo": "bar",
			},
			want: map[string]any{
				"foo": "bar",
			},
			modifications: func(obj any) {
				obj.(map[string]any)["foo"] = "baz"
			},
		},
		{
			name: "nested slice",
			src:  []any{1, 2, []any{3, 4}},
			want: []any{1, 2, []any{3, 4}},
			modifications: func(obj any) {
				obj.([]any)[2].([]any)[0] = 5
			},
		},
		{
			name: "nested map",
			src: map[string]any{
				"foo": map[string]any{
					"bar": "baz",
				},
			},
			want: map[string]any{
				"foo": map[string]any{
					"bar": "baz",
				},
			},
			modifications: func(obj any) {
				obj.(map[string]any)["foo"].(map[string]any)["bar"] = "qux"
			},
		},
		{
			name: "nested slice and map",
			src: []any{
				1,
				map[string]any{
					"foo": []any{
						2,
						map[string]any{
							"bar": 3,
						},
					},
				},
			},
			want: []any{
				1,
				map[string]any{
					"foo": []any{
						2,
						map[string]any{
							"bar": 3,
						},
					},
				},
			},
			modifications: func(obj any) {
				obj.([]any)[1].(map[string]any)["foo"].([]any)[1].(map[string]any)["bar"] = 4
			},
		},
		{
			name: "nil slice",
			src:  nil,
			want: nil,
		},
		{
			name: "nil map",
			src:  nil,
			want: nil,
		},
		{
			name: "empty slice",
			src:  []any{},
			want: []any{},
		},
		{
			name: "empty map",
			src:  map[string]any{},
			want: map[string]any{},
			modifications: func(obj any) {
				obj.(map[string]any)["foo"] = "bar"
			},
		},
		{
			name: "struct ptr",
			src: &testStruct{
				Foo: "foo",
				Bar: 1,
				Nested: &testStruct{
					Foo: "nested",
					Bar: 2,
					Nested: &testStruct{
						Foo: "nested2",
						Bar: 3,
					},
				},
			},
			want: &testStruct{
				Foo: "foo",
				Bar: 1,
				Nested: &testStruct{
					Foo: "nested",
					Bar: 2,
					Nested: &testStruct{
						Foo: "nested2",
						Bar: 3,
					},
				},
			},
			modifications: func(obj any) {
				obj.(*testStruct).Nested.Nested.Foo = "qux2"
			},
		},
		{
			name: "struct",
			src: testStruct{
				Foo: "foo",
				Bar: 1,
				Nested: &testStruct{
					Foo: "nested",
					Bar: 2,
					Nested: &testStruct{
						Foo: "nested2",
						Bar: 3,
					},
				},
			},
			want: testStruct{
				Foo: "foo",
				Bar: 1,
				Nested: &testStruct{
					Foo: "nested",
					Bar: 2,
					Nested: &testStruct{
						Foo: "nested2",
						Bar: 3,
					},
				},
			},
			modifications: func(obj any) {
				obj.(testStruct).Nested.Nested.Foo = "qux2"
			},
		},
		{
			name: "struct ptr with nil field",
			src: &testStruct{
				Foo: "foo",
				Bar: 1,
			},
			want: &testStruct{
				Foo: "foo",
				Bar: 1,
			},
			modifications: func(obj any) {
				obj.(*testStruct).Foo = "bar"
			},
		},
		{
			name: "struct with nil field",
			src: testStruct{
				Foo: "foo",
				Bar: 1,
			},
			want: testStruct{
				Foo: "foo",
				Bar: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Copy(tt.src)
			assert.Equal(t, tt.want, got)
			if tt.modifications != nil {
				tt.modifications(got)
				assert.NotEqual(t, tt.src, got)
			}
		})
	}
}
