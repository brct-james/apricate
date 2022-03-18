// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

type MainDictionary struct {
	Goods map[string]interface{} `yaml:"Goods" json:"goods" binding:"required"`
	Seeds map[string]string`yaml:"Seeds" json:"seeds" binding:"required"`
	Plants map[string]PlantDefinition`yaml:"Plants" json:"plants" binding:"required"`
}