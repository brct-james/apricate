// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

import (
	"apricate/filemngr"
	"apricate/log"

	"gopkg.in/yaml.v3"
)

// PLACEHOLDER to get basic market function into early alpha build - will revisit later

// Define a market
type Market struct {
	Name string `yaml:"Name" json:"name" binding:"required"`
	LocationSymbol string `yaml:"Location" json:"location_symbol" binding:"required"`
	Imports MarketIOField `yaml:"Imports" json:"imports" binding:"required"`
	Exports MarketIOField `yaml:"Exports" json:"exports" binding:"required"`
}

// Define a market import or export field
type MarketIOField struct {
	Produce map[string]uint64 `yaml:"Produce" json:"produce,omitempty"`
	Seeds map[string]uint64 `yaml:"Seeds" json:"seeds,omitempty"`
	Goods map[string]uint64 `yaml:"Goods" json:"goods,omitempty"`
	Tools map[ToolTypes]uint64 `yaml:"Tools" json:"tools,omitempty"`
}

// Load market struct by unmarhsalling given yaml file
func Markets_load(path_to_markets_yaml string) map[string]Market {
	marketsBytes := filemngr.ReadFileToBytes(path_to_markets_yaml)
	var markets map[string]Market
	err := yaml.Unmarshal(marketsBytes, &markets)
	if err != nil {
		log.Error.Fatalln(err)
	}
	
	return markets
}