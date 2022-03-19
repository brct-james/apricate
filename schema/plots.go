// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

import (
	"apricate/log"
	"apricate/rdb"
	"apricate/responses"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
)

// Defines a plot
type Plot struct {
	UUID string `json:"uuid" binding:"required"`
	LocationSymbol string `json:"location_symbol" binding:"required"`
	PlotSize Size `json:"size" binding:"required"`
	Quantity uint16 `json:"plant_quantity" binding:"required"`
	PlantedPlant *Plant `json:"plant" binding:"required"`
	GrowthCompleteTimestamp int64 `json:"growth_complete_timestamp" binding:"required"`
}

func NewPlot(username string, countOfPlots uint16, locationSymbol string, capacity Size) *Plot {
	return &Plot{
		UUID: username + "|Farm-" + locationSymbol + "|Plot-" + fmt.Sprintf("%d", countOfPlots),
		LocationSymbol: locationSymbol,
		PlotSize: capacity,
		Quantity: 0,
		PlantedPlant: nil,
		GrowthCompleteTimestamp: 0,
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

// Defines HarvestProduce data
type HarvestProduce struct {
	Produce []Produce `json:"produce" binding:"required"`
	Seeds map[string]uint64 `json:"seeds" binding:"required"`
	Goods map[string]uint64 `json:"goods" binding:"required"`
}

// Defines a plot plant request body
type PlotPlantBody struct {
	SeedName string `json:"name" binding:"required"`
	SeedQuantity uint16`json:"quantity" binding:"required"`
	SeedSize Size `json:"size" binding:"required"`
}

// Defines a plot action request body
type PlotInteractBody struct {
	Action string `json:"action" binding:"required"`
	Consumable string `json:"consumable,omitempty"`
}

// Defines a plot plant response body
type PlotPlantResponse struct {
	Warehouse *Warehouse `json:"warehouse" binding:"required"`
	Plot *Plot `json:"plot" binding:"required"`
	NextStage *GrowthStage `json:"next_stage" binding:"required"`
}

// Defines a plot action response body
type PlotActionResponse struct {
	Warehouse *Warehouse `json:"warehouse,omitempty"`
	Plot *Plot `json:"plot" binding:"required"`
	NextStage *GrowthStage `json:"next_stage" binding:"required"`
}

func (p *Plot) IsPlantable(ppb PlotPlantBody) responses.ResponseCode {
	if p.PlantedPlant != nil {
		return responses.Plot_Already_Planted
	}
	if ppb.SeedSize <= 0 {
		return responses.Bad_Request
	}
	if ppb.SeedQuantity > uint16(p.PlotSize/ppb.SeedSize) {
		return responses.Plot_Too_Small
	}
	return responses.Generic_Success
}

// returns ResponseCode, AddedYield, ConsumableQuantityUsed, GrowthHarvest, Cooldown/GrowthTime
func (p *Plot) IsInteractable(pib PlotInteractBody, plantDef PlantDefinition, consumableQuantityAvailable uint64, tools map[ToolTypes]uint8) (responses.ResponseCode, float64, uint64, *GrowthHarvest, int64) {
	consumableName := strings.Title(strings.ToLower(pib.Consumable))
	pib.Action = strings.Title(strings.ToLower(pib.Action))
	log.Debug.Println(pib.Action)
	growthStage := plantDef.GrowthStages[p.PlantedPlant.CurrentStage]
	// if blank action
	if pib.Action == string("") {
		// action sent is missing or invalid
		return responses.Invalid_Plot_Action, 0, 0, nil, 0
	}
	// if has and is skip action
	if growthStage.ActionToSkip != nil && pib.Action == (*growthStage.ActionToSkip).String() {
		if _, countOk := tools[growthStage.ActionToSkip.ToolType()]; countOk {
			// Success, have correct tool
			return responses.Generic_Success, 0, 0, nil, 0
		}
		// Failure, don't have correct tool
		return responses.Tool_Not_Found, 0, 0, nil, 0
	}
	// if action action
	if pib.Action == (*growthStage.Action).String() {
		if _, countOk := tools[growthStage.Action.ToolType()]; !countOk {
			// Failure, don't have correct tool
			return responses.Tool_Not_Found, 0, 0, nil, 0
		}
		// Growth stage contains no consumables, return success
		if len(growthStage.ConsumableOptions) == 0 {
			// if harvest step, return harvest data, else just return added yield
			if growthStage.Harvestable != nil {
				return responses.Generic_Success, growthStage.AddedYield, 0, growthStage.Harvestable, 0
			}
			return responses.Generic_Success, growthStage.AddedYield, 0, nil, *growthStage.GrowthTime
		}
		// Check consumables
		if consumableName == string("") {
			// No consumables included in request body, fail
			return responses.Missing_Consumable_Selection, 0, 0, nil, 0
		}
		// Check all consumables for option matching request, return in loop if passes, else fail after
		scaledGrowthStage, sGSErr := plantDef.GetScaledGrowthStage(p.PlantedPlant.CurrentStage, uint64(p.Quantity), p.PlantedPlant.Size)
		if sGSErr != nil {
			// internal server error, could not get scaled growth stage
			return responses.Internal_Server_Error, 0, 0, nil, 0
		}
		for _, consumableOption := range scaledGrowthStage.ConsumableOptions {
			if consumableOption.Name == consumableName {
				// found matching consumable option
				if consumableOption.Quantity <= consumableQuantityAvailable {
					// have enough, return success
					if growthStage.Harvestable != nil {
						// if harvest step, return harvest data, else just return added yield
						return responses.Generic_Success, growthStage.AddedYield + consumableOption.AddedYield, consumableOption.Quantity, growthStage.Harvestable, *growthStage.GrowthTime
					}
					return responses.Generic_Success, growthStage.AddedYield + consumableOption.AddedYield, consumableOption.Quantity, nil, *growthStage.GrowthTime
				}
				// insufficient quantity in local warehouse
				return responses.Not_Enough_Items_In_Warehouse, 0, 0, nil, 0
			}
		}
		return responses.Consumable_Not_In_Options, 0, 0, nil, 0
	}
	// else, invalid action specified
	return responses.Invalid_Plot_Action, 0, 0, nil, 0
}

func (p *Plot) CalculateProduce(growthHarvest *GrowthHarvest) HarvestProduce {
	harvest := HarvestProduce{
		Produce: make([]Produce, 0),
		Seeds: make(map[string]uint64),
		Goods: make(map[string]uint64),
	}
	totalYield := p.PlantedPlant.Yield
	size := p.PlantedPlant.Size
	rand.Seed(time.Now().UnixNano())
	randMin := 0.8
	randMax := 1.2
	yieldRNG := randMin + rand.Float64() * (randMax - randMin)
	log.Debug.Println(growthHarvest)
	// Calculate Produce - Quantity Affected By AddedYield NOT Size
	for produceName, yieldModifier := range growthHarvest.Produce {
		harvest.Produce = append(harvest.Produce, *NewProduce(produceName, size, uint64(math.Round(float64(p.Quantity) * yieldRNG * totalYield * yieldModifier))))
	}
	// Calculate Seeds - NOT Affected By AddedYield OR Size (Affected by Seed Yield Modifier and Yield RNG, however)
	for seedName, yieldModifier := range growthHarvest.Seeds {
		harvest.Seeds[seedName] = uint64(math.Round(float64(p.Quantity) * yieldRNG * yieldModifier))
	}
	// Calculate Goods - Quantity Affected By AddedYield AND Size
	for goodName, yieldModifier := range growthHarvest.Goods {
		harvest.Goods[goodName] = uint64(math.Round(float64(size) * float64(p.Quantity) * yieldRNG * totalYield * yieldModifier))
	}
	return harvest
}