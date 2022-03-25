// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

import (
	"apricate/filemngr"
	"apricate/log"
	"math"

	"gopkg.in/yaml.v3"
)

// Defines a location
type Location struct {
	Name string `yaml:"Name" json:"name" binding:"required"`
	Symbol string `yaml:"Symbol" json:"symbol" binding:"required"`
	IslandName string `yaml:"IslandName" json:"island_name" binding:"required"`
	X int8 `yaml:"X" json:"x" binding:"required"` //-100:100
	Y int8 `yaml:"Y" json:"y" binding:"required"` //-100:100
	Description string `yaml:"Description" json:"description" binding:"required"`
	NPCs []string `yaml:"NPCs" json:"npcs" binding:"required"`
}

// Calculate travel time for a caravan between two locations
// Return validation map incase not same island and not ports, and travel time in seconds
func CalculateTravelTime(world World, a string, b string, slowestSpeed int) (map[string]string, int) {
	log.Debug.Printf("CalculateTravelTime a: %s b: %s, slowestSpeed: %d", a, b, slowestSpeed)
	validationMap := make(map[string]string)
	aLoc, aOk := world.Locations[a]
	bLoc, bOk := world.Locations[b]
	if !aOk {
		validationMap["origin"] = "Could not map origin to known location, ensure matches expected structure like example TS-PR-HF"
	}
	if !bOk {
		validationMap["destination"] = "Could not map destination to known location, ensure matches expected structure like example TS-PR-HF"
	}
	var travelTime int
	if aLoc.IslandName != bLoc.IslandName {
		// TODO check if there is port connection between islands, and if the origin and destination are the connected ports, rather than simply erroring
		travelTime = 0
		validationMap["destination"] = "Cannot travel between islands yet"
	} else {
		travelTime = int(math.Ceil(math.Sqrt(math.Pow(float64(bLoc.X - aLoc.X), 2) + math.Pow(float64(bLoc.Y - aLoc.Y), 2)))) * 10 / slowestSpeed
	}
	return validationMap, travelTime
}

// Load location struct by unmarhsalling given yaml file
func Locations_load(path_to_locations_yaml string) map[string]Location {
	locationsBytes, readErr := filemngr.ReadFilesToBytes(path_to_locations_yaml)
	if readErr != nil {
		// Essential to server start
		panic(readErr)
	}
	locations := make(map[string]Location)
	for _, byte := range locationsBytes {
		var location map[string]Location
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