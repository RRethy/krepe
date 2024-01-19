package pipeline

import (
	"testing"

	"github.com/Shopify/krepe/cli/pkg/pkg/function"
	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert"
)

func TestStepUnmarshallYAML(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		wantName string
		wantFn   function.Function
		wantErr  bool
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Step{}
			err := s.UnmarshalYAML(func(into any) error {
				return yaml.Unmarshal([]byte(tt.yaml), &into)
			})
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, s.fn)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantName, s.name)
				assert.Equal(t, tt.wantFn, s.fn)
			}
		})
	}
}
