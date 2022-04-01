// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

// Defines about info
type AboutInfo struct {
	Name string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Information map[string]interface{} `json:"information" binding:"required"`
}