// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

import (
	"apricate/log"
	"apricate/rdb"
	"apricate/uuid"
	"encoding/json"
	"fmt"
)

// Defines a warehouse
type Warehouse struct {
	UUID string `json:"uuid" binding:"required"`
	RegionLocation string `json:"region_location" binding:"required"` // Location format: Region|Location
	Goods map[Goods]uint64 `json:"goods" binding:"required"`
}

func NewEmptyWarehouse(regionLocation string) *Warehouse {
	return NewWarehouse(regionLocation, make(map[Goods]uint64))
}

func NewWarehouse(regionLocation string, starting_goods map[Goods]uint64) *Warehouse {
	uuid := uuid.NewUUID()
	return &Warehouse{
		UUID: uuid,
		RegionLocation: regionLocation,
		Goods: starting_goods,
	}
}

// Check DB for existing warehouse with given uuid and return bool for if exists, and error if error encountered
func CheckForExistingWarehouse (uuid string, tdb rdb.Database) (bool, error) {
	// Get warehouse
	_, getError := tdb.GetJsonData(uuid, ".")
	if getError != nil {
		if fmt.Sprint(getError) != "redis: nil" {
			// error
			return false, getError
		}
		// warehouse not found
		return false, nil
	}
	// Got successfully
	return true, nil
}

// Get warehouse from DB, bool is warehouse found
func GetWarehouseFromDB (uuid string, tdb rdb.Database) (Warehouse, bool, error) {
	// Get warehouse json
	someJson, getError := tdb.GetJsonData(uuid, ".")
	if getError != nil {
		if fmt.Sprint(getError) != "redis: nil" {
			// warehouse not found
			return Warehouse{}, false, nil
		}
		// error
		return Warehouse{}, false, getError
	}
	// Got successfully, unmarshal
	someData := Warehouse{}
	unmarshalErr := json.Unmarshal(someJson, &someData)
	if unmarshalErr != nil {
		log.Error.Fatalf("Could not unmarshal warehouse json from DB: %v", unmarshalErr)
		return Warehouse{}, false, unmarshalErr
	}
	return someData, true, nil
}

// Get warehouse from DB, bool is warehouse found
func GetWarehousesFromDB (uuids []string, tdb rdb.Database) ([]Warehouse, bool, error) {
	// Get warehouse json
	someJson, getError := tdb.MGetJsonData(".", uuids)
	if getError != nil {
		if fmt.Sprint(getError) != "redis: nil" {
			// warehouse not found
			return []Warehouse{}, false, nil
		}
		// error
		return []Warehouse{}, false, getError
	}
	// Got successfully, unmarshal
	someData := make([]Warehouse, len(someJson))
	for i, tempjson := range someJson {
		data := Warehouse{}
		unmarshalErr := json.Unmarshal(tempjson, &data)
		if unmarshalErr != nil {
			log.Error.Fatalf("Could not unmarshal warehouse json from DB: %v", unmarshalErr)
			return []Warehouse{}, false, unmarshalErr
		}
		someData[i] = data
	}
	
	return someData, true, nil
}

// Get warehousedata at path from DB, bool is warehouse found
func GetWarehouseDataAtPathFromDB (uuid string, path string, tdb rdb.Database) (interface{}, bool, error) {
	// Get warehouse json
	someJson, getError := tdb.GetJsonData(uuid, path)
	if getError != nil {
		if fmt.Sprint(getError) != "redis: nil" {
			// warehouse not found
			return nil, false, nil
		}
		// error
		return nil, false, getError
	}
	// Got successfully, unmarshal
	var someData interface{}
	unmarshalErr := json.Unmarshal(someJson, &someData)
	if unmarshalErr != nil {
		log.Error.Fatalf("Could not unmarshal warehouse json from DB: %v", unmarshalErr)
		return nil, false, unmarshalErr
	}
	return someData, true, nil
}

// Attempt to save warehouse, returns error or nil if successful
func SaveWarehouseToDB(tdb rdb.Database, warehouseData *Warehouse) error {
	log.Debug.Printf("Saving warehouse %s to DB", warehouseData.UUID)
	err := tdb.SetJsonData(warehouseData.UUID, ".", warehouseData)
	// creationSuccess := rdb.CreateWarehouse(tdb, warehousename, uuid, 0)
	return err
}

// Attempt to save warehouse data at path, returns error or nil if successful
func SaveWarehouseDataAtPathToDB(tdb rdb.Database, uuid string, path string, newValue interface{}) error {
	log.Debug.Printf("Saving warehouse data at path %s to DB for uuid %s", path, uuid)
	err := tdb.SetJsonData(uuid, path, newValue)
	return err
}