package function

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunFail(t *testing.T) {
	rf := &RunFail{}
	_, err := rf.WithConfigMap(map[string]any{
		"foo": "bar",
	})
	assert.NoError(t, err)
	assert.Error(t, rf.Run(nil))
}
