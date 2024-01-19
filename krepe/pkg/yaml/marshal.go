package yaml

import (
	"bytes"

	"github.com/goccy/go-yaml"
	yamlv3 "gopkg.in/yaml.v3"
)

const (
	indent = 2
)

// TODO: try out the custom decoder for ordered maps and comments

type BytesMarshaler yaml.BytesMarshaler
type InterfaceMarshaler yaml.InterfaceMarshaler

// Marshal is a wrapper for yaml.MarshalWithOptions which sets the indent
func Marshal(v any) ([]byte, error) {
	return yaml.MarshalWithOptions(v, yaml.IndentSequence(true), yaml.Indent(indent))
}

// MarshalCompatibilityShim is a shim for yaml.v3 which is needed for types
// which don't support go-yaml's Marshaler interface.
// E.g. orderedmap.OrderedMap
func MarshalCompatibilityShim(v any) ([]byte, error) {
	var buf bytes.Buffer
	encoder := yamlv3.NewEncoder(&buf)
	encoder.SetIndent(indent)
	err := encoder.Encode(v)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
