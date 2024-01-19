package reporef

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseRepoRef(t *testing.T) {
	tests := []struct {
		name    string
		repoRef string
		want    *RepoRef
		wantErr bool
	}{
		{
			name:    "valid repoRef",
			repoRef: "github.com/RRethy/krepe@v0.0.1",
			want: &RepoRef{
				URL:  "github.com/RRethy/krepe",
				Path: "",
				Name: "krepe",
				Tag:  "v0.0.1",
			},
			wantErr: false,
		},
		{
			name:    "valid repoRef with path",
			repoRef: "github.com/RRethy/krepe/pkg/repoRef@v0.0.1",
			want: &RepoRef{
				URL:  "github.com/RRethy/krepe",
				Path: "pkg/repoRef",
				Name: "repoRef",
				Tag:  "v0.0.1",
			},
			wantErr: false,
		},
		{
			name:    "invalid repoRef with missing tag",
			repoRef: "github.com/RRethy/krepe",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid repoRef with incomplete url",
			repoRef: "github.com/RRethy@v0.0.1",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRepoRef(tt.repoRef)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
