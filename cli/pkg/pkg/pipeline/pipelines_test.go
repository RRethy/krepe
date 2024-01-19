package pipeline

import (
	"testing"

	"github.com/Shopify/krepe/cli/pkg/pkg/function"
	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert"
	"github.com/wk8/go-ordered-map/v2"
)

const yamlStr = `foo:
  - function: succeed
baz:
  - function: succeed
zxy:
  - function: succeed
abc:
  - function: succeed
`

func TestPipelinesUnmarshalYAML(t *testing.T) {
	p := &Pipelines{}
	err := yaml.Unmarshal([]byte(yamlStr), p)
	assert.NoError(t, err)
	assert.Equal(t, 4, p.Len())
	expectedOrderedKeys := []string{"foo", "baz", "zxy", "abc"}
	var gotOrderedKeys []string
	for pair := p.Oldest(); pair != nil; pair = pair.Next() {
		gotOrderedKeys = append(gotOrderedKeys, pair.Key)
	}
	assert.Equal(t, expectedOrderedKeys, gotOrderedKeys)
}

func TestPipelinesMarshalYAML(t *testing.T) {
	p := Pipelines{
		*orderedmap.New[string, Pipeline](),
	}
	steps := Pipeline{
		{
			name: "succeed",
			fn:   &function.Succeed{},
		},
	}
	p.Set("foo", steps)
	p.Set("baz", steps)
	p.Set("zxy", steps)
	p.Set("abc", steps)
	got, err := yaml.Marshal(&p)
	assert.NoError(t, err)
	assert.Equal(t, yamlStr, string(got))
}
