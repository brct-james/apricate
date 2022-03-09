// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

// Defines a route
type Route struct {
	UUID string `json:"uuid" binding:"required"`
	Name string `json:"name" binding:"required"`
	IsPortTravel bool `json:"is_port_travel" binding:"required"`
	Start Location `json:"start" binding:"required"`
	End Location `json:"end" binding:"required"`
	ArrivalTime int64 `json:"arrival_time" binding:"required"`
}