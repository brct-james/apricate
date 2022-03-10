// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

import (
	"apricate/filemngr"
	"apricate/log"

	"gopkg.in/yaml.v2"
)

// Defines a location
type Location struct {
	Name string `yaml:"Name" json:"name" binding:"required"`
	RegionName string `yaml:"RegionName" json:"region_name" binding:"required"`
	X int8 `yaml:"X" json:"x" binding:"required"`
	Y int8 `yaml:"Y" json:"y" binding:"required"`
	Description string `yaml:"Description" json:"description" binding:"required"`
	NPCs []string `yaml:"NPCs" json:"npcs" binding:"required"`
}

// Load location struct by unmarhsalling given yaml file
func Locations_load(path_to_locations_yaml string) map[string]map[string]Location {
	locationsBytes := filemngr.ReadFilesToBytes(path_to_locations_yaml)
	locations := make(map[string]map[string]Location)
	for _, byte := range locationsBytes {
		var location map[string]map[string]Location
		err := yaml.Unmarshal(byte, &location)
		if err != nil {
			log.Error.Fatalln(err)
		}
		for k, v := range location {
			locations[k] = v
		}
	}
	return locations
}