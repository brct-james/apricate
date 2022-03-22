// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

import (
	"apricate/filemngr"
	"apricate/log"

	"gopkg.in/yaml.v3"
)

// Defines a island
type Island struct {
	Name string `yaml:"Name" json:"name" binding:"required"`
	Description string `yaml:"Description" json:"description" binding:"required"`
	Ports map[string]Port `yaml:"Ports" json:"ports" binding:"required"`
}

// Load island struct by unmarhsalling given yaml file
func Islands_load(path_to_islands_yaml string) map[string]Island {
	islandsBytes, readErr := filemngr.ReadFilesToBytes(path_to_islands_yaml)
	if readErr != nil {
		// Essential to server start
		panic(readErr)
	}
	islands := make(map[string]Island)
	for _, byte := range islandsBytes {
		var island map[string]Island
		err := yaml.Unmarshal(byte, &island)
		if err != nil {
			log.Error.Fatalln(err)
		}
		for k, v := range island {
			islands[k] = v
		}
	}
	return islands
}

// Load location struct by unmarhsalling given yaml file
// func Locations_load(path_to_locations_yaml string) map[string]map[string]Location {
// 	locationsBytes := filemngr.ReadFilesToBytes(path_to_locations_yaml)
// 	locations := make(map[string]map[string]Location)
// 	for _, byte := range locationsBytes {
// 		var location map[string]map[string]Location
// 		err := yaml.Unmarshal(byte, &location)
// 		if err != nil {
// 			log.Error.Fatalln(err)
// 		}
// 		for k, v := range location {
// 			locations[k] = v
// 		}
// 	}
// 	return locations
// }