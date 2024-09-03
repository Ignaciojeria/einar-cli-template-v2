package openapi

import (
	"testing"
)

func TestHelloWorldSchema(t *testing.T) {
	schema, err := NewSchemaComponent()
	if err != nil {
		t.Fatalf("Failed to create HelloWorldSchema: %v", err)
	}

	tests := []struct {
		name      string
		json      []byte
		wantError bool
	}{
		{
			name:      "Valid greeting",
			json:      []byte(`{"message": "Hello, World!"}`),
			wantError: false,
		},
		{
			name:      "Invalid greeting with missing message",
			json:      []byte(`{}`),
			wantError: true,
		},
		{
			name:      "Invalid greeting with wrong type",
			json:      []byte(`{"message": 12345}`),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := schema.HelloEndpoint.ValidateBodyBytes(tt.json)
			if (err != nil) != tt.wantError {
				t.Errorf("Test %s: expected error = %v, got %v", tt.name, tt.wantError, err != nil)
			}
		})
	}
}
