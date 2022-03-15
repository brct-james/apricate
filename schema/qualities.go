// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

import (
	"bytes"
	"encoding/json"
)

// enum for qualities
type Quality uint8
const (
	Quality_Unremarkable Quality = 0 // 60% of items in the world
	Quality_Fine Quality = 1 // 25% of items in the world
	Quality_Exquisite Quality = 2 // 10% of items in the world
	Quality_Supreme Quality = 3 // 4% of items in the world
	Quality_Perfect Quality = 4 // 1% of items in the world
)

var qualitysToString = map[Quality]string {
	Quality_Unremarkable: "Unremarkable",
	Quality_Fine: "Fine",
	Quality_Exquisite: "Exquisite",
	Quality_Supreme: "Supreme",
	Quality_Perfect: "Perfect",
}

var qualitysToID = map[string]Quality {
	"Unremarkable": Quality_Unremarkable,
	"Fine": Quality_Fine,
	"Exquisite": Quality_Exquisite,
	"Supreme": Quality_Supreme,
	"Perfect": Quality_Perfect,
}

func (s Quality) String() string {
	return qualitysToString[s]
}

// MarshalJSON marshals the enum as a text string
func (s Quality) MarshalText() ([]byte, error) {
	buffer := bytes.NewBufferString(qualitysToString[s])
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a text string to the enum value
func (s *Quality) UnmarshalText(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = qualitysToID[j]
	return nil
}

// MarshalJSON marshals the enum as a quoted json string
func (s Quality) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(qualitysToString[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *Quality) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = qualitysToID[j]
	return nil
}

// MarshalYAML marshals the enum as a quoted yaml string
func (s Quality) MarshalYAML() (interface{}, error) {
	buffer := bytes.NewBufferString(qualitysToString[s])
	return buffer.Bytes(), nil
}

// UnmarshalYAML unmashals a quoted yaml string to the enum value
func (s *Quality) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var j string
	if err := unmarshal(&j); err != nil {
		return err
	}

	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = qualitysToID[j]
	return nil
}