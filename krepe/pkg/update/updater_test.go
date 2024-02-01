package update

import (
	"path/filepath"
	"testing"

	"github.com/RRethy/krepe/krepe/pkg/exec"
	"github.com/RRethy/krepe/krepe/pkg/git"
	"github.com/RRethy/krepe/krepe/pkg/pkg"
	"github.com/RRethy/krepe/krepe/pkg/writer"
	"github.com/stretchr/testify/assert"
)

const (
	packagesDirName = "../../testdata/packages"
)

type MockMerger struct {
	origin   *pkg.Package
	local    *pkg.Package
	upstream *pkg.Package
	cnt      int
	success  bool
}

func (m *MockMerger) Merge(origin, local, upstream *pkg.Package) (*pkg.Package, error) {
	if !m.success {
		return nil, assert.AnError
	}

	m.origin = origin
	m.local = local
	m.upstream = upstream
	m.cnt++
	return local, nil
}

func (m *MockMerger) Assert(t *testing.T, origin, local, upstream string, cnt int) {
	assert.Equal(t, origin, m.origin.Labels["version"])
	assert.Equal(t, local, m.local.Labels["version"])
	assert.Equal(t, upstream, m.upstream.Labels["version"])
	assert.Equal(t, cnt, m.cnt)
}

func TestUpdaterUpdate(t *testing.T) {
	tests := []struct {
		name         string
		pkgPath      string
		url          string
		newPkgName   string
		cmd          string
		wantMergeErr bool
		wantWriteErr bool
		wantErr      bool
	}{
		{
			name:       "success",
			pkgPath:    filepath.Join(packagesDirName, "updater_test_pkg_root"),
			url:        "github.com/RRethy/updater_test_pkg_upstream@v1",
			newPkgName: "foobar",
			cmd:        "true",
			wantErr:    false,
		},
		{
			name:       "invalid url",
			pkgPath:    filepath.Join(packagesDirName, "updater_test_pkg_root"),
			url:        "invalid-url",
			newPkgName: "foobar",
			cmd:        "true",
			wantErr:    true,
		},
		{
			name:       "fail clone",
			pkgPath:    filepath.Join(packagesDirName, "updater_test_pkg_root"),
			url:        "github.com/RRethy/updater_test_pkg_upstream@v1",
			newPkgName: "foobar",
			cmd:        "false",
			wantErr:    true,
		},
		{
			name:       "invalid pkg",
			pkgPath:    filepath.Join(packagesDirName, "updater_test_pkg_root"),
			url:        "github.com/RRethy/bad_krepe_file_pkg@v1",
			newPkgName: "foobar",
			cmd:        "true",
			wantErr:    true,
		},
		{
			name:       "not existing pkg import",
			pkgPath:    filepath.Join(packagesDirName, "updater_test_pkg_root"),
			url:        "github.com/RRethy/sample_pkg@v1",
			newPkgName: "",
			cmd:        "true",
			wantErr:    true,
		},
		{
			name:       "not existing origin pkg import",
			pkgPath:    filepath.Join(packagesDirName, "udpdater_test_pkg_unknown_import"),
			url:        "github.com/RRethy/updater_test_pkg_upstream@v1",
			newPkgName: "foobar",
			cmd:        "true",
			wantErr:    true,
		},
		{
			name:       "can't clone",
			pkgPath:    filepath.Join(packagesDirName, "updater_test_pkg_root"),
			url:        "github.com/RRethy/updater_test_pkg_upstream@v1",
			newPkgName: "foobar",
			cmd:        "false",
			wantErr:    true,
		},
		{
			name:         "merge error",
			pkgPath:      filepath.Join(packagesDirName, "updater_test_pkg_root"),
			url:          "github.com/RRethy/updater_test_pkg_upstream@v1",
			newPkgName:   "foobar",
			cmd:          "true",
			wantMergeErr: true,
			wantErr:      true,
		},
		{
			name:         "write fails",
			pkgPath:      filepath.Join(packagesDirName, "updater_test_pkg_root"),
			url:          "github.com/RRethy/updater_test_pkg_upstream@v1",
			newPkgName:   "foobar",
			cmd:          "true",
			wantWriteErr: true,
			wantErr:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p, err := pkg.NewPackageFromPath(test.pkgPath)
			assert.NoError(t, err)

			cmd := exec.NewExec(exec.WithCmd(test.cmd))
			git, err := git.NewGit(
				git.WithExec(cmd),
				git.WithDir(packagesDirName),
			)
			assert.Nil(t, err)

			merger := &MockMerger{
				success: !test.wantMergeErr,
			}
			writer := &writer.Mock{
				Success: !test.wantWriteErr,
			}

			updater, err := NewUpdater(
				WithGit(git),
				WithWriter(writer),
				WithMerger(merger),
			)
			assert.Nil(t, err)

			err = updater.Update(p, test.url, test.newPkgName)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				merger.Assert(t, "origin", "local", "upstream", 1)
			}
		})
	}
}
