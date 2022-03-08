// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

// Defines a ledger
type Ledger struct {
	Currencies map[string]uint64 `json:"currencies" binding:"required"`
	Favor map[string]int8 `json:"favor" binding:"required"`
	Escrow map[string]uint64 `json:"escrow" binding:"required"`
}