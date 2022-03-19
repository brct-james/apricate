// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

import (
	"apricate/log"
	"apricate/rdb"
	"encoding/json"
	"fmt"
)

// Defines a warehouse
type Warehouse struct {
	UUID string `json:"uuid" binding:"required"`
	LocationSymbol string `json:"location_symbol" binding:"required"`
	Tools map[ToolTypes]uint8 `json:"tools" binding:"required"`
	Produce map[string]uint64 `json:"produce" binding:"required"`
	Seeds map[string]uint64 `json:"seeds" binding:"required"`
	Goods map[string]uint64 `json:"goods" binding:"required"`
}

func NewEmptyWarehouse(username string, locationSymbol string) *Warehouse {
	return NewWarehouse(username, locationSymbol, make(map[ToolTypes]uint8), make(map[string]uint64), make(map[string]uint64), make(map[string]uint64))
}

func NewWarehouse(username string, locationSymbol string, starting_tools map[ToolTypes]uint8, starting_produce map[string]uint64, starting_seeds map[string]uint64, starting_goods map[string]uint64) *Warehouse {
	return &Warehouse{
		UUID: username + "|Warehouse-" + locationSymbol,
		LocationSymbol: locationSymbol,
		Tools: starting_tools,
		Produce: starting_produce,
		Seeds: starting_seeds,
		Goods: starting_goods,
	}
}

func (w *Warehouse) AddTools(name ToolTypes, quantity uint8) *Warehouse {
	w.Tools[name] += quantity
	return w
}

func (w *Warehouse) RemoveTools(name ToolTypes, quantity uint8) *Warehouse {
	w.Tools[name] -= quantity
	log.Important.Printf("%v", w.Tools[name])
	if w.Tools[name] <= 0 {
		delete(w.Tools, name)
	}
	return w
}

func (w *Warehouse) AddProduce(name string, quantity uint64) *Warehouse {
	w.Produce[name] += quantity
	return w
}

func (w *Warehouse) RemoveProduce(name string, quantity uint64) *Warehouse {
	w.Produce[name] -= quantity
	log.Important.Printf("%v", w.Produce[name])
	if w.Produce[name] <= 0 {
		delete(w.Produce, name)
	}
	return w
}

func (w *Warehouse) AddGoods(name string, quantity uint64) *Warehouse {
	w.Goods[name] += quantity
	return w
}

func (w *Warehouse) RemoveGoods(name string, quantity uint64) *Warehouse {
	w.Goods[name] -= quantity
	log.Important.Printf("%v", w.Goods[name])
	if w.Goods[name] <= 0 {
		delete(w.Goods, name)
	}
	return w
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