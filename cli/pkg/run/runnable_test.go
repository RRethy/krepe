package run

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	nonExistentPkgPath = "../../testdata/non_existent_pkg"
)

func TestNewRunnable(t *testing.T) {
	tests := []struct {
		name     string
		pkgPath  string
		pipeline string
		function string
		want     func(got runnable) bool
		wantErr  bool
	}{
		{
			name:     "succeeds with valid package and pipeline",
			pkgPath:  samplePkgPath,
			pipeline: "mypipeline",
			function: "",
			want: func(got runnable) bool {
				switch got.(type) {
				case *pipeline:
					return got.(*pipeline).name == "mypipeline"
				default:
					return false
				}
			},
			wantErr: false,
		},
		{
			name:     "succeeds with valid package and pipeline=. and function",
			pkgPath:  samplePkgPath,
			pipeline: ".",
			function: "foo",
			want: func(got runnable) bool {
				switch got.(type) {
				case *function:
					return got.(*function).name == "foo"
				default:
					return false
				}
			},
			wantErr: false,
		},
		{
			name:     "succeeds with valid package and function",
			pkgPath:  samplePkgPath,
			pipeline: "",
			function: "foo",
			want: func(got runnable) bool {
				switch got.(type) {
				case *function:
					return got.(*function).name == "foo"
				default:
					return false
				}
			},
			wantErr: false,
		},
		{
			name:     "fails with both pipeline and function",
			pkgPath:  samplePkgPath,
			pipeline: "mypipeline",
			function: "foo",
			want:     nil,
			wantErr:  true,
		},
		{
			name:    "fails with non-existent package",
			pkgPath: nonExistentPkgPath,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := newRunnable(tt.pkgPath, tt.pipeline, tt.function)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, tt.want(r))
			}
		})
	}
}
