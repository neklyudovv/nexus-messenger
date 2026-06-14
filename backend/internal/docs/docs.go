package docs

import _ "embed"

//go:embed openapi.yaml
var OpenAPISpec []byte

//go:embed scalar.html
var ScalarHTML []byte
