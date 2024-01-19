package pkg

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	krepeFile            = "../../testdata/packages/sample_pkg/krepe.yaml"
	badKrepeFile         = "../../testdata/packages/bad_krepe_file_pkg/krepe.yaml"
	nonExistentKrepeFile = "../../testdata/packages/non_existent_pkg/krepe.yaml"
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
				"mypipeline":     {"set-labels"},
				"badpipeline":    {"run-fail"},
				"no-op-pipeline": {"succeed"},
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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			k, err := NewKrepeFromPath(test.file)
			if test.wantErr {
				assert.Error(t, err)
				assert.Nil(t, k)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, k)
				assert.Equal(t, test.wantImportFiles, k.Imports.Files)
				var gotPkgs []string
				for _, pkg := range k.Imports.Packages {
					gotPkgs = append(gotPkgs, pkg.Name())
				}
				assert.Equal(t, test.wantImportPackages, gotPkgs)
				gotPipelines := make(map[string][]string)
				for pair := k.Pipelines.Oldest(); pair != nil; pair = pair.Next() {
					pname := pair.Key
					p := pair.Value
					for _, step := range p {
						gotPipelines[pname] = append(gotPipelines[pname], step.Name())
					}
				}
				assert.Equal(t, test.wantPipelines, gotPipelines)
			}
		})
	}
}

func TestKrepeWrite(t *testing.T) {
	krepe, _ := NewKrepeFromPath(krepeFile)
	wantYaml, _ := os.ReadFile(krepeFile)

	tmpDir := t.TempDir()
	err := krepe.Write(tmpDir)
	assert.NoError(t, err)

	got, err := os.ReadFile(tmpDir + "/krepe.yaml")
	assert.NoError(t, err)
	assert.Equal(t, string(wantYaml), string(got))
}
