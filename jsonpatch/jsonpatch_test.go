package jsonpatch

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
