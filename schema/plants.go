// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

import (
	"apricate/filemngr"
	"apricate/log"
	"fmt"

	"gopkg.in/yaml.v3"
)

// Defines a plant object to be stored on plots
type Plant struct {
	PlantType string `json:"type" binding:"required"`
	Size Size `json:"size" binding:"required"`
	CurrentStage int16 `json:"current_stage" binding:"required"`
	Yield float64 `json:"yield" binding:"required"`
}

func NewPlant(ptype string, size Size) *Plant {
	return &Plant{
		PlantType: ptype,
		Size: size,
		CurrentStage: 0,
		Yield: 1.0,
	}
}

// Defines a plant definition for the plant dictionary
type PlantDefinition struct {
	Name string `yaml:"Name" json:"name" binding:"required"`
	Description string `yaml:"Description" json:"description" binding:"required"`
	GrowthStages []GrowthStage `yaml:"GrowthStages" json:"growth_stages" binding:"required"`
}

func (d *PlantDefinition) GetScaledGrowthStage(gsIndex int16, plantQuantity uint64, plantSize Size) (*GrowthStage, error) {
	if int(gsIndex) >= len(d.GrowthStages) {
		return nil, fmt.Errorf("growth stage index out of bounds: %d of %d", gsIndex, len(d.GrowthStages))
	}
	res := d.GrowthStages[gsIndex]
	if len(res.ConsumableOptions) == 0 {
		// No consumables to scale
		return &res, nil
	}
	for index, option := range res.ConsumableOptions {
		option.Quantity *= plantQuantity * uint64(plantSize)
		res.ConsumableOptions[index] = option
	}
	return &res, nil
}

// Load seed struct by unmarhsalling given yaml file
func Seeds_load(path_to_seeds_yaml string) map[string]string {
	seedsBytes, readErr := filemngr.ReadFileToBytes(path_to_seeds_yaml)
	if readErr != nil {
		// Essential to server start
		panic(readErr)
	}
	var seeds map[string]string
	err := yaml.Unmarshal(seedsBytes, &seeds)
	if err != nil {
		log.Error.Fatalf("%v", err)
		// log.Error.Fatalf("%v", err.(*json.SyntaxError))
		// log.Error.Fatalf("%v", err.(*yaml.TypeError))
	}
	return seeds
}

// Load plant struct by unmarhsalling given yaml file
func Plants_load(path_to_plants_yaml string) map[string]PlantDefinition {
	plantsBytes, readErr := filemngr.ReadFileToBytes(path_to_plants_yaml)
	if readErr != nil {
		// Essential to server start
		panic(readErr)
	}
	var plants map[string]PlantDefinition
	err := yaml.Unmarshal(plantsBytes, &plants)
	if err != nil {
		log.Error.Fatalf("%v", err)
		// log.Error.Fatalf("%v", err.(*json.SyntaxError))
		// log.Error.Fatalf("%v", err.(*yaml.TypeError))
	}
	return plants
}