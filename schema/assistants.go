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
	Hireling AssistantTypes = 0
	Familiar AssistantTypes = 1
	Golem AssistantTypes = 2
)

// Defines an assistant
type Assistant struct {
	UUID string `json:"uuid" binding:"required"`
	Archetype AssistantTypes `json:"archetype" binding:"required"`
	Improvements map[string]uint8 `json:"improvements" binding:"required"`
	LocationSymbol string `json:"location_symbol" binding:"required"`
	Route string `json:"route" binding:"required"`
}

func NewAssistant(username string, countOfUserAssistants int16, archetype AssistantTypes, locationSymbol string) *Assistant {
	return &Assistant{
		UUID: username + "-Assistant-" + fmt.Sprintf("%d", countOfUserAssistants),
		Archetype: archetype,
		Improvements: make(map[string]uint8),
		LocationSymbol: locationSymbol,
		Route: "",
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
func GetAssistantsFromDB (uuids []string, tdb rdb.Database) ([]Assistant, bool, error) {
	// Get assistant json
	someJson, getError := tdb.MGetJsonData(".", uuids)
	if getError != nil {
		if fmt.Sprint(getError) != "redis: nil" {
			// assistant not found
			return []Assistant{}, false, nil
		}
		// error
		return []Assistant{}, false, getError
	}
	// Got successfully, unmarshal
	someData := make([]Assistant, len(someJson))
	for i, tempjson := range someJson {
		data := Assistant{}
		unmarshalErr := json.Unmarshal(tempjson, &data)
		if unmarshalErr != nil {
			log.Error.Fatalf("Could not unmarshal assistant json from DB: %v", unmarshalErr)
			return []Assistant{}, false, unmarshalErr
		}
		someData[i] = data
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
	Hireling: "Hireling",
	Familiar: "Familiar",
	Golem: "Golem",
}

var assistantTypesToID = map[string]AssistantTypes {
	"Hireling": Hireling,
	"Familiar": Familiar,
	"Golem": Golem,
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