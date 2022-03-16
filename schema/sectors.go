// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

import (
	"apricate/filemngr"
	"apricate/log"

	"gopkg.in/yaml.v3"
)

// Defines a sector
type Sector struct {
	Identifier
	SectorGroup string `yaml:"SectorGroup" json:"sector_group" binding:"required"`
	Description string `yaml:"Description" json:"description" binding:"required"`
	Islands []Identifier `yaml:"Islands" json:"islands" binding:"required"`
}

// Load sector struct by unmarhsalling given yaml file
func Sectors_load(path_to_sectors_yaml string) map[string]Sector {
	sectorsBytes := filemngr.ReadFileToBytes(path_to_sectors_yaml)
	var sectors map[string]Sector
	err := yaml.Unmarshal(sectorsBytes, &sectors)
	if err != nil {
		log.Error.Fatalln(err)
	}
	return sectors
}