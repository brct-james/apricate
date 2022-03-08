// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

// enum for api response codes
type Achievement int16
const (
	Achievement_Owner Achievement = -1
	Achievement_Contributor Achievement = 0
	Achievement_Noob Achievement = 1
)