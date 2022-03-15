// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

import (
	"apricate/log"
	"apricate/rdb"
	"apricate/uuid"
	"bytes"
	"encoding/json"
	"fmt"
)

// enum for contract types
type ContractTypes uint8
const (
	ContractType_Collect ContractTypes = 0
	ContractType_Deliver ContractTypes = 1
	ContractType_Courier ContractTypes = 2
	ContractType_Talk ContractTypes = 3
)

// Defines a contract
type Contract struct {
	UUID string `json:"uuid" binding:"required"`
	ContractType ContractTypes `json:"type" binding:"required"`
	RegionLocation string `json:"region_location" binding:"required"` // Location format: Region|Location
	NPC string `json:"NPC" binding:"required"`
	Terms []ContractTerms `json:"terms" binding:"required"`
	Reward []ContractReward `json:"reward" binding:"required"`
}

// Defines ContractTerms
type ContractTerms struct {
	NPC string `json:"npc,omitempty"`
	Item string `json:"item,omitempty"`
	Quantity uint64 `json:"quantity,omitempty"`
}

// Defines contract reward types
type RewardType uint8
const (
	RewardType_Currency RewardType = 0
	RewardType_Item RewardType = 1
)

// Defines ContractReward
type ContractReward struct {
	RewardType RewardType `json:"type" binding:"required"` 
	Item string `json:"item" binding:"required"`
	Quantity uint64 `json:"quantity" binding:"required"`
}

func NewContract(regionLocation string, contractType ContractTypes, npc string, terms []ContractTerms, reward []ContractReward) *Contract {
	uuid := uuid.NewUUID()
	return &Contract{
		UUID: uuid,
		ContractType: contractType,
		RegionLocation: regionLocation,
		NPC: npc,
		Terms: terms,
		Reward: reward,
	}
}

// Check DB for existing contract with given uuid and return bool for if exists, and error if error encountered
func CheckForExistingContract (uuid string, tdb rdb.Database) (bool, error) {
	// Get contract
	_, getError := tdb.GetJsonData(uuid, ".")
	if getError != nil {
		if fmt.Sprint(getError) != "redis: nil" {
			// error
			return false, getError
		}
		// contract not found
		return false, nil
	}
	// Got successfully
	return true, nil
}

// Get contract from DB, bool is contract found
func GetContractFromDB (uuid string, tdb rdb.Database) (Contract, bool, error) {
	// Get contract json
	someJson, getError := tdb.GetJsonData(uuid, ".")
	if getError != nil {
		if fmt.Sprint(getError) != "redis: nil" {
			// contract not found
			return Contract{}, false, nil
		}
		// error
		return Contract{}, false, getError
	}
	// Got successfully, unmarshal
	someData := Contract{}
	unmarshalErr := json.Unmarshal(someJson, &someData)
	if unmarshalErr != nil {
		log.Error.Fatalf("Could not unmarshal contract json from DB: %v", unmarshalErr)
		return Contract{}, false, unmarshalErr
	}
	return someData, true, nil
}

// Get contract from DB, bool is contract found
func GetContractsFromDB (uuids []string, tdb rdb.Database) ([]Contract, bool, error) {
	// Get contract json
	someJson, getError := tdb.MGetJsonData(".", uuids)
	if getError != nil {
		if fmt.Sprint(getError) != "redis: nil" {
			// contract not found
			return []Contract{}, false, nil
		}
		// error
		return []Contract{}, false, getError
	}
	// Got successfully, unmarshal
	someData := make([]Contract, len(someJson))
	for i, tempjson := range someJson {
		data := Contract{}
		unmarshalErr := json.Unmarshal(tempjson, &data)
		if unmarshalErr != nil {
			log.Error.Fatalf("Could not unmarshal contract json from DB: %v", unmarshalErr)
			return []Contract{}, false, unmarshalErr
		}
		someData[i] = data
	}
	
	return someData, true, nil
}

// Get contractdata at path from DB, bool is contract found
func GetContractDataAtPathFromDB (uuid string, path string, tdb rdb.Database) (interface{}, bool, error) {
	// Get contract json
	someJson, getError := tdb.GetJsonData(uuid, path)
	if getError != nil {
		if fmt.Sprint(getError) != "redis: nil" {
			// contract not found
			return nil, false, nil
		}
		// error
		return nil, false, getError
	}
	// Got successfully, unmarshal
	var someData interface{}
	unmarshalErr := json.Unmarshal(someJson, &someData)
	if unmarshalErr != nil {
		log.Error.Fatalf("Could not unmarshal contract json from DB: %v", unmarshalErr)
		return nil, false, unmarshalErr
	}
	return someData, true, nil
}

// Attempt to save contract, returns error or nil if successful
func SaveContractToDB(tdb rdb.Database, contractData *Contract) error {
	log.Debug.Printf("Saving contract %s to DB", contractData.UUID)
	err := tdb.SetJsonData(contractData.UUID, ".", contractData)
	// creationSuccess := rdb.CreateContract(tdb, contractname, uuid, 0)
	return err
}

// Attempt to save contract data at path, returns error or nil if successful
func SaveContractDataAtPathToDB(tdb rdb.Database, uuid string, path string, newValue interface{}) error {
	log.Debug.Printf("Saving contract data at path %s to DB for uuid %s", path, uuid)
	err := tdb.SetJsonData(uuid, path, newValue)
	return err
}

func (s ContractTypes) String() string {
	return contractTypesToString[s]
}

var contractTypesToString = map[ContractTypes]string {
	ContractType_Collect: "Collect",
	ContractType_Deliver: "Deliver",
	ContractType_Courier: "Courier",
	ContractType_Talk: "Talk",
}

var contractTypesToID = map[string]ContractTypes {
	"Collect": ContractType_Collect,
	"Deliver": ContractType_Deliver,
	"Courier": ContractType_Courier,
	"Talk": ContractType_Talk,
}

// MarshalJSON marshals the enum as a quoted json string
func (s ContractTypes) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(contractTypesToString[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *ContractTypes) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = contractTypesToID[j]
	return nil
}

func (s RewardType) String() string {
	return rewardToString[s]
}

var rewardToString = map[RewardType]string {
	RewardType_Currency: "Currency",
	RewardType_Item: "Item",
}

var rewardToID = map[string]RewardType {
	"Currency": RewardType_Currency,
	"Item": RewardType_Item,
}

// MarshalJSON marshals the enum as a quoted json string
func (s RewardType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(rewardToString[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *RewardType) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = rewardToID[j]
	return nil
}