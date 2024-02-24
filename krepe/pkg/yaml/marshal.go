package yaml

import (
	"github.com/goccy/go-yaml"
)

const (
	indent = 2
)

// TODO: try out the custom decoder for comments

// Marshal returns [yaml.Marshal] with various configurations.
func Marshal(v any) ([]byte, error) {
	return yaml.MarshalWithOptions(v, yaml.IndentSequence(true), yaml.Indent(indent))
}
