// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

import (
	"bytes"
	"encoding/json"
)

// PLACEHOLDER to get basic market function into early alpha build - will revisit later

// Define a market order
type MarketOrder struct {
	OrderType OrderType `json:"order_type" binding:"required"`
	TXType TXType `json:"transaction_type" binding:"required"`
	ItemType ItemType `json:"item_type" binding:"required"`
	ItemName string `json:"item_name" binding:"required"`
	Quantity uint64 `json:"quantity" binding:"required"`
}

// enum for assistant types
type OrderType uint16
const (
	MARKET OrderType = 0
	// LIMIT OrderType = 1
	// STOP OrderType = 2
	// STOPLIMIT OrderType = 3
)

func (s OrderType) String() string {
	return orderTypeToString[s]
}

var orderTypeToString = map[OrderType]string {
	MARKET: "MARKET",
	// LIMIT: "LIMIT",
	// STOP: "STOP",
	// STOPLIMIT: "STOPLIMIT",
}

var orderTypeToID = map[string]OrderType {
	"MARKET": MARKET,
	// "LIMIT": LIMIT,
	// "STOP": STOP,
	// "STOPLIMIT": STOPLIMIT,
}

// MarshalJSON marshals the enum as a quoted json string
func (s OrderType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(orderTypeToString[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *OrderType) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = orderTypeToID[j]
	return nil
}

// enum for assistant types
type ItemType uint16
const (
	GOOD ItemType = 0
	SEED ItemType = 1
	PRODUCE ItemType = 2
	TOOL ItemType = 3
)

func (s ItemType) String() string {
	return itemTypeToString[s]
}

var itemTypeToString = map[ItemType]string {
	GOOD: "GOODS",
	SEED: "SEEDS",
	PRODUCE: "PRODUCE",
	TOOL: "TOOLS",
}

var itemTypeToID = map[string]ItemType {
	"GOODS": GOOD,
	"SEEDS": SEED,
	"PRODUCE": PRODUCE,
	"TOOLS": TOOL,
}

// MarshalJSON marshals the enum as a quoted json string
func (s ItemType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(itemTypeToString[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *ItemType) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = itemTypeToID[j]
	return nil
}

// enum for assistant types
type TXType uint16
const (
	BUY TXType = 0
	SELL TXType = 1
)

func (s TXType) String() string {
	return tXTypeToString[s]
}

var tXTypeToString = map[TXType]string {
	BUY: "BUY",
	SELL: "SELL",
}

var tXTypeToID = map[string]TXType {
	"BUY": BUY,
	"SELL": SELL,
}

// MarshalJSON marshals the enum as a quoted json string
func (s TXType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(tXTypeToString[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *TXType) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = tXTypeToID[j]
	return nil
}

