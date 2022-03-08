// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

// Defines a region
type Region struct {
	Name string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	PortConnections []string `json:"port_connections" binding:"required"`
}