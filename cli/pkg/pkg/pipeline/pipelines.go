package pipeline

import (
	"github.com/Shopify/krepe/cli/pkg/yaml"
	"github.com/wk8/go-ordered-map/v2"
)

var (
	_ yaml.BytesUnmarshaler = &Pipelines{}
	_ yaml.BytesMarshaler   = &Pipelines{}
)

type Pipelines struct {
	orderedmap.OrderedMap[string, Pipeline] `yaml:",inline"`
}

func (p *Pipelines) UnmarshalYAML(data []byte) error {
	*p = Pipelines{orderedmap.OrderedMap[string, Pipeline]{}}
	return yaml.UnmarshalCompatibilityShim(data, p)
}

func (p *Pipelines) MarshalYAML() ([]byte, error) {
	return yaml.MarshalCompatibilityShim(&p.OrderedMap)
}
