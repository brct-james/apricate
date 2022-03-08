// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

// Defines a npc
type NPC struct {
	Name string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Location Location `json:"location" binding:"required"`
	// TODO: Add AvailableContracts and figure out how to restrict access depending on player favor and pre-requisite contracts/items_owned/currency_owned
}