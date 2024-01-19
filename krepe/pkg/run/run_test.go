package run

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	samplePkgPath      = "../../testdata/packages/sample_pkg"
	nonExistentPkgPath = "../../testdata/packages/non_existent_pkg"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name     string
		pkg      string
		pipeline string
		wantErr  bool
	}{
		{
			name:     "succeeds with valid run",
			pkg:      samplePkgPath,
			pipeline: "no-op-pipeline",
			wantErr:  false,
		},
		{
			name:     "fails with invalid package",
			pkg:      nonExistentPkgPath,
			pipeline: "mypipeline",
			wantErr:  true,
		},
		{
			name:     "fails with unknown pipeline",
			pkg:      samplePkgPath,
			pipeline: "unknown",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Run(tt.pkg, tt.pipeline)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
