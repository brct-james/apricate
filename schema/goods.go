// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

import (
	"bytes"
	"encoding/json"
)

// enum for good types
type Goods uint8
const (
	Good_Water Goods = 0
	Good_Cabbage Goods = 1
	Good_CabbageSeeds Goods = 2
	Good_Shelvis Goods = 3
	Good_ShelvisSeeds Goods = 4
	Good_Potato Goods = 5
)

var goodsToString = map[Goods]string {
	Good_Water: "Water",
	Good_Cabbage: "Cabbage",
	Good_CabbageSeeds: "Cabbage Seeds",
	Good_Shelvis: "Shelvis",
	Good_ShelvisSeeds: "Shelvis Seeds",
	Good_Potato: "Potato",
}

var goodsToID = map[string]Goods {
	"Water": Good_Water,
	"Cabbage": Good_Cabbage,
	"Cabbage Seeds": Good_CabbageSeeds,
	"Shelvis": Good_Shelvis,
	"Shelvis Seeds": Good_ShelvisSeeds,
	"Potato": Good_Potato,
}

func (s Goods) String() string {
	return goodsToString[s]
}

// MarshalJSON marshals the enum as a text string
func (s Goods) MarshalText() ([]byte, error) {
	buffer := bytes.NewBufferString(goodsToString[s])
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a text string to the enum value
func (s *Goods) UnmarshalText(b []byte) error {
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
func (s Goods) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(goodsToString[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *Goods) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = goodsToID[j]
	return nil
}