package pkg

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	krepeFile            = "../../testdata/sample_pkg/krepe.yaml"
	badKrepeFile         = "../../testdata/bad_krepe_file_pkg/krepe.yaml"
	nonExistentKrepeFile = "../../testdata/non_existent_pkg/krepe.yaml"
)

func TestNewKrepeFromPath(t *testing.T) {
	tests := []struct {
		name               string
		file               string
		wantImportFiles    []string
		wantImportPackages []string
		wantPipelines      map[string][]string
		wantErr            bool
	}{
		{
			name: "succeeds with valid krepe file",
			file: krepeFile,
			wantImportFiles: []string{
				"deployment.yaml",
				"service.yaml",
				"ingress.yaml",
			},
			wantImportPackages: nil,
			wantPipelines: map[string][]string{
				"mypipeline":  {"set-labels"},
				"badpipeline": {"run-fail"},
			},
			wantErr: false,
		},
		{
			name:               "fails with invalid krepe file",
			file:               badKrepeFile,
			wantImportFiles:    nil,
			wantImportPackages: nil,
			wantPipelines:      nil,
			wantErr:            true,
		},
		{
			name:               "fails with non-existent krepe file",
			file:               nonExistentKrepeFile,
			wantImportFiles:    nil,
			wantImportPackages: nil,
			wantPipelines:      nil,
			wantErr:            true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, err := NewKrepeFromPath(tt.file)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, k)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, k)
				assert.Equal(t, tt.wantImportFiles, k.Imports.Files)
				var gotPkgs []string
				for _, pkg := range k.Imports.Packages {
					gotPkgs = append(gotPkgs, *pkg.Url)
				}
				assert.Equal(t, tt.wantImportPackages, gotPkgs)
				gotPipelines := make(map[string][]string)
				for pair := k.Pipelines.Oldest(); pair != nil; pair = pair.Next() {
					pname := pair.Key
					p := pair.Value
					for _, step := range p {
						gotPipelines[pname] = append(gotPipelines[pname], step.Name())
					}
				}
				assert.Equal(t, tt.wantPipelines, gotPipelines)
			}
		})
	}
}

func XTestKrepeWrite(t *testing.T) {
	krepe, _ := NewKrepeFromPath(krepeFile)
	wantYaml, _ := os.ReadFile(krepeFile)

	tmpDir := t.TempDir()
	err := krepe.Write(tmpDir)
	assert.NoError(t, err)

	got, err := os.ReadFile(tmpDir + "/krepe.yaml")
	assert.NoError(t, err)
	assert.Equal(t, string(wantYaml), string(got))
}
