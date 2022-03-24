// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

import (
	"apricate/log"
	"apricate/rdb"
	"apricate/timecalc"
	"encoding/json"
	"fmt"
	"time"
)

// Define CaravanWares
type CaravanWares struct {
	Tools map[string]uint64 `json:"tools,omitempty"`
	Produce map[string]uint64 `json:"produce,omitempty"`
	Seeds map[string]uint64 `json:"seeds,omitempty"`
	Goods map[string]uint64 `json:"goods,omitempty"`
}

// Defines a caravan charter
type CaravanCharter struct {
	Origin string `json:"origin" binding:"required"`
	Destination string `json:"destination" binding:"required"`
	Assistants []string `json:"assistants" binding:"required"`
	Wares CaravanWares `json:"wares,omitempty"`
}

// Validate caravan charter, return validation map
func ValidateCaravanCharter(charter CaravanCharter) map[string]string {
	res := make(map[string]string)
	if len(charter.Origin) < 8 {
		res["origin"] = "Too Short, expect minimum 8 characters based on example TS-PR-HF"
	}
	if len(charter.Destination) < 8 {
		res["destination"] = "Too Short, expect minimum 8 characters based on example TS-PR-HF"
	}
	if len(charter.Assistants) < 1 {
		res["assistants"] = "Must specify at least one assistant ID to include in caravan"
	}
	return res
}

// Defines a caravan
type Caravan struct {
	UUID string `json:"uuid" binding:"required"`
	ID int64 `json:"id" binding:"required"`
	CaravanCharter `yaml:",inline"`
	ArrivalTime int64 `json:"arrival_time" binding:"required"`
	SecondsTillArrival int64 `json:"seconds_till_arrival" binding:"required"` // SHOULD BE STORED AS 0, ONLY FOR FORMATTING RESPONSE
}

func NewCaravan(UUID string, timestamp time.Time, origin string, destination string, assistants []string, wares CaravanWares, travelTimeSeconds int) *Caravan {
	return &Caravan{
		UUID: UUID,
		ID: timestamp.UnixNano(),
		CaravanCharter: CaravanCharter{
			Origin: origin,
			Destination: destination,
			Assistants: assistants,
			Wares: wares,
		},
		ArrivalTime: timecalc.AddSecondsToTimestamp(timestamp, travelTimeSeconds).Unix(),
		SecondsTillArrival: int64(travelTimeSeconds),
	}
}

// Check DB for existing caravan with given uuid and return bool for if exists, and error if error encountered
func CheckForExistingCaravan (uuid string, tdb rdb.Database) (bool, error) {
	// Get caravan
	_, getError := tdb.GetJsonData(uuid, ".")
	if getError != nil {
		if fmt.Sprint(getError) != "redis: nil" {
			// error
			return false, getError
		}
		// caravan not found
		return false, nil
	}
	// Got successfully
	return true, nil
}

// Get caravan from DB, bool is caravan found
func GetCaravanFromDB (uuid string, tdb rdb.Database) (Caravan, bool, error) {
	// Get caravan json
	someJson, getError := tdb.GetJsonData(uuid, ".")
	if getError != nil {
		if fmt.Sprint(getError) != "redis: nil" {
			// caravan not found
			return Caravan{}, false, nil
		}
		// error
		return Caravan{}, false, getError
	}
	// Got successfully, unmarshal
	someData := Caravan{}
	unmarshalErr := json.Unmarshal(someJson, &someData)
	if unmarshalErr != nil {
		log.Error.Fatalf("Could not unmarshal caravan json from DB: %v", unmarshalErr)
		return Caravan{}, false, unmarshalErr
	}
	return someData, true, nil
}

// Get caravan from DB, bool is caravan found
func GetCaravansFromDB (uuids []string, tdb rdb.Database) ([]Caravan, bool, error) {
	// Get caravan json
	someJson, getError := tdb.MGetJsonData(".", uuids)
	if getError != nil {
		if fmt.Sprint(getError) != "redis: nil" {
			// caravan not found
			return []Caravan{}, false, nil
		}
		// error
		return []Caravan{}, false, getError
	}
	// Got successfully, unmarshal
	someData := make([]Caravan, len(someJson))
	for i, tempjson := range someJson {
		data := Caravan{}
		unmarshalErr := json.Unmarshal(tempjson, &data)
		if unmarshalErr != nil {
			log.Error.Fatalf("Could not unmarshal caravan json from DB: %v", unmarshalErr)
			return []Caravan{}, false, unmarshalErr
		}
		someData[i] = data
	}
	
	return someData, true, nil
}

// Get caravandata at path from DB, bool is caravan found
func GetCaravanDataAtPathFromDB (uuid string, path string, tdb rdb.Database) (interface{}, bool, error) {
	// Get caravan json
	someJson, getError := tdb.GetJsonData(uuid, path)
	if getError != nil {
		if fmt.Sprint(getError) != "redis: nil" {
			// caravan not found
			return nil, false, nil
		}
		// error
		return nil, false, getError
	}
	// Got successfully, unmarshal
	var someData interface{}
	unmarshalErr := json.Unmarshal(someJson, &someData)
	if unmarshalErr != nil {
		log.Error.Fatalf("Could not unmarshal caravan json from DB: %v", unmarshalErr)
		return nil, false, unmarshalErr
	}
	return someData, true, nil
}

// Attempt to save caravan, returns error or nil if successful
func SaveCaravanToDB(tdb rdb.Database, caravanData *Caravan) error {
	log.Debug.Printf("Saving caravan %s to DB", caravanData.UUID)
	err := tdb.SetJsonData(caravanData.UUID, ".", caravanData)
	// creationSuccess := rdb.CreateCaravan(tdb, caravanname, uuid, 0)
	return err
}

// Attempt to save caravan data at path, returns error or nil if successful
func SaveCaravanDataAtPathToDB(tdb rdb.Database, uuid string, path string, newValue interface{}) error {
	log.Debug.Printf("Saving caravan data at path %s to DB for uuid %s", path, uuid)
	err := tdb.SetJsonData(uuid, path, newValue)
	return err
}