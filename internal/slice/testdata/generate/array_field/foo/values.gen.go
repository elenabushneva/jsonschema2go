// Code generated by jsonschema2go. DO NOT EDIT.
package foo

import (
	"encoding/json"
	"fmt"
	"github.com/jwilner/jsonschema2go/pkg/boxed"
)

// generated from https://example.com/testdata/generate/array_field/foo/example.json
type Example struct {
	Options ExampleOptions `json:"options,omitempty"`
}

func (m *Example) Validate() error {
	if err := m.Options.Validate(); err != nil {
		if err, ok := err.(valErr); ok {
			return &validationError{
				errType:  err.ErrType(),
				message:  err.Message(),
				path:     append([]interface{}{"Options"}, err.Path()...),
				jsonPath: append([]interface{}{"options"}, err.JSONPath()...),
			}
		}
		return err
	}
	return nil
}

// generated from https://example.com/testdata/generate/array_field/foo/inner.json
type Inner struct {
	Name  boxed.String `json:"name"`
	Value interface{}  `json:"value,omitempty"`
}

func (m *Inner) Validate() error {
	return nil
}

func (m *Inner) MarshalJSON() ([]byte, error) {
	inner := struct {
		Name  *string     `json:"name,omitempty"`
		Value interface{} `json:"value,omitempty"`
	}{
		Value: m.Value,
	}
	if m.Name.Set {
		inner.Name = &m.Name.String
	}
	return json.Marshal(inner)
}

// generated from https://example.com/testdata/generate/array_field/foo/example.json#/properties/options
type ExampleOptions []Inner

func (m ExampleOptions) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte(`[]`), nil
	}
	return json.Marshal([]Inner(m))
}

func (m ExampleOptions) Validate() error {
	for i := range m {
		if err := m[i].Validate(); err != nil {
			if err, ok := err.(valErr); ok {
				return &validationError{
					errType:  err.ErrType(),
					message:  err.Message(),
					path:     append([]interface{}{i}, err.Path()...),
					jsonPath: append([]interface{}{i}, err.JSONPath()...),
				}
			}
			return err
		}
	}
	return nil
}

type valErr interface {
	ErrType() string
	JSONPath() []interface{}
	Path() []interface{}
	Message() string
}

type validationError struct {
	errType, message string
	jsonPath, path   []interface{}
}

func (e *validationError) ErrType() string {
	return e.errType
}

func (e *validationError) JSONPath() []interface{} {
	return e.jsonPath
}

func (e *validationError) Path() []interface{} {
	return e.path
}

func (e *validationError) Message() string {
	return e.message
}

func (e *validationError) Error() string {
	return fmt.Sprintf("%v: %v", e.path, e.message)
}

var _ valErr = new(validationError)
