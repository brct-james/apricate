// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

import (
	"apricate/filemngr"
	"apricate/log"

	"gopkg.in/yaml.v3"
)

// Load produce list by unmarhsalling given yaml file
func Produce_load(path_to_produce_yaml string) map[string]string {
	produceBytes, readErr := filemngr.ReadFileToBytes(path_to_produce_yaml)
	if readErr != nil {
		// Essential to server start
		panic(readErr)
	}
	var rawProduce map[string]string
	err := yaml.Unmarshal(produceBytes, &rawProduce)
	if err != nil {
		log.Error.Fatalf("%v", err)
	}
	return rawProduce
}