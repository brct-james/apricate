// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

// Defines a location
type Location struct {
	Name string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	RegionName string `json:"region_name" binding:"required"`
	X uint8 `json:"x" binding:"required"`
	Y uint8 `json:"y" binding:"required"`
	NPCs []string `json:"npcs" binding:"required"`
}