// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

// Defines a port
type Port struct {
	Name string `yaml:"Name" json:"name" binding:"required"`
	Fare uint64 `yaml:"Fare" json:"fare" binding:"required"`
	Duration uint16 `yaml:"Duration" json:"duration" binding:"required"`
}