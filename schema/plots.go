// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

import (
	"apricate/log"
	"apricate/rdb"
	"apricate/uuid"
	"encoding/json"
	"fmt"
)

// Defines a plot
type Plot struct {
	UUID string `json:"uuid" binding:"required"`
	Capacity Size `json:"capacity" binding:"required"`
	Plants []Plant `json:"plants" binding:"required"`
}

func NewPlot(capacity Size) *Plot {
	uuid := uuid.NewUUID()
	return &Plot{
		UUID: uuid,
		Capacity: capacity,
		Plants: make([]Plant, 0),
	}
}

// Check DB for existing plot with given uuid and return bool for if exists, and error if error encountered
func CheckForExistingPlot (uuid string, tdb rdb.Database) (bool, error) {
	// Get plot
	_, getError := tdb.GetJsonData(uuid, ".")
	if getError != nil {
		if fmt.Sprint(getError) != "redis: nil" {
			// error
			return false, getError
		}
		// plot not found
		return false, nil
	}
	// Got successfully
	return true, nil
}

// Get plot from DB, bool is plot found
func GetPlotFromDB (uuid string, tdb rdb.Database) (Plot, bool, error) {
	// Get plot json
	someJson, getError := tdb.GetJsonData(uuid, ".")
	if getError != nil {
		if fmt.Sprint(getError) != "redis: nil" {
			// plot not found
			return Plot{}, false, nil
		}
		// error
		return Plot{}, false, getError
	}
	// Got successfully, unmarshal
	someData := Plot{}
	unmarshalErr := json.Unmarshal(someJson, &someData)
	if unmarshalErr != nil {
		log.Error.Fatalf("Could not unmarshal plot json from DB: %v", unmarshalErr)
		return Plot{}, false, unmarshalErr
	}
	return someData, true, nil
}

// Get plot from DB, bool is plot found
func GetPlotsFromDB (uuids []string, tdb rdb.Database) ([]Plot, bool, error) {
	// Get plot json
	someJson, getError := tdb.MGetJsonData(".", uuids)
	if getError != nil {
		if fmt.Sprint(getError) != "redis: nil" {
			// plot not found
			return []Plot{}, false, nil
		}
		// error
		return []Plot{}, false, getError
	}
	// Got successfully, unmarshal
	someData := make([]Plot, len(someJson))
	for i, tempjson := range someJson {
		data := Plot{}
		unmarshalErr := json.Unmarshal(tempjson, &data)
		if unmarshalErr != nil {
			log.Error.Fatalf("Could not unmarshal plot json from DB: %v", unmarshalErr)
			return []Plot{}, false, unmarshalErr
		}
		someData[i] = data
	}
	
	return someData, true, nil
}

// Get plotdata at path from DB, bool is plot found
func GetPlotDataAtPathFromDB (uuid string, path string, tdb rdb.Database) (interface{}, bool, error) {
	// Get plot json
	someJson, getError := tdb.GetJsonData(uuid, path)
	if getError != nil {
		if fmt.Sprint(getError) != "redis: nil" {
			// plot not found
			return nil, false, nil
		}
		// error
		return nil, false, getError
	}
	// Got successfully, unmarshal
	var someData interface{}
	unmarshalErr := json.Unmarshal(someJson, &someData)
	if unmarshalErr != nil {
		log.Error.Fatalf("Could not unmarshal plot json from DB: %v", unmarshalErr)
		return nil, false, unmarshalErr
	}
	return someData, true, nil
}

// Attempt to save plot, returns error or nil if successful
func SavePlotToDB(tdb rdb.Database, plotData *Plot) error {
	log.Debug.Printf("Saving plot %s to DB", plotData.UUID)
	err := tdb.SetJsonData(plotData.UUID, ".", plotData)
	// creationSuccess := rdb.CreatePlot(tdb, plotname, uuid, 0)
	return err
}

// Attempt to save plot data at path, returns error or nil if successful
func SavePlotDataAtPathToDB(tdb rdb.Database, uuid string, path string, newValue interface{}) error {
	log.Debug.Printf("Saving plot data at path %s to DB for uuid %s", path, uuid)
	err := tdb.SetJsonData(uuid, path, newValue)
	return err
}