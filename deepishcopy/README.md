# deepishcopy

TODO(RRethy): Add godoc link

Package deepishcopy provides a deep copy for data structures that would be produced from parsing YAML.

It is purpose-built for krepe, but can be used for other purposes.
Originally, we were using https://github.com/golang-design/reflect since it implements deep copy.
However, it has a deal-breaking bug that arises when parsing YAML.
See https://github.com/golang-design/reflect/issues/2 for more details.
