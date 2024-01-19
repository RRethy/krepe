package jsonpatch

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.design/x/reflect"
)

type patchTest struct {
	name    string
	obj     map[string]any
	patch   JsonPatch
	want    any
	wantErr bool
}

func (pt *patchTest) run(t *testing.T) {
	t.Run(pt.name, func(t *testing.T) {
		objCopy := reflect.DeepCopy(pt.obj)
		got, err := pt.patch.Apply(pt.obj)
		if pt.wantErr {
			assert.Error(t, err)
			assert.Equal(t, objCopy, pt.obj)
			assert.Nil(t, got)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, pt.want, got)
		}
	})
}

func runPatchTests(t *testing.T, tests []*patchTest) {
	t.Helper()

	for _, tt := range tests {
		tt.run(t)
	}
}

func TestNewJsonPatch(t *testing.T) {
	tests := []struct {
		name    string
		op      string
		from    string
		path    string
		value   any
		want    JsonPatch
		wantErr bool
	}{
		{
			name:    "Test with add op",
			op:      "add",
			from:    "",
			path:    "/foo/bar",
			value:   "baz",
			want:    &Add{path: []string{"foo", "bar"}, value: "baz"},
			wantErr: false,
		},
		{
			name:    "Test with remove op",
			op:      "remove",
			from:    "",
			path:    "/foo/bar",
			value:   nil,
			want:    &Remove{path: []string{"foo", "bar"}},
			wantErr: false,
		},
		{
			name:    "Test with replace op",
			op:      "replace",
			from:    "",
			path:    "/foo/bar",
			value:   "baz",
			want:    &Replace{path: []string{"foo", "bar"}, value: "baz"},
			wantErr: false,
		},
		{
			name:    "Test with move op",
			op:      "move",
			from:    "/foo/bar",
			path:    "/foo/baz",
			value:   nil,
			want:    &Move{from: []string{"foo", "bar"}, path: []string{"foo", "baz"}},
			wantErr: false,
		},
		{
			name:    "Test with copy op",
			op:      "copy",
			from:    "/foo/bar",
			path:    "/foo/baz",
			value:   nil,
			want:    &Copy{from: []string{"foo", "bar"}, path: []string{"foo", "baz"}},
			wantErr: false,
		},
		{
			name:    "Test with test op",
			op:      "test",
			from:    "",
			path:    "/foo/bar",
			value:   5,
			want:    &Test{path: []string{"foo", "bar"}, value: 5},
			wantErr: false,
		},
		{
			name:    "Test with unknown op",
			op:      "unknown",
			from:    "/foo/bar",
			path:    "/foo/baz",
			value:   nil,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewJsonPatch(tt.op, tt.from, tt.path, tt.value)
			assert.Equal(t, tt.want, got)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
