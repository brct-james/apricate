// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

import (
	"apricate/filemngr"
	"apricate/log"

	"gopkg.in/yaml.v3"
)

// Defines a region
type Region struct {
	Identifier
	RegionGroup string `yaml:"RegionGroup" json:"region_group" binding:"required"`
	Description string `yaml:"Description" json:"description" binding:"required"`
	Islands []Identifier `yaml:"Islands" json:"islands" binding:"required"`
}

// Load region struct by unmarhsalling given yaml file
func Regions_load(path_to_regions_yaml string) map[string]Region {
	regionsBytes := filemngr.ReadFileToBytes(path_to_regions_yaml)
	var regions map[string]Region
	err := yaml.Unmarshal(regionsBytes, &regions)
	if err != nil {
		log.Error.Fatalln(err)
	}
	return regions
}