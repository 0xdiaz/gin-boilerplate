package api

import _ "embed"

// OpenAPIYAML is the OpenAPI 3.0 spec for the API. Served at GET /swagger/openapi.yaml.
//
//go:embed openapi.yaml
var OpenAPIYAML []byte
