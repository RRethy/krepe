package functions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSucceed(t *testing.T) {
	s := &Succeed{}
	_, err := s.WithConfigMap(map[string]any{
		"foo": "bar",
	})
	assert.NoError(t, err)
	assert.NoError(t, s.Run(nil))
}
