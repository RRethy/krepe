package run

import (
	"path/filepath"
	"testing"

	"github.com/RRethy/krepe/krepe/pkg/pkg"
	"github.com/RRethy/krepe/krepe/pkg/writer"
	"github.com/stretchr/testify/assert"
)

const (
	packagesDirName = "../../testdata/packages"
)

func TestPipelineRun(t *testing.T) {
	tests := []struct {
		name         string
		pipelineName string
		wantErr      bool
		wantWriteErr bool
		assert       func(t *testing.T, pkg *pkg.Package)
	}{
		{
			name:         "valid pipeline",
			pipelineName: "mypipeline",
			wantErr:      false,
			assert: func(t *testing.T, pkg *pkg.Package) {
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
		{
			name:         "write fails",
			pipelineName: "mypipeline",
			wantErr:      true,
			wantWriteErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := &writer.Mock{Success: !test.wantWriteErr}
			pipeline := &Pipeline{Writer: w}
			err := pipeline.Run(filepath.Join(packagesDirName, "sample_pkg"), test.pipelineName)
			if test.wantErr {
				assert.Error(t, err)
				assert.Equal(t, 0, w.Cnt)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, 1, w.Cnt)
				test.assert(t, w.Pkg)
			}
		})
	}
}
