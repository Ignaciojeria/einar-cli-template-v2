package openapi

import (
	_ "embed"

	contract "github.com/Ignaciojeria/einar-contracts"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/MarceloPetrucio/go-scalar-api-reference"
)

type SchemaComponent struct {
	APIReferenceHTML string
	HelloEndpoint    *contract.Endpoint
	GreetEndpoint    *contract.Endpoint
}

//go:embed schema_file.json
var schema_file []byte

func init() {
	ioc.Registry(NewSchemaComponent)
}

func NewSchemaComponent() (SchemaComponent, error) {
	apiReferenceHTML, err := scalar.ApiReferenceHTML(&scalar.Options{
		SpecContent: string(schema_file),
	})
	if err != nil {
		return SchemaComponent{}, err
	}
	helloEndpoint, err := contract.LoadSpecEndpoint(
		contract.EndpointDetails{
			ContractData: schema_file, 
			Path:         "/hello",
			HTTPMethod:   "POST",
			ContentType:  "application/json",
		},
	)
	if err != nil {
		return SchemaComponent{}, err
	}
	greetEndpoint, err := contract.LoadSpecEndpoint(
		contract.EndpointDetails{
			ContractData: schema_file,
			Path:         "/greet",
			HTTPMethod:   "POST",
			ContentType:  "application/json",
		},
	)
	if err != nil {
		return SchemaComponent{}, err
	}
	return SchemaComponent{
		APIReferenceHTML: apiReferenceHTML,
		HelloEndpoint:    helloEndpoint,
		GreetEndpoint:    greetEndpoint}, nil
}
