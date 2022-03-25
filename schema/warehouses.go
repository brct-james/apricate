// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

import (
	"apricate/log"
	"apricate/rdb"
	"encoding/json"
	"fmt"
	"strings"
)

// Defines a warehouse
type Warehouse struct {
	UUID string `json:"uuid" binding:"required"`
	LocationSymbol string `json:"location_symbol" binding:"required"`
	Tools map[string]uint64 `json:"tools" binding:"required"`
	Produce map[string]Produce `json:"produce" binding:"required"`
	Seeds map[string]uint64 `json:"seeds" binding:"required"`
	Goods map[string]uint64 `json:"goods" binding:"required"`
}

func NewEmptyWarehouse(username string, locationSymbol string) *Warehouse {
	return NewWarehouse(username, locationSymbol, make(map[string]uint64), make(map[string]Produce), make(map[string]uint64), make(map[string]uint64))
}

func NewWarehouse(username string, locationSymbol string, starting_tools map[string]uint64, starting_produce map[string]Produce, starting_seeds map[string]uint64, starting_goods map[string]uint64) *Warehouse {
	return &Warehouse{
		UUID: username + "|Warehouse-" + locationSymbol,
		LocationSymbol: locationSymbol,
		Tools: starting_tools,
		Produce: starting_produce,
		Seeds: starting_seeds,
		Goods: starting_goods,
	}
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

func (w *Warehouse) GetProduce(name string, size Size) *Produce {
	produceName := name + "|" + size.String()
	if entry, ok := w.Produce[produceName]; ok {
		return &entry
	}
	return nil
}

func (w *Warehouse) GetSimpleProduceDict() map[string]uint64 {
	res := make(map[string]uint64, len(w.Produce))
	for key, entry := range w.Produce {
		res[key] = entry.Quantity
	}
	return res
}

func (w *Warehouse) SetSimpleProduceDict(produce map[string]uint64) {
	res := make(map[string]Produce, len(produce))
	for key, quantity := range produce {
		res[key] = *NewProduce(strings.Split(key, "|")[0], SizeToID[strings.Split(key, "|")[1]], quantity)
	}
	w.Produce = res
}

func (w *Warehouse) AddProduce(name string, size Size, quantity uint64) {
	produceName := name + "|" + size.String()
	if entry, ok := w.Produce[produceName]; ok {
		entry.Quantity += quantity
		w.Produce[produceName] = entry
	} else {
		w.Produce[produceName] = Produce{
			Good: Good{
				Name:name,
				Quantity:quantity,
			},
			Size: size,
		}
	}
}

func (w *Warehouse) RemoveProduce(name string, size Size, quantity uint64) {
	produceName := name + "|" + size.String()
	if entry, ok := w.Produce[produceName]; ok {
		entry.Quantity -= quantity
		if(entry.Quantity <= 0) {
			delete (w.Produce, produceName)
		} else {
			w.Produce[produceName] = entry
		}
	} else {
		log.Error.Printf("Cannot add produce, !ok for name: %s", name)
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