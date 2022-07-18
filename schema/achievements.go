// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

import (
	"bytes"
	"encoding/json"
)

// enum for api response codes
type Achievement int16
const (
	Achievement_Owner Achievement = -1
	Achievement_Contributor Achievement = 0
	Achievement_Noob Achievement = 1
)

func (s Achievement) String() string {
	return achievementToString[s]
}

var achievementToString = map[Achievement]string {
	Achievement_Owner: "Owner",
	Achievement_Contributor: "Contributor",
	Achievement_Noob: "Noob",
}

var achievementToID = map[string]Achievement {
	"Owner": Achievement_Owner,
	"Contributor": Achievement_Contributor,
	"Noob": Achievement_Noob,
}

// MarshalJSON marshals the enum as a quoted json string
func (s Achievement) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(achievementToString[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *Achievement) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = achievementToID[j]
	return nil
}