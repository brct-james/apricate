// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

// Defines the world
type World struct {
	Name string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Regions map[string]Region `json:"regions" binding:"required"`
	Locations map[string]map[string]Location `json:"locations" binding:"required"`
}

// Load world struct by unmarhsalling given yaml file
func World_load(path_to_regions_yaml string, path_to_locations_directory string) World {
	regions := Regions_load(path_to_regions_yaml)
	locations := Locations_load(path_to_locations_directory)
	return World{
		Name: "Astrid",
		Description: "A fantasy world, torn apart by magical warfare. Continents reduced to islands, and oceans with few navigable routes due to residual magic storms.",
		Regions: regions,
		Locations: locations,
	}
}