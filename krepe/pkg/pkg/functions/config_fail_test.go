package functions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigFail(t *testing.T) {
	cf := &ConfigFail{}
	_, err := cf.WithConfigMap(map[string]any{
		"foo": "bar",
	})
	assert.Error(t, err)
	assert.NoError(t, cf.Run(nil))
}
