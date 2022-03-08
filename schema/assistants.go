// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

import (
	"bytes"
	"encoding/json"
)

// enum for assistant types
type AssistantTypes uint8
const (
	Hireling AssistantTypes = 0
	Familiar AssistantTypes = 1
	Golem AssistantTypes = 2
)

// Defines an assistant
type Assistant struct {
	UID uint64 `json:"uid" binding:"required"`
	Archetype AssistantTypes `json:"archetype" binding:"required"`
	Improvements map[string]uint8 `json:"improvements" binding:"required"`
	Location *Location `json:"location" binding:"required"`
	Route *Route `json:"route" binding:"required"`
}

func (s AssistantTypes) String() string {
	return assistantTypesToString[s]
}

var assistantTypesToString = map[AssistantTypes]string {
	Hireling: "Hireling",
	Familiar: "Familiar",
	Golem: "Golem",
}

var assistantTypesToID = map[string]AssistantTypes {
	"Hireling": Hireling,
	"Familiar": Familiar,
	"Golem": Golem,
}

// MarshalJSON marshals the enum as a quoted json string
func (s AssistantTypes) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(assistantTypesToString[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *AssistantTypes) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = assistantTypesToID[j]
	return nil
}