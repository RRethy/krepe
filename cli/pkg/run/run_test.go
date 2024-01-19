package run

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name     string
		pkg      string
		pipeline string
		function string
		wantErr  bool
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Run(tt.pkg, tt.pipeline, tt.function)
			if tt.wantErr {
				assert.Error(t, err)
			} else {

			}
		})
	}
}
