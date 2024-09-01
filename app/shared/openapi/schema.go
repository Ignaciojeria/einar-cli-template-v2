package openapi

import (
	_ "embed"

	contract "github.com/Ignaciojeria/einar-contracts"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type SchemaComponent struct {
	*contract.APISpec
}

//go:embed schema_file.json
var schema_file []byte

func init() {
	ioc.Registry(NewSchemaComponent)
}

func NewSchemaComponent() (SchemaComponent, error) {
	spec, err := contract.NewAPIAPISpec(
		contract.Contract{
			Data:        schema_file,
			Path:        "/hello",
			HTTPMethod:  "POST",
			ContentType: "application/json",
		},
	)
	if err != nil {
		return SchemaComponent{}, err
	}
	return SchemaComponent{spec}, nil
}
