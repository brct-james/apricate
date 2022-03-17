// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

import (
	"apricate/rdb"
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

func NewPlot(username string, countOfPlots uint16, locationSymbol string, capacity Size) *Plot {
	return &Plot{
		UUID: username + "|Farm-" + locationSymbol + "|Plot-" + fmt.Sprintf("%d", countOfPlots),
		LocationSymbol: locationSymbol,
		PlotSize: capacity,
		PlantedPlant: nil,
		Quantity: 0,
	}
}

func NewPlots(pdb rdb.Database, username string, countOfPlots uint16, locationSymbol string, capacities []Size) map[string]Plot {
	res := make(map[string]Plot, len(capacities))
	for i, size := range capacities {
		plot := NewPlot(username, countOfPlots + uint16(i), locationSymbol, size)
		res[plot.UUID] = *plot
	}
	return res
}

// Defines a plot plant request body
type PlotPlantBody struct {
	SeedName string `json:"seed_name" binding:"required"`
	SeedQuantity uint64`json:"seed_quantity" binding:"required"`
	SeedSize Size `json:"size" binding:"required"`
}

// Defines a plot action request body
type PlotInteractBody struct {
	Action GrowthAction `json:"action" binding:"required"`
	Consumables Good `json:"consumables,omitempty"`
}

// Defines a plot action response body
type PlotActionResponse struct {
	Plot *Plot `json:"plot" binding:"required"`
	NextStage *GrowthStage `json:"next_stage" binding:"required"`
}

// // Handles planting a plant in the plot
// func (p *Plot) Plant(w http.ResponseWriter, dbs map[string]rdb.Database, plantDictionary map[string]PlantDefinition, body PlotPlantBody) Plot {
	// if p.PlantedPlant != nil {
	// 	// fail case, already planted, clear first
	// 	log.Error.Println("already planted, clear first")
	// 	return Plot{}
	// }
	// if consumables.Quantity > uint64(p.PlotSize) {
	// 	// fail case, not large enough for specified quantity
	// 	log.Error.Println("Plot not large enough for specified quantity")
	// 	return Plot{}
	// }
	// plantDefinition, ok := plantDictionary[strings.Split(consumables.Name, " Seeds")[0]]
	// if !ok {
	// 	// fail case, consumables.Name not in plantDictionary
	// 	log.Error.Printf("consumables.Name %s not in plantDictionary", consumables.Name)
	// 	return Plot{}
	// }
	// plantingGrowthStage := plantDefinition.GrowthStages[0]
	// seedQuantityNeeded := consumables.Quantity * plantingGrowthStage.ConsumableOptions[0].Quantity
	// seedQuantityOwned := farmWarehouse.Goods[plantDefinition.SeedName]
	// if seedQuantityOwned < uint64(seedQuantityNeeded) {
	// 	// fail case, not enough seeds in local warehouse
	// 	log.Error.Println("Not enough of specified seeds in local warehouse")
	// 	return Plot{}
	// }
	// // Action should never be tool-less (never /wait or /clear)
	// toolTypeNeeded := plantingGrowthStage.Action.ToolType()
	// if _, ok := farmTools[toolTypeNeeded]; !ok {
	// 	// fail case, don't have necessary tool
	// 	log.Error.Println("Dont have necessary tool")
	// 	return Plot{}
	// }

	// // Plant Plant
	// p.PlantedPlant = NewPlant(plantDefinition.Name, size)
	// p.PlantedPlant.CurrentStage ++
	// p.Quantity = uint16(consumables.Quantity)
	// //Save to DB
	// SavePlotToDB(fdb, p)

	// // Send response
	// nextStage, nextStageErr := plantDefinition.GetScaledGrowthStage(int(p.PlantedPlant.CurrentStage), uint64(p.Quantity), p.PlantedPlant.Size)
	// if nextStageErr != nil {
	// 	log.Error.Printf("Error in Interact, could not get scaled growth stage.")
	// 	responses.SendRes(w, responses.JSON_Marshal_Error, nil, "")
	// 	return *p
	// }
	// res := PlotActionResponse{Plot: p, NextStage: nextStage}
	// getPlotActionResponseJsonString, getPlotActionResponseJsonStringErr := responses.JSON(res)
	// if getPlotActionResponseJsonStringErr != nil {
	// 	log.Error.Printf("Error in Interact, could not format plot action response as JSON. res: %v, error: %v", res, getPlotActionResponseJsonStringErr)
	// 	responses.SendRes(w, responses.JSON_Marshal_Error, nil, getPlotActionResponseJsonStringErr.Error())
	// 	return *p
	// }
	// log.Debug.Printf("Sending response for Interact:\n%v", getPlotActionResponseJsonString)
	// responses.SendRes(w, responses.Generic_Success, res, "")
	// return *p
// }