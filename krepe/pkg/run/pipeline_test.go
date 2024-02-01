package run

import (
	"path/filepath"
	"testing"

	"github.com/RRethy/krepe/krepe/pkg/pkg"
	"github.com/stretchr/testify/assert"
)

const (
	packagesDirName = "../../testdata/packages"
)

type mockWriter struct {
	writeCalled bool
}

func (m *mockWriter) Write(pkg *pkg.Package) error {
	m.writeCalled = true
	return nil
}

func TestPipelineRun(t *testing.T) {
	tests := []struct {
		name         string
		pipelineName string
		wantErr      bool
		want         func(t *testing.T, pkg *pkg.Package)
	}{
		{
			name:         "valid pipeline",
			pipelineName: "mypipeline",
			wantErr:      false,
			want: func(t *testing.T, pkg *pkg.Package) {
				for _, importedPacakge := range pkg.PackageImports {
					assert.Equal(t, map[string]string{"foo": "bar"}, importedPacakge.Package.GetLabels())
				}
			},
		},
		{
			name:         "non-existent pipeline",
			pipelineName: "nonexistent",
			wantErr:      true,
		},
		{
			name:         "error pipeline",
			pipelineName: "badpipeline",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkg, err := pkg.NewPackageFromPath(filepath.Join(packagesDirName, "sample_pkg"))
			assert.NoError(t, err)
			w := &mockWriter{}
			pipeline := newPipeline(pkg, tt.pipelineName, withWriter(w))
			err = pipeline.run()
			if tt.wantErr {
				assert.Error(t, err)
				assert.False(t, w.writeCalled)
			} else {
				assert.NoError(t, err)
				tt.want(t, pkg)
				assert.True(t, w.writeCalled)
			}
		})
	}
}
