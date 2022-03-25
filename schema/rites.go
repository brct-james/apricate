// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

import (
	"apricate/filemngr"
	"apricate/log"

	"gopkg.in/yaml.v3"
)

// Define rite dictionary entry
type Rite struct {
	RunicSymbol string `yaml:"RunicSymbol" json:"runic_symbol" binding:"required"`
	Name string `yaml:"Name" json:"name" binding:"required"`
	Description string `yaml:"Description" json:"description" binding:"required"` // describe effects
	MinimumDistortion int `yaml:"MinimumDistortion" json:"minimum_distortion_tier" binding:"required"`
	MaximumDistortion int `yaml:"MaximumDistortion" json:"maximum_distortion_tier" binding:"required"`
	ArcaneFlux int `yaml:"ArcaneFlux" json:"arcane_flux" binding:"required"`
	Currencies map[string]uint64 `yaml:"Currencies" json:"currencies" binding:"required"`
	Materials Wareset `yaml:"Materials" json:"materials" binding:"required"`
}

// Load rite struct by unmarhsalling given yaml file
func Rites_load(path_to_rites_yaml string) map[string]Rite {
	ritesBytes, readErr := filemngr.ReadFileToBytes(path_to_rites_yaml)
	if readErr != nil {
		// Essential to server start
		panic(readErr)
	}
	var rites map[string]Rite
	err := yaml.Unmarshal(ritesBytes, &rites)
	if err != nil {
		log.Error.Fatalln(err)
	}
	
	return rites
}