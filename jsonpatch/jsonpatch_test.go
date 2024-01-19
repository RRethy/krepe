package jsonpatch

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewJsonPatch(t *testing.T) {
	tests := []struct {
		name    string
		op      string
		path    string
		value   any
		want    *JsonPatch
		wantErr bool
	}{
		{
			name:    "Test with valid path",
			op:      "add",
			path:    "/foo/bar",
			value:   "quux",
			want:    &JsonPatch{op: "add", path: []string{"foo", "bar"}, value: "quux"},
			wantErr: false,
		},
		{
			name:    "Test with invalid path",
			op:      "add",
			path:    "foo/bar",
			value:   "quux",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Test with unknown op",
			op:      "foo",
			path:    "/foo/bar",
			value:   "quux",
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		got, err := NewJsonPatch(tt.op, tt.path, tt.value)
		assert.Equal(t, tt.want, got, tt.name)
		if tt.wantErr {
			assert.Error(t, err, tt.name)
		} else {
			assert.Equal(t, tt.want, got, tt.name)
		}
	}
}

func TestApply(t *testing.T) {
	tests := []struct {
		name    string
		obj     map[string]any
		patch   *JsonPatch
		want    map[string]any
		wantErr bool
	}{
		{
			name: "add without error",
			obj: map[string]any{
				"foo": "bar",
			},
			patch: &JsonPatch{
				op:    "add",
				path:  []string{"foo"},
				value: "qux",
			},
			want: map[string]any{
				"foo": "qux",
			},
			wantErr: false,
		},
		{
			name: "add with error",
			obj: map[string]any{
				"foo": "bar",
			},
			patch: &JsonPatch{
				op:    "add",
				path:  []string{"3", "foo"},
				value: "qux",
			},
			want: map[string]any{
				"foo": "bar",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.patch.Apply(tt.obj)
			assert.Equal(t, tt.want, got)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPatchToJsonPtrs(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    []string
		wantErr bool
	}{
		{
			name:    "Test with empty path",
			path:    "",
			want:    nil,
			wantErr: false,
		},
		{
			name:    "Test with valid path",
			path:    "/foo/bar",
			want:    []string{"foo", "bar"},
			wantErr: false,
		},
		{
			name:    "Test with path not starting with /",
			path:    "foo/bar",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Test with path containing ~1",
			path:    "/foo~1bar/baz",
			want:    []string{"foo/bar", "baz"},
			wantErr: false,
		},
		{
			name:    "Test with path containing ~0",
			path:    "/foo~0bar/baz",
			want:    []string{"foo~bar", "baz"},
			wantErr: false,
		},
		{
			name:    "Test with path containing ~1 and ~0",
			path:    "/foo~1bar~0baz/baz",
			want:    []string{"foo/bar~baz", "baz"},
			wantErr: false,
		},
		{
			name:    "Test with valid path",
			path:    "/foo/bar",
			want:    []string{"foo", "bar"},
			wantErr: false,
		},
		{
			name:    "Test with just /",
			path:    "/",
			want:    []string{""},
			wantErr: false,
		},
		{
			name:    "Test with path ending with /",
			path:    "/foo/bar/",
			want:    []string{"foo", "bar", ""},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pathToJsonPtrs(tt.path)
			assert.Equal(t, tt.want, got)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
