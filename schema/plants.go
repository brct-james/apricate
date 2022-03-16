// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

import (
	"apricate/filemngr"
	"apricate/log"
	"bytes"
	"encoding/json"

	"gopkg.in/yaml.v3"
)

// Defines a plant object to be stored on plots
type Plant struct {
	PlantType PlantType `json:"type" binding:"required"`
	CurrentStage int16 `json:"current_stage" binding:"required"`
	Yield float32 `json:"yield" binding:"required"`
	NextStageTimestamp uint64 `json:"next_stage_timestamp" binding:"required"`
}

func NewPlant(ptype PlantType) *Plant {
	return &Plant{
		PlantType: ptype,
		CurrentStage: 0,
		Yield: 1.0,
		NextStageTimestamp: 0,
	}
}

// Defines a plant definition for the plant dictionary
type PlantDefinition struct {
	Name PlantType `yaml:"Name" json:"name" binding:"required"`
	ProduceName GoodType `yaml:"ProduceName" json:"produce_name" binding:"required"`
	SeedName GoodType `yaml:"SeedName" json:"seed_name" binding:"required"`
	Description string `yaml:"Description" json:"description" binding:"required"`
	GrowthStages []GrowthStage `yaml:"GrowthStages" json:"growth_stages" binding:"required"`
}

// Load plant struct by unmarhsalling given yaml file
func Plants_load(path_to_plants_yaml string) map[string]PlantDefinition {
	plantsBytes := filemngr.ReadFileToBytes(path_to_plants_yaml)
	var plants map[string]PlantDefinition
	err := yaml.Unmarshal(plantsBytes, &plants)
	if err != nil {
		log.Error.Fatalf("%v", err.(*json.SyntaxError))
		// log.Error.Fatalf("%v", err.(*yaml.TypeError))
	}
	return plants
}

// enum for plant types
type PlantType uint8
const (
	Plant_Cabbage PlantType = 0
	Plant_ShelvisFig PlantType = 1
	Plant_Potato PlantType = 2
)

var plantsToString = map[PlantType]string {
	Plant_Cabbage: "Cabbage",
	Plant_ShelvisFig: "Shelvis Fig",
	Plant_Potato: "Potato",
}

var plantsToID = map[string]PlantType {
	"Cabbage": Plant_Cabbage,
	"Shelvis Fig": Plant_ShelvisFig,
	"Potato": Plant_Potato,
}

func (s PlantType) String() string {
	return plantsToString[s]
}

// MarshalJSON marshals the enum as a text string
func (s PlantType) MarshalText() ([]byte, error) {
	buffer := bytes.NewBufferString(plantsToString[s])
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a text string to the enum value
func (s *PlantType) UnmarshalText(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = plantsToID[j]
	return nil
}

// MarshalJSON marshals the enum as a quoted json string
func (s PlantType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(plantsToString[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *PlantType) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = plantsToID[j]
	return nil
}

// MarshalYAML marshals the enum as a quoted yaml string
func (s PlantType) MarshalYAML() (interface{}, error) {
	buffer := bytes.NewBufferString(plantsToString[s])
	return buffer.Bytes(), nil
}

// UnmarshalYAML unmashals a quoted yaml string to the enum value
func (s *PlantType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var j string
	if err := unmarshal(&j); err != nil {
		return err
	}

	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = plantsToID[j]
	return nil
}