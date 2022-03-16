// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

import (
	"apricate/log"
	"apricate/rdb"
	"encoding/json"
	"fmt"
)

// Defines a plot
type Plot struct {
	UUID string `json:"uuid" binding:"required"`
	LocationSymbol string `json:"location_symbol" binding:"required"`
	PlotSize Size `json:"size" binding:"required"`
	PlantedPlant *Plant `json:"plant" binding:"required"`
	Quantity uint16 `json:"plant_quantity" binding:"required"`
}

func NewPlot(username string, countOfPlots int16, locationSymbol string, capacity Size) *Plot {
	return &Plot{
		UUID: username + "|" + locationSymbol + "|" + fmt.Sprintf("%d", countOfPlots),
		LocationSymbol: locationSymbol,
		PlotSize: capacity,
		PlantedPlant: nil,
		Quantity: 0,
	}
}

func NewPlots(pdb rdb.Database, username string, countOfPlots int16, locationSymbol string, capacities []Size) []string {
	res := make([]string, len(capacities))
	for i, size := range capacities {
		plot := NewPlot(username, countOfPlots + int16(i), locationSymbol, size)
		SavePlotToDB(pdb, plot)
		res[i] = plot.UUID
	}
	return res
}

// Defines a plot action request body
type PlotActionBody struct {
	Action GrowthAction `json:"action" binding:"required"`
	Consumables Good `json:"consumables,omitempty"`
}

// Handles planting a plant in the plot
func (p *Plot) Plant(pdb rdb.Database, plantDictionary map[string]PlantDefinition, farmWarehouse Warehouse, farmTools map[ToolTypes]uint8, newPlant PlantType, quantity uint16) (Plot, bool) {
	if p.PlantedPlant != nil {
		// fail case, already planted, clear first
		return Plot{}, false
	}
	if quantity > uint16(p.PlotSize) {
		// fail case, not large enough for specified quantity
		log.Error.Println("Plot not large enough for specified quantity")
		return Plot{}, false
	}
	plantDefinition := plantDictionary[newPlant.String()]
	plantingGrowthStage := plantDefinition.GrowthStages[0]
	seedQuantityNeeded := quantity * uint16(plantingGrowthStage.ConsumableOptions[0].Quantity)
	seedQuantityOwned := farmWarehouse.Goods[plantDefinition.SeedName].Quantity
	if seedQuantityOwned < uint64(seedQuantityNeeded) {
		// fail case, not enough seeds in local warehouse
		log.Error.Println("Not enough of specified seeds in local warehouse")
		return Plot{}, false
	}
	// Action should never be tool-less (never /wait or /clear)
	toolTypeNeeded := plantingGrowthStage.Action.ToolType()
	if _, ok := farmTools[toolTypeNeeded]; !ok {
		// fail case, don't have necessary tool
		log.Error.Println("Dont have necessary tool")
		return Plot{}, false
	}

	// Plant Plant
	p.PlantedPlant = NewPlant(newPlant)
	p.Quantity = quantity
	//Save to DB
	SavePlotToDB(pdb, p)
	return *p, true
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