// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

import (
	"bytes"
	"encoding/json"
)

// Defines a growth stage
type GrowthStage struct {
	Name string `yaml:"Name" json:"name" binding:"required"`
	Description string `yaml:"Description" json:"description" binding:"required"`
	Action GrowthAction `yaml:"Action" json:"action" binding:"required"`
	ActionToSkip GrowthAction `yaml:"ActionToSkip" json:"action_to_skip,omitempty"`
	ConsumableOptions []GrowthConsumable `yaml:"Consumables" json:"consumable_options,omitempty"` // One of the requirements from this list must be specified in action request. Goods used from local warehouse. Quantity multiplied by plant size.
	Optional bool `yaml:"Optional" json:"optional,omitempty"` // May send Wait action to skip optional steps (growth time of optional steps is skipped as well). 
	AddedYield float32 `yaml:"AddedYield" json:"added_yield,omitempty"` // For Gigantic, Colossal and Titanic sizes, yield exclusively impacts Quality (but too a much higher extent), rather than Quantity like with smaller varietals
	GrowthTime uint64 `yaml:"GrowthTime" json:"growth_time,omitempty"`
	Harvestable GrowthHarvest `yaml:"Harvestable" json:"harvestable,omitempty"` // Plants may be harvested at any stage where Harvestable is present, some may have additional stages beyond first harvest opportunity
}

// Defines a growth consumable
type GrowthConsumable struct {
	Good
	AddedYield float32 `yaml:"AddedYield" json:"added_yield,omitempty"`
}

// Defines a growth harvest
type GrowthHarvest struct { // If there's room in warehouse, harvest sets Harvested to true and adds harvest to warehouse. If not final, next action may be sent instantly - no growth time
	Good GoodType `yaml:"Good" json:"good" binding:"required"`
	Seeds GoodType `yaml:"Seeds" json:"seeds,omitempty"`
	FinalHarvest bool `yaml:"FinalHarvest" json:"final_harvest,omitempty"` // If true, when harvested, clears the plot after
}

// enum for growthaction types
type GrowthAction uint8
const (
	GA_Wait GrowthAction = 0 // Special: Skips optional actions
	GA_Clear GrowthAction = 1 // Special: Clears plot, destroying plants in-progress
	GA_Water GrowthAction = 2 // Water Wand
	GA_Trim GrowthAction = 3 // Shears
	GA_Dig GrowthAction = 4 // Spade
	GA_Weed GrowthAction = 5 // Hoe
	GA_Fertilize GrowthAction = 6 // Pitchfork
	GA_Hill GrowthAction = 7 // Rake
	GA_Sprout GrowthAction = 8 // Pot
)

var growthActionsToToolTypes = map[GrowthAction]ToolTypes {
	GA_Water: Tool_WaterWand,
	GA_Trim: Tool_Shears,
	GA_Dig: Tool_Spade,
	GA_Weed: Tool_Hoe,
	GA_Fertilize: Tool_Pitchfork,
	GA_Hill: Tool_Rake,
	GA_Sprout: Tool_SproutingPot,
}

// var toolTypesToGrowthActions = map[ToolTypes]GrowthAction {
// 	Tool_WaterWand: GA_Water,
// 	Tool_Shears: GA_Trim,
// 	Tool_Spade: GA_Dig,
// 	Tool_Hoe: GA_Weed,
// 	Tool_Pitchfork: GA_Fertilize,
// 	Tool_Rake: GA_Hill,
// 	Tool_SproutingPot: GA_Sprout,
// }

var growthActionsToString = map[GrowthAction]string {
	GA_Wait: "Wait",
	GA_Clear: "Clear",
	GA_Water: "Water",
	GA_Trim: "Trim",
	GA_Dig: "Dig",
	GA_Weed: "Weed",
	GA_Fertilize: "Fertilize",
	GA_Hill: "Hill",
	GA_Sprout: "Sprout",
}

var growthActionsToID = map[string]GrowthAction {
	"Wait": GA_Wait,
	"Clear": GA_Clear,
	"Water": GA_Water,
	"Trim": GA_Trim,
	"Dig": GA_Dig,
	"Weed": GA_Weed,
	"Fertilize": GA_Fertilize,
	"Hill": GA_Hill,
	"Sprout": GA_Sprout,
}

func (s GrowthAction) String() string {
	return growthActionsToString[s]
}

func (s GrowthAction) ToolType() ToolTypes {
	return growthActionsToToolTypes[s]
}

// MarshalJSON marshals the enum as a text string
func (s GrowthAction) MarshalText() ([]byte, error) {
	buffer := bytes.NewBufferString(growthActionsToString[s])
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a text string to the enum value
func (s *GrowthAction) UnmarshalText(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = growthActionsToID[j]
	return nil
}

// MarshalJSON marshals the enum as a quoted json string
func (s GrowthAction) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(growthActionsToString[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *GrowthAction) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = growthActionsToID[j]
	return nil
}

// MarshalYAML marshals the enum as a quoted yaml string
func (s GrowthAction) MarshalYAML() (interface{}, error) {
	buffer := bytes.NewBufferString(growthActionsToString[s])
	return buffer.Bytes(), nil
}

// UnmarshalYAML unmashals a quoted yaml string to the enum value
func (s *GrowthAction) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var j string
	if err := unmarshal(&j); err != nil {
		return err
	}

	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = growthActionsToID[j]
	return nil
}