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
	Action *GrowthAction `yaml:"Action" json:"action" binding:"required"`
	Skippable bool `yaml:"Skippable" json:"skippable" binding:"required"` // send Skip action to skip stage (growth time of optional steps is skipped as well). 
	Repeatable bool `yaml:"Repeatable" json:"repeatable" binding:"required"` // If true, growth stage remains the same unless skipped (if not skippable, must CLEAR plot to escape the loop). Allows infinite harvests and infinite yield boosts
	ConsumableOptions []GrowthConsumable `yaml:"Consumables" json:"consumable_options,omitempty"` // One of the requirements from this list must be specified in action request. Goods used from local warehouse. Quantity multiplied by plant size.
	AddedYield float64 `yaml:"AddedYield" json:"added_yield" binding:"required"` // For Gigantic, Colossal and Titanic sizes, yield exclusively impacts Quality (but too a much higher extent), rather than Quantity like with smaller varietals
	GrowthTime *int64 `yaml:"GrowthTime" json:"growth_time" binding:"required"` // Cannot make omitempty else intentional 0s will be omitted
	Harvestable *GrowthHarvest `yaml:"Harvestable" json:"harvestable,omitempty"` // Plants may be harvested at any stage where Harvestable is present, some may have additional stages beyond first harvest opportunity
}

// Defines a growth consumable
type GrowthConsumable struct {
	Good `yaml:",inline"`
	AddedYield float64 `yaml:"AddedYield" json:"added_yield" binding:"required"`
}

// Defines a growth harvest
type GrowthHarvest struct { // If there's room in warehouse, harvest sets Harvested to true and adds harvest to warehouse. If not final, next action may be sent instantly - no growth time
	Produce map[string]float64 `yaml:"Produce" json:"produce,omitempty"`
	Seeds map[string]float64 `yaml:"Seeds" json:"seeds,omitempty"`
	Goods map[string]float64 `yaml:"Goods" json:"goods,omitempty"`
	FinalHarvest bool `yaml:"FinalHarvest" json:"final_harvest" binding:"required"` // If true, when harvested, clears the plot after
}

// enum for growthaction types
type GrowthAction uint8
const (
	GA_Skip GrowthAction = 0 // Special: Skips optional actions
	GA_Wait GrowthAction = 1 // None
	GA_Water GrowthAction = 2 // Water Wand
	GA_Trim GrowthAction = 3 // Shears
	GA_Dig GrowthAction = 4 // Spade
	GA_Weed GrowthAction = 5 // Hoe
	GA_Fertilize GrowthAction = 6 // Pitchfork
	GA_Hill GrowthAction = 7 // Rake
	GA_Sprout GrowthAction = 8 // Pot
	GA_Shade GrowthAction = 9 // Shade Cloth
	GA_Reap GrowthAction = 10 // Sickle
	GA_Tap GrowthAction = 11 // Tap
)

var growthActionsToToolTypes = map[GrowthAction]ToolTypes {
	GA_Water: Tool_WaterWand,
	GA_Trim: Tool_Shears,
	GA_Dig: Tool_Spade,
	GA_Weed: Tool_Hoe,
	GA_Fertilize: Tool_Pitchfork,
	GA_Hill: Tool_Rake,
	GA_Sprout: Tool_SproutingPot,
	GA_Shade: Tool_ShadeScroll,
	GA_Reap: Tool_Sickle,
	GA_Tap: Tool_Tap,
}

// var toolTypesToGrowthActions = map[ToolTypes]GrowthAction {
// 	Tool_WaterWand: GA_Water,
// 	Tool_Shears: GA_Trim,
//	...
// }

var growthActionsToString = map[GrowthAction]string {
	GA_Skip: "Skip",
	GA_Wait: "Wait",
	GA_Water: "Water",
	GA_Trim: "Trim",
	GA_Dig: "Dig",
	GA_Weed: "Weed",
	GA_Fertilize: "Fertilize",
	GA_Hill: "Hill",
	GA_Sprout: "Sprout",
	GA_Shade: "Shade",
	GA_Reap: "Reap",
	GA_Tap: "Tap",
}

var growthActionsToID = map[string]GrowthAction {
	"Skip": GA_Skip,
	"Wait": GA_Wait,
	"Water": GA_Water,
	"Trim": GA_Trim,
	"Dig": GA_Dig,
	"Weed": GA_Weed,
	"Fertilize": GA_Fertilize,
	"Hill": GA_Hill,
	"Sprout": GA_Sprout,
	"Shade": GA_Shade,
	"Reap": GA_Reap,
	"Tap": GA_Tap,
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