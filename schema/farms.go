// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

import (
	"bytes"
	"encoding/json"
)

// enum for farm bonuses
type FarmBonuses uint8
const (
	FarmBonus_GoodSoil FarmBonuses = 0
	FarmBonus_Portals FarmBonuses = 1
	FarmBonus_Forested FarmBonuses = 2
)

// Defines a farm
type Farm struct {
	UID uint64 `json:"uid" binding:"required"`
	Location *Location `json:"location" binding:"required"`
	Bonuses []FarmBonuses `json:"bonuses" binding:"required"`
	Tools map[string]uint8 `json:"tools" binding:"required"`
	Buildings map[string]uint8 `json:"buildings" binding:"required"`
	Plots []uint64 `json:"plots" binding:"required"`
}

func (s FarmBonuses) String() string {
	return farmBonusesToString[s]
}

var farmBonusesToString = map[FarmBonuses]string {
	FarmBonus_GoodSoil: "Good Soil | Doubles yields for plants of Average size and below. Plants Large and above grow 1.5x as fast.",
	FarmBonus_Portals: "Portals | Halve the travel time for any Assistant travelling to OR from the farm.",
	FarmBonus_Forested: "Forested | Two unique Buildings, the Logging Camp and Lumber Mill, are available for construction. Also comes with a free Axe tool, allowing the on-site harvest and processing of logs, lumber, and planks.",
}

var farmBonusesToID = map[string]FarmBonuses {
	"Good Soil | Doubles yields for plants of Average size and below. Plants Large and above grow 1.5x as fast.": FarmBonus_GoodSoil,
	"Portals | Halve the travel time for any Assistant travelling to OR from the farm.": FarmBonus_Portals,
	"Forested | Two unique Buildings, the Logging Camp and Lumber Mill, are available for construction. Also comes with a free Axe tool, allowing the on-site harvest and processing of logs, lumber, and planks.": FarmBonus_Forested,
}

// MarshalJSON marshals the enum as a quoted json string
func (s FarmBonuses) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(farmBonusesToString[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *FarmBonuses) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = farmBonusesToID[j]
	return nil
}