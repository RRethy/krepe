package jsonpatch

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
