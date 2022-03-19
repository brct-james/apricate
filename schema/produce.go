// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

import (
	"apricate/filemngr"
	"apricate/log"

	"gopkg.in/yaml.v3"
)

// Define produce
type Produce struct {
	Good `yaml:",inline"`
	Size Size `yaml:"Size" json:"size" binding:"required"`
}

// New produce
func NewProduce(name string, size Size, quantity uint64) *Produce {
	return &Produce{
		Good: Good{
			Name: name,
			Quantity: quantity,
		},
		Size: size,
	}
}

// Load produce list by unmarhsalling given yaml file
func Produce_load(path_to_produce_yaml string) map[string]string {
	produceBytes := filemngr.ReadFileToBytes(path_to_produce_yaml)
	var rawProduce map[string]string
	err := yaml.Unmarshal(produceBytes, &rawProduce)
	if err != nil {
		log.Error.Fatalf("%v", err)
	}
	return rawProduce
}