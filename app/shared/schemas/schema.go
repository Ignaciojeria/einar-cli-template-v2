package schemas

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/qri-io/jsonschema"
)

type SchemaComponent struct {
	schema *jsonschema.Schema
}

//go:embed schema_file.json
var schema_file []byte

func init() {
	ioc.Registry(NewSchemaComponent)
}

func NewSchemaComponent() (SchemaComponent, error) {
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
