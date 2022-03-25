// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

import (
	"apricate/log"
	"apricate/rdb"
	"bytes"
	"encoding/json"
	"fmt"
)

// enum for assistant types
type AssistantTypes uint8
const (
	Imp AssistantTypes = 0
	Familiar AssistantTypes = 1
	Golem AssistantTypes = 2
	Oni AssistantTypes = 3
	Sprite AssistantTypes = 4
	Dragon AssistantTypes = 5
)

// define map of assistantType to base speed
// note: largest land distance is 283, more typical travel between 40-90 so travel time will be Math.Ceil(distance_factor * distance) in seconds
// distance factor = (1800 / Math.Sqrt(80000)) so max is 30 minutes from -100,-100 to 100,100 with speed 1
var aTypeToBaseSpeed = map[AssistantTypes]int {
	Imp: 5,
	Familiar: 7,
	Golem: 3,
	Oni: 1, // slowest base cap
	Sprite: 10, // fastest base cap
	Dragon: 6,
}

// define map of assistantType to base carry cap
var aTypeToBaseCarryCap = map[AssistantTypes]int {
	Imp: 16,
	Familiar: 8,
	Golem: 64,
	Oni: 256, // largest base cap
	Sprite: 1, // smallest base cap - enough to carry a letter and not much else
	Dragon: 128,
}

// Defines an assistant
type Assistant struct {
	UUID string `json:"uuid" binding:"required"`
	ID uint64 `json:"id" binding:"required"`
	Archetype AssistantTypes `json:"archetype" binding:"required"`
	Speed int `json:"speed" binding:"required"`
	CarryCap int `json:"carrying_capacity" binding:"required"`
	Improvements map[string]uint8 `json:"improvements" binding:"required"`
	Location string `json:"location" binding:"required"` // EITHER the location symbol OR the caravan UUID
}

func NewAssistant(username string, countOfUserAssistants uint64, archetype AssistantTypes, locationSymbol string) *Assistant {
	return &Assistant{
		UUID: username + "|Assistant-" + fmt.Sprintf("%d", countOfUserAssistants),
		ID: countOfUserAssistants,
		Archetype: archetype,
		Speed: aTypeToBaseSpeed[archetype],
		CarryCap: aTypeToBaseCarryCap[archetype],
		Improvements: make(map[string]uint8),
		Location: locationSymbol,
	}
}

// Check DB for existing assistant with given uuid and return bool for if exists, and error if error encountered
func CheckForExistingAssistant (uuid string, tdb rdb.Database) (bool, error) {
	// Get assistant
	_, getError := tdb.GetJsonData(uuid, ".")
	if getError != nil {
		if fmt.Sprint(getError) != "redis: nil" {
			// error
			return false, getError
		}
		// assistant not found
		return false, nil
	}
	// Got successfully
	return true, nil
}

// Get assistant from DB, bool is assistant found
func GetAssistantFromDB (uuid string, tdb rdb.Database) (Assistant, bool, error) {
	// Get assistant json
	someJson, getError := tdb.GetJsonData(uuid, ".")
	if getError != nil {
		if fmt.Sprint(getError) != "redis: nil" {
			// assistant not found
			return Assistant{}, false, nil
		}
		// error
		return Assistant{}, false, getError
	}
	// Got successfully, unmarshal
	someData := Assistant{}
	unmarshalErr := json.Unmarshal(someJson, &someData)
	if unmarshalErr != nil {
		log.Error.Fatalf("Could not unmarshal assistant json from DB: %v", unmarshalErr)
		return Assistant{}, false, unmarshalErr
	}
	return someData, true, nil
}

// Get assistant from DB, bool is assistant found
func GetAssistantsFromDB (uuids []string, tdb rdb.Database) (map[string]Assistant, bool, error) {
	// Get assistant json
	someJson, getError := tdb.MGetJsonData(".", uuids)
	if getError != nil {
		if fmt.Sprint(getError) != "redis: nil" {
			// assistant not found
			return map[string]Assistant{}, false, nil
		}
		// error
		return map[string]Assistant{}, false, getError
	}
	// Got successfully, unmarshal
	someData := make(map[string]Assistant, len(someJson))
	for _, tempjson := range someJson {
		data := Assistant{}
		unmarshalErr := json.Unmarshal(tempjson, &data)
		if unmarshalErr != nil {
			log.Error.Fatalf("Could not unmarshal assistant json from DB: %v", unmarshalErr)
			return map[string]Assistant{}, false, unmarshalErr
		}
		someData[data.UUID] = data
	}
	
	return someData, true, nil
}

// Get assistantdata at path from DB, bool is assistant found
func GetAssistantDataAtPathFromDB (uuid string, path string, tdb rdb.Database) (interface{}, bool, error) {
	// Get assistant json
	someJson, getError := tdb.GetJsonData(uuid, path)
	if getError != nil {
		if fmt.Sprint(getError) != "redis: nil" {
			// assistant not found
			return nil, false, nil
		}
		// error
		return nil, false, getError
	}
	// Got successfully, unmarshal
	var someData interface{}
	unmarshalErr := json.Unmarshal(someJson, &someData)
	if unmarshalErr != nil {
		log.Error.Fatalf("Could not unmarshal assistant json from DB: %v", unmarshalErr)
		return nil, false, unmarshalErr
	}
	return someData, true, nil
}

// Attempt to save assistant, returns error or nil if successful
func SaveAssistantToDB(tdb rdb.Database, assistantData *Assistant) error {
	log.Debug.Printf("Saving assistant %s to DB", assistantData.UUID)
	err := tdb.SetJsonData(assistantData.UUID, ".", assistantData)
	// creationSuccess := rdb.CreateAssistant(tdb, assistantname, uuid, 0)
	return err
}

// Attempt to save assistant data at path, returns error or nil if successful
func SaveAssistantDataAtPathToDB(tdb rdb.Database, uuid string, path string, newValue interface{}) error {
	log.Debug.Printf("Saving assistant data at path %s to DB for uuid %s", path, uuid)
	err := tdb.SetJsonData(uuid, path, newValue)
	return err
}

func (s AssistantTypes) String() string {
	return assistantTypesToString[s]
}

var assistantTypesToString = map[AssistantTypes]string {
	Imp: "Imp",
	Familiar: "Familiar",
	Golem: "Golem",
	Oni: "Oni",
	Sprite: "Sprite",
	Dragon: "Dragon",
}

var assistantTypesToID = map[string]AssistantTypes {
	"Imp": Imp,
	"Familiar": Familiar,
	"Golem": Golem,
	"Oni": Oni,
	"Sprite": Sprite,
	"Dragon": Dragon,
}

// MarshalJSON marshals the enum as a quoted json string
func (s AssistantTypes) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(assistantTypesToString[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *AssistantTypes) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = assistantTypesToID[j]
	return nil
}