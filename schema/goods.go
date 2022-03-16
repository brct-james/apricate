// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

import (
	"bytes"
	"encoding/json"
)

// Define a good
type Good struct {
	Name GoodType `yaml:"Name" json:"name" binding:"required"`
	Quality Quality `yaml:"Quality" json:"quality" binding:"required"`
	Quantity uint64 `yaml:"Quantity" json:"quantity" binding:"required"`
	Enchantment *Enchantment `yaml:"Enchantment" json:"enchantment,omitempty"`
}

// enum for good types
type GoodType uint8
const (
	Good_None GoodType = 0 // Default value ""
	Good_Water GoodType = 1
	Good_WildSeeds GoodType = 2
	Good_Cabbage GoodType = 3
	Good_CabbageSeeds GoodType = 4
	Good_ShelvisFig GoodType = 5
	Good_ShelvisFigSeeds GoodType = 6
	Good_ShelvisFigAle GoodType = 7
	Good_Potato GoodType = 8
	Good_Fertilizer GoodType = 9
)

var goodsToString = map[GoodType]string {
	Good_Water: "Water",
	Good_WildSeeds: "Wild Seeds",
	Good_Cabbage: "Cabbage",
	Good_CabbageSeeds: "Cabbage Seeds",
	Good_ShelvisFig: "Shelvis Fig",
	Good_ShelvisFigSeeds: "Shelvis Fig Seeds",
	Good_ShelvisFigAle: "Shelvis Fig Ale",
	Good_Potato: "Potato",
	Good_Fertilizer: "Fertilizer",
}

var goodsToID = map[string]GoodType {
	"Water": Good_Water,
	"Wild Seeds": Good_WildSeeds,
	"Cabbage": Good_Cabbage,
	"Cabbage Seeds": Good_CabbageSeeds,
	"Shelvis Fig": Good_ShelvisFig,
	"Shelvis Fig Seeds": Good_ShelvisFigSeeds,
	"Shelvis Fig Ale": Good_ShelvisFigAle,
	"Potato": Good_Potato,
	"Fertilizer": Good_Fertilizer,
}

func (s GoodType) String() string {
	return goodsToString[s]
}

// MarshalJSON marshals the enum as a text string
func (s GoodType) MarshalText() ([]byte, error) {
	buffer := bytes.NewBufferString(goodsToString[s])
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a text string to the enum value
func (s *GoodType) UnmarshalText(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = goodsToID[j]
	return nil
}

// MarshalJSON marshals the enum as a quoted json string
func (s GoodType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(goodsToString[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *GoodType) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = goodsToID[j]
	return nil
}

// MarshalYAML marshals the enum as a quoted yaml string
func (s GoodType) MarshalYAML() (interface{}, error) {
	buffer := bytes.NewBufferString(goodsToString[s])
	return buffer.Bytes(), nil
}

// UnmarshalYAML unmashals a quoted yaml string to the enum value
func (s *GoodType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var j string
	if err := unmarshal(&j); err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = goodsToID[j]
	return nil
}