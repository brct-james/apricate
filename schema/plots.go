// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

import (
	"apricate/rdb"
	"apricate/responses"
	"fmt"
)

// Defines a plot
type Plot struct {
	UUID string `json:"uuid" binding:"required"`
	LocationSymbol string `json:"location_symbol" binding:"required"`
	PlotSize Size `json:"size" binding:"required"`
	Quantity uint16 `json:"plant_quantity" binding:"required"`
	PlantedPlant *Plant `json:"plant" binding:"required"`
}

func NewPlot(username string, countOfPlots uint16, locationSymbol string, capacity Size) *Plot {
	return &Plot{
		UUID: username + "|Farm-" + locationSymbol + "|Plot-" + fmt.Sprintf("%d", countOfPlots),
		LocationSymbol: locationSymbol,
		PlotSize: capacity,
		Quantity: 0,
		PlantedPlant: nil,
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
	SeedName string `json:"name" binding:"required"`
	SeedQuantity uint16`json:"quantity" binding:"required"`
	SeedSize Size `json:"size" binding:"required"`
}

// Defines a plot action request body
type PlotInteractBody struct {
	Action GrowthAction `json:"action" binding:"required"`
	Consumables Good `json:"consumables,omitempty"`
}

// Defines a plot plant response body
type PlotPlantResponse struct {
	Warehouse *Warehouse `json:"warehouse" binding:"required"`
	Plot *Plot `json:"plot" binding:"required"`
	NextStage *GrowthStage `json:"next_stage" binding:"required"`
}

// Defines a plot action response body
type PlotActionResponse struct {
	Plot *Plot `json:"plot" binding:"required"`
	NextStage *GrowthStage `json:"next_stage" binding:"required"`
}

func (p *Plot) IsPlantable(ppb PlotPlantBody) responses.ResponseCode {
	if p.PlantedPlant != nil {
		return responses.Plot_Already_Planted
	}
	if ppb.SeedQuantity > uint16(p.PlotSize/ppb.SeedSize) {
		return responses.Plot_Too_Small
	}
	return responses.Generic_Success
}