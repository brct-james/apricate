// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

import (
	"bytes"
	"encoding/json"
)

// enum for farm bonuses
type BuildingTypes uint8
const (
	Building_Home BuildingTypes = 0
	Building_Field BuildingTypes = 1
	Building_Altar BuildingTypes = 2
	Building_SurveyOffice BuildingTypes = 3
	Building_Pasture BuildingTypes = 4
	Building_Barn BuildingTypes = 5
	Building_Kitchen BuildingTypes = 6
	Building_Silo BuildingTypes = 7
)

func (s BuildingTypes) String() string {
	return buildingsToString[s]
}

var buildingsToString = map[BuildingTypes]string {
	Building_Home: "Home",
	Building_Field: "Field",
	Building_Altar: "Altar",
	Building_SurveyOffice: "Survey Office",
	Building_Pasture: "Pasture",
	Building_Barn: "Barn",
	Building_Kitchen: "Kitchen",
	Building_Silo: "Silo",
}

var buildingsToID = map[string]BuildingTypes {
	"Home": Building_Home,
	"Field": Building_Field,
	"Altar": Building_Altar,
	"Survey Office": Building_SurveyOffice,
	"Pasture": Building_Pasture,
	"Barn": Building_Barn,
	"Kitchen": Building_Kitchen,
	"Silo": Building_Silo,
}

// MarshalJSON marshals the enum as a text string
func (s BuildingTypes) MarshalText() ([]byte, error) {
	buffer := bytes.NewBufferString(buildingsToString[s])
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a text string to the enum value
func (s *BuildingTypes) UnmarshalText(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = buildingsToID[j]
	return nil
}

// MarshalJSON marshals the enum as a quoted json string
func (s BuildingTypes) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(buildingsToString[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *BuildingTypes) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = buildingsToID[j]
	return nil
}