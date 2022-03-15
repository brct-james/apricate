// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

import (
	"bytes"
	"encoding/json"
)

// enum for assistant types
type Size uint16
const (
	Miniature Size = 1
	Tiny Size = 2
	Small Size = 4
	Modest Size = 8
	Average Size = 16
	Large Size = 32
	Huge Size = 64
	Gigantic Size = 256
	Colossal Size = 1024
	Titanic Size = 4096
)

func (s Size) String() string {
	return sizeToString[s]
}

var sizeToString = map[Size]string {
	Miniature: "Miniature (1)",
	Tiny: "Tiny (2)",
	Small: "Small (4)",
	Modest: "Modest (8)",
	Average: "Average (16)",
	Large: "Large (32)",
	Huge: "Huge (64)",
	Gigantic: "Gigantic (256)",
	Colossal: "Colossal (1024)",
	Titanic: "Titanic (4096)",
}

var sizeToID = map[string]Size {
	"Miniature (1)": Miniature,
	"Tiny (2)": Tiny,
	"Small (4)": Small,
	"Modest (8)": Modest,
	"Average (16)": Average,
	"Large (32)": Large,
	"Huge (64)": Huge,
	"Gigantic (256)": Gigantic,
	"Colossal (1024)": Colossal,
	"Titanic (4096)": Titanic,
}

// MarshalJSON marshals the enum as a quoted json string
func (s Size) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(sizeToString[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *Size) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = sizeToID[j]
	return nil
}

// MarshalYAML marshals the enum as a quoted yaml string
func (s Size) MarshalYAML() (interface{}, error) {
	buffer := bytes.NewBufferString(sizeToString[s])
	return buffer.Bytes(), nil
}

// UnmarshalYAML unmashals a quoted yaml string to the enum value
func (s *Size) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var j string
	if err := unmarshal(&j); err != nil {
		return err
	}

	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = sizeToID[j]
	return nil
}