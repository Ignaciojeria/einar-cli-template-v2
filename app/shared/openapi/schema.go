package openapi

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/qri-io/jsonschema"
)

type SchemaComponent struct {
	schema *jsonschema.Schema
}

//go:embed schema_file.yaml
var schema_file []byte

func init() {
	ioc.Registry(NewSchemaComponent)
}

func NewSchemaComponent() (SchemaComponent, error) {
	const SCHEMA_KEY = "REPLACE_BY_YOUR_SCHEMA_KEY"
	loader := openapi3.NewLoader()
	data, err := loader.LoadFromData(schema_file)
	if err != nil {
		return SchemaComponent{}, err
	}
	schemaComponents := data.Components
	if schemaComponents == nil {
		return SchemaComponent{}, errors.New("schema components not found")
	}
	schemaRef := schemaComponents.Schemas[SCHEMA_KEY]
	if schemaRef == nil {
		return SchemaComponent{}, errors.New(SCHEMA_KEY + " schemaRef not found")
	}
	schema_file, err = schemaRef.MarshalJSON()
	if err != nil {
		return SchemaComponent{}, err
	}

	schema := &jsonschema.Schema{}
	if err := json.Unmarshal(schema_file, schema); err != nil {
		return SchemaComponent{}, err
	}
	return SchemaComponent{schema: schema}, nil
}

func (v *SchemaComponent) ValidateBytes(ctx context.Context, json []byte) error {
	errs, err := v.schema.ValidateBytes(ctx, json)
	if err != nil {
		return err
	}
	if len(errs) == 0 {
		return nil
	}
	var sb strings.Builder
	for _, e := range errs {
		sb.WriteString(fmt.Sprintf("Validation error: %s\n", e.Error()))
	}
	return errors.New(sb.String())
}
