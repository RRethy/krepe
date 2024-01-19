package threewaymerge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsAssociativeSlice(t *testing.T) {
	tests := []struct {
		name  string
		slice []any
		want  bool
	}{
		{
			name:  "nil slice",
			slice: nil,
			want:  true,
		},
		{
			name:  "empty slice",
			slice: []any{},
			want:  true,
		},
		{
			name:  "non-associative slice",
			slice: []any{1, 2, 3},
			want:  false,
		},
		{
			name:  "mixed slice values",
			slice: []any{1, map[string]any{"a": 1, "b": 2, "c": 3}, 3},
			want:  false,
		},
		{
			name: "associative slice",
			slice: []any{
				map[string]any{"a": 1, "b": 2, "c": 3},
				map[string]any{"d": 1, "e": 2, "f": 3},
			},
			want: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := isAssociativeSlice(test.slice)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestGetAssociativeKeys(t *testing.T) {
	tests := []struct {
		name  string
		slice []any
		want  []string
	}{
		{
			name:  "nil slice",
			slice: nil,
			want:  nil,
		},
		{
			name:  "empty slice",
			slice: []any{},
			want:  nil,
		},
		{
			name:  "non-associative slice",
			slice: []any{1, 2, 3},
			want:  nil,
		},
		{
			name:  "associative slice with no associative key",
			slice: []any{map[string]any{"a": 1, "b": 2, "c": 3}},
			want:  nil,
		},
		{
			name: "associative slice with associative key",
			slice: []any{
				map[string]any{"name": "foo", "a": 2, "b": 3, "c": 4},
				map[string]any{"name": "bar", "a": 2, "b": 3, "c": 4},
				map[string]any{"name": "baz", "a": 2, "b": 3, "c": 4},
			},
			want: []string{"name"},
		},
		{
			name: "associative slice with multiple associative keys with duplicate values",
			slice: []any{
				map[string]any{"name": "foo", "a": 2, "b": 3, "c": 4},
				map[string]any{"name": "bar", "a": 2, "b": 3, "c": 4, "d": 5},
				map[string]any{"name": "bar", "a": 2, "b": 3, "c": 4, "d": 5},
			},
			want: nil,
		},
		{
			name: "associative slice with multiple associative keys",
			slice: []any{
				map[string]any{"name": "foo", "type": "foo", "a": 2, "b": 3, "c": 4},
				map[string]any{"name": "bar", "type": "bar", "a": 2, "b": 3, "c": 4, "d": 5},
			},
			want: []string{"type", "name"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := getAssociativeKeys(test.slice)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestHasAssociativeKey(t *testing.T) {
	tests := []struct {
		name  string
		slice []any
		key   string
		want  bool
	}{
		{
			name:  "nil slice",
			slice: nil,
			key:   "name",
			want:  true,
		},
		{
			name:  "empty slice",
			slice: []any{},
			key:   "name",
			want:  true,
		},
		{
			name:  "non-associative slice",
			slice: []any{1, 2, 3},
			key:   "name",
			want:  false,
		},
		{
			name:  "associative slice with no associative key",
			slice: []any{map[string]any{"a": 1, "b": 2, "c": 3}},
			key:   "name",
			want:  false,
		},
		{
			name: "associative slice with associative key",
			slice: []any{
				map[string]any{"name": "foo", "a": 2, "b": 3, "c": 4},
				map[string]any{"name": "bar", "a": 2, "b": 3, "c": 4},
				map[string]any{"name": "baz", "a": 2, "b": 3, "c": 4},
			},
			key:  "name",
			want: true,
		},
		{
			name: "associative slice with multiple associative keys with duplicate values",
			slice: []any{
				map[string]any{"name": "foo", "a": 2, "b": 3, "c": 4},
				map[string]any{"name": "bar", "a": 2, "b": 3, "c": 4, "d": 5},
				map[string]any{"name": "bar", "a": 2, "b": 3, "c": 4, "d": 5},
			},
			key:  "name",
			want: false,
		},
		{
			name: "associative slice with multiple associative keys",
			slice: []any{
				map[string]any{"name": "foo", "type": "foo", "a": 2, "b": 3, "c": 4},
				map[string]any{"name": "bar", "type": "bar", "a": 2, "b": 3, "c": 4, "d": 5},
			},
			key:  "type",
			want: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := hasAssociativeKey(test.slice, test.key)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestGetCommonAssociativeKey(t *testing.T) {
	tests := []struct {
		name  string
		keys1 []string
		keys2 []string
		keys3 []string
		want  string
	}{
		{
			name:  "empty keys1 and keys2",
			keys1: []string{},
			keys2: []string{},
			want:  "",
		},
		{
			name:  "empty keys1",
			keys1: []string{},
			keys2: []string{"name"},
			want:  "",
		},
		{
			name:  "empty keys2",
			keys1: []string{"name"},
			keys2: []string{},
			want:  "",
		},
		{
			name:  "no common keys",
			keys1: []string{"name"},
			keys2: []string{"type"},
			want:  "",
		},
		{
			name:  "one common key",
			keys1: []string{"name"},
			keys2: []string{"name"},
			want:  "name",
		},
		{
			name:  "multiple common keys",
			keys1: []string{"name", "type"},
			keys2: []string{"type", "name"},
			want:  "name",
		},
		{
			name:  "multiple common keys with duplicate values",
			keys1: []string{"name", "type"},
			keys2: []string{"type", "name", "name"},
			want:  "name",
		},
		{
			name:  "three slices",
			keys1: []string{"name", "type"},
			keys2: []string{"type", "mountPath", "name"},
			keys3: []string{"type", "name"},
			want:  "name",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var got string
			if len(test.keys3) > 0 {
				got = getCommonAssociativeKey(test.keys1, test.keys2, test.keys3)
			} else {
				got = getCommonAssociativeKey(test.keys1, test.keys2)
			}
			assert.Equal(t, test.want, got)
		})
	}
}

func TestGetCommonAssociativeKeys(t *testing.T) {
	tests := []struct {
		name  string
		keyss [][]string
		want  []string
	}{
		{
			name:  "empty keyss",
			keyss: [][]string{},
			want:  nil,
		},
		{
			name:  "one keyss",
			keyss: [][]string{{"name", "type"}},
			want:  []string{"name", "type"},
		},
		{
			name:  "two keyss",
			keyss: [][]string{{"name", "type"}, {"type", "name", "mountPath"}},
			want:  []string{"name", "type"},
		},
		{
			name:  "three keyss",
			keyss: [][]string{{"name", "type"}, {"type", "name", "mountPath"}, {"devicePath", "type", "name"}},
			want:  []string{"name", "type"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := getCommonAssociativeKeys(test.keyss)
			assert.Equal(t, test.want, got)
		})
	}
}
