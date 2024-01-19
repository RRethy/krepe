package pipeline

import (
	"bytes"

	"github.com/goccy/go-yaml"
	"github.com/wk8/go-ordered-map/v2"
	yamlv3 "gopkg.in/yaml.v3"
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
	return yamlv3.Unmarshal([]byte(data), p)
}

func (p *Pipelines) MarshalYAML() ([]byte, error) {
	var buf bytes.Buffer
	encoder := yamlv3.NewEncoder(&buf)
	encoder.SetIndent(2)
	err := encoder.Encode(&p.OrderedMap)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
