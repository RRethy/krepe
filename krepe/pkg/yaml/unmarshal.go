package yaml

import (
	"github.com/goccy/go-yaml"
)

var DisallowUnknownFieldOption = yaml.DecodeOption(yaml.DisallowUnknownField())

// Unmarshal return [yaml.Unmarshal].
func Unmarshal(data []byte, v any, opts ...yaml.DecodeOption) error {
	return yaml.UnmarshalWithOptions(data, v, opts...)
}
