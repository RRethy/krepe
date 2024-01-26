package merger

import (
	"github.com/RRethy/krepe/krepe/pkg/pkg"
)

// This interface is a bit stupid but it's to show what
// TODO
type Mergeable interface {
	any | map[string]any | []any | *pkg.Pkg | *pkg.Krepe
}
