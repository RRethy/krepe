package yaml

import (
	"github.com/goccy/go-yaml"
	yamlv3 "gopkg.in/yaml.v3"
)

type BytesUnmarshaler yaml.BytesUnmarshaler
type InterfaceUnmarshaler yaml.InterfaceUnmarshaler

// Unmarshal is a wrapper for yaml.Unmarshal
func Unmarshal(data []byte, v any) error {
	return yaml.Unmarshal(data, v)
}

// UnmarshalCompatibilityShim is a shim for yaml.v3 which is needed for types
// which don't support go-yaml's Unmarshaler interface.
// E.g. orderedmap.OrderedMap
func UnmarshalCompatibilityShim(data []byte, out any) error {
	return yamlv3.Unmarshal([]byte(data), out)
}
