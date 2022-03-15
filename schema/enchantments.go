// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

import (
	"bytes"
	"encoding/json"
)

// enum for good types
type Enchantment uint8
const (
	Enchantment_Imbued Enchantment = 0
)

var enchantmentsToString = map[Enchantment]string {
	Enchantment_Imbued: "Imbued",
}

var enchantmentsToID = map[string]Enchantment {
	"Imbued": Enchantment_Imbued,
}

func (s Enchantment) String() string {
	return enchantmentsToString[s]
}

// MarshalJSON marshals the enum as a text string
func (s Enchantment) MarshalText() ([]byte, error) {
	buffer := bytes.NewBufferString(enchantmentsToString[s])
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a text string to the enum value
func (s *Enchantment) UnmarshalText(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = enchantmentsToID[j]
	return nil
}

// MarshalJSON marshals the enum as a quoted json string
func (s Enchantment) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(enchantmentsToString[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *Enchantment) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = enchantmentsToID[j]
	return nil
}

// MarshalYAML marshals the enum as a quoted yaml string
func (s Enchantment) MarshalYAML() (interface{}, error) {
	buffer := bytes.NewBufferString(enchantmentsToString[s])
	return buffer.Bytes(), nil
}

// UnmarshalYAML unmashals a quoted yaml string to the enum value
func (s *Enchantment) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var j string
	if err := unmarshal(&j); err != nil {
		return err
	}

	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = enchantmentsToID[j]
	return nil
}