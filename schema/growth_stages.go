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
	Action GrowthAction `yaml:"Action" json:"action,omitempty"`
	RequiredConsumableOptions []GrowthConsumable `yaml:"Consumables" json:"required_consumable_options,omitempty"`
	Optional bool `yaml:"Optional" json:"optional,omitempty"` // May send Skip action to skip optional steps (growth time of optional steps is skipped as well)
	AddedYield float32 `yaml:"AddedYield" json:"added_yield,omitempty"`
	GrowthTime uint64 `yaml:"GrowthTime" json:"growth_time,omitempty"`
	Harvestable GrowthHarvest `yaml:"Harvestable" json:"harvestable,omitempty"` // Plants may be harvested at any stage where Harvestable is present, some may have additional stages beyond first harvest opportunity
}

// Defines a growth harvest
type GrowthHarvest struct {
	Good GoodType `yaml:"Good" json:"good" binding:"required"`
	BaseYield float32 `yaml:"BaseYield" json:"base_yield" binding:"required"` // Multiplied by bonus yield for total. For Gigantic, Colossal and Titanic sizes, yield exclusively impacts Quality (but too a much higher extent), rather than Quantity like with smaller varietals
	HarvestAction GrowthAction `yaml:"HarvestAction" json:"harvest_action" binding:"required"`
}

// Defines a growth consumable
type GrowthConsumable struct {
	Good GoodType `yaml:"Good" json:"good" binding:"required"`
	Quantity uint64 `yaml:"Quantity" json:"quantity" binding:"required"`
}

// enum for growthaction types
type GrowthAction uint8
const (
	GA_Water GrowthAction = 0 // Water Wand
	GA_Trim GrowthAction = 1 // Shears
	GA_Dig GrowthAction = 2 // Spade
	GA_Weed GrowthAction = 3 // Hoe
	GA_Fertilize GrowthAction = 4 // Pitchfork
	GA_Hill GrowthAction = 5 // Rake
)

var growthActionsToString = map[GrowthAction]string {
	GA_Water: "Water",
	GA_Trim: "Trim",
	GA_Dig: "Dig",
	GA_Weed: "Weed",
	GA_Fertilize: "Fertilize",
	GA_Hill: "Hill",
}

var growthActionsToID = map[string]GrowthAction {
	"Water": GA_Water,
	"Trim": GA_Trim,
	"Dig": GA_Dig,
	"Weed": GA_Weed,
	"Fertilize": GA_Fertilize,
	"Hill": GA_Hill,
}

func (s GrowthAction) String() string {
	return growthActionsToString[s]
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