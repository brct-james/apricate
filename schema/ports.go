// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

// Defines a port
type Port struct {
	Name string `yaml:"Name" json:"name" binding:"required"`
	Symbol string `yaml:"Symbol" json:"symbol" binding:"required"`
	Connection string `yaml:"Connection" json:"connection" binding:"required"`
	ConnectedLocation string `yaml:"ConnectedLocation" json:"connected_locations" binding:"required"`
	Fare uint64 `yaml:"Fare" json:"fare" binding:"required"`
	Duration int `yaml:"Duration" json:"duration" binding:"required"`
}