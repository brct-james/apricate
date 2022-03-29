// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

import (
	"apricate/filemngr"
	"apricate/log"
	"fmt"
	"math"
	"strings"

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
// Return validation map incase not same island and not ports, and travel time in seconds, and fare cost in coins
func CalculateTravelTime(world World, a string, b string, slowestSpeed int) (map[string]string, int, uint64) {
	log.Debug.Printf("CalculateTravelTime a: %s b: %s, slowestSpeed: %d", a, b, slowestSpeed)
	validationMap := make(map[string]string)
	// Ensure Symbols conforms to expectations
	originSymbolSlice := strings.Split(a, "-")
	if len(originSymbolSlice) < 3 {
		// Origin Symbol not 3-part, invalid
		validationMap["origin"] = "Origin symbol must be 3-part, like TS-PR-HF"
	}
	destSymbolSlice := strings.Split(a, "-")
	if len(destSymbolSlice) < 3 {
		// Dest Symbol not 3-part, invalid
		validationMap["destination"] = "Destination symbol must be 3-part, like TS-PR-HF"
	}
	if len(validationMap) > 0 {
		// Return early if found validation issues
		return validationMap, 0, 0
	}
	originIslandSymbol := strings.Join(originSymbolSlice[0:2], "-")
	originIsland, oiOK := world.Islands[originIslandSymbol]
	if !oiOK {
		// Island symbol not found in world dictionary
		validationMap["origin"] = "Origin symbol island component not found in world dict, ensure matches expected structure like example TS-PR-HF and island abbreviation (PR) is correct"
	}
	destIslandSymbol := strings.Join(destSymbolSlice[0:2], "-")
	_, disOK := world.Islands[destIslandSymbol]
	if !disOK {
		// Island symbol not found in world dictionary
		validationMap["destination"] = "Destination symbol island component not found in world dict, ensure matches expected structure like example TS-PR-HF and island abbreviation (PR) is correct"
	}
	if len(validationMap) > 0 {
		// Return early if found validation issues
		return validationMap, 0, 0
	}
	// Get Locations Info from World
	aLoc, aOk := world.Locations[a]
	bLoc, bOk := world.Locations[b]
	if !aOk {
		validationMap["origin"] = "Could not map origin to known location, ensure matches expected structure like example TS-PR-HF and location exists"
	}
	if !bOk {
		validationMap["destination"] = "Could not map destination to known location, ensure matches expected structure like example TS-PR-HF and location exists"
	}
	if len(validationMap) > 0 {
		// Return early if found validation issues
		return validationMap, 0, 0
	}
	var travelTime int
	if aLoc.IslandName != bLoc.IslandName {
		// Validate origin is a port in the originIsland definition
		portListString := "["
		for _, port := range originIsland.Ports {
			if port.Symbol == a {
				// Origin Found in Port List, Check that Dest is ConnectedLocation
				if port.ConnectedLocation == b {
					// Dest is ConnectedLocation, return travel time and fare
					return validationMap, port.Duration, port.Fare
				}
				// Dest Not connected location
				validationMap["destination"] = fmt.Sprintf("Origin and Destination are on different islands. Origin was found to be a port, but destination was not the connected location. Did you mean %s?", port.ConnectedLocation)
				return validationMap, 0, 0
			}
			if len(portListString) == 1 {
				portListString += port.Symbol
			} else {
				portListString += ", " + port.Symbol
			}
		}
		// Origin not found in island ports list
		portListString += "]"
		validationMap["origin"] = fmt.Sprintf("Origin and Destination are on different islands. Origin was not found to be a port. Ports at this island are from the list: %s", portListString)
		travelTime = 0
	} else {
		travelTime = int(math.Ceil(math.Sqrt(math.Pow(float64(bLoc.X - aLoc.X), 2) + math.Pow(float64(bLoc.Y - aLoc.Y), 2)))) * 10 / slowestSpeed
	}
	return validationMap, travelTime, 0
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