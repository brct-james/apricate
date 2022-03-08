// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

import (
	"bytes"
	"encoding/json"
)

// enum for assistant types
type Size uint16
const (
	Miniscule Size = 1
	Tiny Size = 2
	Small Size = 4
	Average Size = 8
	Large Size = 16
	Huge Size = 32
	Gigantic Size = 64
	Colossal Size = 128
	Titanic Size = 256
)

func (s Size) String() string {
	return sizeToString[s]
}

var sizeToString = map[Size]string {
	Miniscule: "Miniscule (1)",
	Tiny: "Tiny (2)",
	Small: "Small (4)",
	Average: "Average (8)",
	Large: "Large (16)",
	Huge: "Huge (32)",
	Gigantic: "Gigantic (64)",
	Colossal: "Colossal (128)",
	Titanic: "Titanic (256)",
}

var sizeToID = map[string]Size {
	"Miniscule (1)": Miniscule,
	"Tiny (2)": Tiny,
	"Small (4)": Small,
	"Average (8)": Average,
	"Large (16)": Large,
	"Huge (32)": Huge,
	"Gigantic (64)": Gigantic,
	"Colossal (128)": Colossal,
	"Titanic (256)": Titanic,
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