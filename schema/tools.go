// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

import (
	"bytes"
	"encoding/json"
)

// enum for farm bonuses
type ToolTypes uint8
const (
	Tool_Spade ToolTypes = 0
	Tool_Hoe ToolTypes = 1
	Tool_Rake ToolTypes = 2
	Tool_Pitchfork ToolTypes = 3
	Tool_Shears ToolTypes = 4
	Tool_WaterWand ToolTypes = 5
	Tool_Knife ToolTypes = 6
	Tool_PestleAndMortar ToolTypes = 7
	Tool_DryingRack ToolTypes = 8
	Tool_SproutingPot ToolTypes = 9
	Tool_ShadeScroll ToolTypes = 10
	Tool_Sickle ToolTypes = 11
	Tool_SpiritFlute ToolTypes = 12
	Tool_BindEvilScroll ToolTypes = 13
)

func (s ToolTypes) String() string {
	return toolTypesToString[s]
}

func StringToTool(s string) ToolTypes {
	return toolTypesToID[s]
}

var toolTypesToString = map[ToolTypes]string {
	Tool_Spade: "Spade",
	Tool_Hoe: "Hoe",
	Tool_Rake: "Rake",
	Tool_Pitchfork: "Pitchfork",
	Tool_Shears: "Shears",
	Tool_WaterWand: "Water Wand",
	Tool_Knife: "Knife",
	Tool_PestleAndMortar: "Pestle and Mortar",
	Tool_DryingRack: "Drying Rack",
	Tool_SproutingPot: "Sprouting Pot",
	Tool_ShadeScroll: "Scroll of Hyperspecific Cloud Cover",
	Tool_Sickle: "Sickle",
	Tool_SpiritFlute: "Spirit Flute",
	Tool_BindEvilScroll: "Scroll of Bind Evil",
}

var toolTypesToID = map[string]ToolTypes {
	"Spade": Tool_Spade,
	"Hoe": Tool_Hoe,
	"Rake": Tool_Rake,
	"Pitchfork": Tool_Pitchfork,
	"Shears": Tool_Shears,
	"Water Wand": Tool_WaterWand,
	"Knife": Tool_Knife,
	"Pestle and Mortar": Tool_PestleAndMortar,
	"Drying Rack": Tool_DryingRack,
	"Sprouting Pot": Tool_SproutingPot,
	"Scroll of Hyperspecific Cloud Cover": Tool_ShadeScroll,
	"Sickle": Tool_Sickle,
	"Spirit Flute": Tool_SpiritFlute,
	"Scroll of Bind Evil": Tool_BindEvilScroll,
}

// MarshalJSON marshals the enum as a text string
func (s ToolTypes) MarshalText() ([]byte, error) {
	buffer := bytes.NewBufferString(toolTypesToString[s])
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a text string to the enum value
func (s *ToolTypes) UnmarshalText(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = toolTypesToID[j]
	return nil
}

// MarshalJSON marshals the enum as a quoted json string
func (s ToolTypes) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toolTypesToString[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *ToolTypes) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = toolTypesToID[j]
	return nil
}

// MarshalYAML marshals the enum as a quoted yaml string
func (s ToolTypes) MarshalYAML() (interface{}, error) {
	buffer := bytes.NewBufferString(toolTypesToString[s])
	return buffer.Bytes(), nil
}

// UnmarshalYAML unmashals a quoted yaml string to the enum value
func (s *ToolTypes) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var j string
	if err := unmarshal(&j); err != nil {
		return err
	}

	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = toolTypesToID[j]
	return nil
}