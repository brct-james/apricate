// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

import (
	"apricate/log"
	"apricate/rdb"
	"encoding/json"
	"fmt"
	"strings"
)

// Define Wareset
type Wareset struct {
	Tools map[string]uint64 `yaml:"Tools" json:"tools,omitempty"`
	Produce map[string]uint64 `yaml:"Produce" json:"produce,omitempty"`
	Seeds map[string]uint64 `yaml:"Seeds" json:"seeds,omitempty"`
	Goods map[string]uint64 `yaml:"Goods" json:"goods,omitempty"`
}

// Defines a warehouse
type Warehouse struct {
	UUID string `json:"uuid" binding:"required"`
	LocationSymbol string `json:"location_symbol" binding:"required"`
	Wareset
}

func NewEmptyWarehouse(username string, locationSymbol string) *Warehouse {
	return NewWarehouse(username, locationSymbol, make(map[string]uint64), make(map[string]uint64), make(map[string]uint64), make(map[string]uint64))
}

func NewWarehouse(username string, locationSymbol string, starting_tools map[string]uint64, starting_produce map[string]uint64, starting_seeds map[string]uint64, starting_goods map[string]uint64) *Warehouse {
	return &Warehouse{
		UUID: username + "|Warehouse-" + locationSymbol,
		LocationSymbol: locationSymbol,
		Wareset: Wareset{
			Tools: starting_tools,
			Produce: starting_produce,
			Seeds: starting_seeds,
			Goods: starting_goods,
		},
	}
}

func (w *Warehouse) TotalSize() uint64 {
	size := uint64(0)
	for _, q := range w.Goods {
		size += q
	}
	for _, q := range w.Tools {
		size += q
	}
	for _, q := range w.Produce {
		size += q
	}
	for _, q := range w.Seeds {
		size += q
	}
	return size
}

func (w *Warehouse) AddTools(name string, quantity uint64) {
	w.Tools[name] += quantity
}

func (w *Warehouse) RemoveTools(name string, quantity uint64) {
	w.Tools[name] -= quantity
	if w.Tools[name] <= 0 {
		delete(w.Tools, name)
	}
}

func (w *Warehouse) GetProduceNameSizeSlice(name string) (string, string, bool) {
	slice := strings.Split(name, "|")
	if len(slice) < 2 {
		return "", "", false
	}
	return slice[0], slice[1], true
}

func (w *Warehouse) AddProduce(name string, quantity uint64) {
	w.Produce[name] += quantity
}

func (w *Warehouse) RemoveProduce(name string, quantity uint64) {
	w.Produce[name] -= quantity
	if w.Produce[name] <= 0 {
		delete(w.Produce, name)
	}
}

func (w *Warehouse) AddSeeds(name string, quantity uint64) {
	w.Seeds[name] += quantity
}

func (w *Warehouse) RemoveSeeds(name string, quantity uint64) {
	w.Seeds[name] -= quantity
	if w.Seeds[name] <= 0 {
		delete(w.Seeds, name)
	}
}

func (w *Warehouse) AddGoods(name string, quantity uint64) {
	w.Goods[name] += quantity
}

func (w *Warehouse) RemoveGoods(name string, quantity uint64) {
	w.Goods[name] -= quantity
	if w.Goods[name] <= 0 {
		delete(w.Goods, name)
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
		if fmt.Sprint(getError) == "redis: nil" {
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
		if fmt.Sprint(getError) == "redis: nil" {
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
		if fmt.Sprint(getError) == "redis: nil" {
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

// Attempt to delete warehouse, returns error or nil if successful
func DeleteWarehouseFromDB(tdb rdb.Database, uuid string) error {
	log.Debug.Printf("Saving warehouse %s to DB", uuid)
	_, err := tdb.DelJsonData(uuid, ".")
	// creationSuccess := rdb.CreateWarehouse(tdb, warehousename, uuid, 0)
	return err
}