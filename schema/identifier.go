// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

// Defines an Identifier
type Identifier struct {
	Name string `yaml:"Name" json:"name" binding:"required"`
	Symbol string `yaml:"Symbol" json:"symbol" binding:"required"`
}