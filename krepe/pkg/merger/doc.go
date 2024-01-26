// Package threewaymerge provides a 3-way merge for data structures that would be produced from parsing YAML.
//
// At a high level, the 3-way merge algorithm will return the local value if it is different than origin iff upstream != origin.
// Otherwise, it will return upstream.
//
// Maps and associative slices are merged recursively by key.
// Scalars are merged by returning upstream iff upstream == origin, otherwise local.
// Non-associative slices are treated as scalars.
// Associative slices are recursively merged by associative key.
//
// Values of different types, including nil, are treated as scalars.
// They are still still merged recursively when possible.
//
// An associative slice is a slice of maps where each map has an associative key that is unique across all maps.
// The following are possible associative keys:
//   - mountPath
//   - devicePath
//   - ip
//   - type
//   - topologyKey
//   - name
//   - containerPort
package merger
