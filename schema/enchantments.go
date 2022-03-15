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

var enchantmentToString = map[Enchantment]string {
	Enchantment_Imbued: "Imbued",
}

var enchantmentToId = map[string]Enchantment {
	"Imbued": Enchantment_Imbued,
}

func (s Enchantment) String() string {
	return enchantmentToString[s]
}

// MarshalJSON marshals the enum as a text string
func (s Enchantment) MarshalText() ([]byte, error) {
	buffer := bytes.NewBufferString(enchantmentToString[s])
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
	*s = enchantmentToId[j]
	return nil
}

// MarshalJSON marshals the enum as a quoted json string
func (s Enchantment) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(enchantmentToString[s])
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
	*s = enchantmentToId[j]
	return nil
}