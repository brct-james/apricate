// Package handlers provides functions for handling web routes
package handlers

import (
	"apricate/auth"
	"apricate/log"
	"apricate/metrics"
	"apricate/rdb"
	"apricate/responses"
	"apricate/schema"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// HANDLER FUNCTIONS

// Handler function for the secure route: /api/my/account
type AccountInfo struct {
	Dbs *map[string]rdb.Database
}
func (h *AccountInfo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- accountInfo --"))
	udb := (*h.Dbs)["users"]
	OK, userData, _ := secureGetUser(w, r, udb)
	if !OK {
		return // Failure states handled by secureGetUser, simply return
	}
	getUserJsonString, getUserJsonStringErr := responses.JSON(userData)
	if getUserJsonStringErr != nil {
		log.Important.Printf("in AccountInfo, could not format thisUser as JSON. userData: %v, error: %v", userData, getUserJsonStringErr)
		responses.SendRes(w, responses.JSON_Marshal_Error, userData, getUserJsonStringErr.Error())
		return
	}
	log.Debug.Printf("Sending response for AccountInfo:\n%v", getUserJsonString)
	responses.SendRes(w, responses.Generic_Success, userData, "")
	log.Debug.Println(log.Cyan("-- End accountInfo --"))
}

// Handler function for the secure route: /api/my/assistants
type AssistantsInfo struct {
	Dbs *map[string]rdb.Database
}
func (h *AssistantsInfo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- AssistantsInfo --"))
	udb := (*h.Dbs)["users"]
	OK, userData, _ := secureGetUser(w, r, udb)
	if !OK {
		return // Failure states handled by secureGetUser, simply return
	}
	adb := (*h.Dbs)["assistants"]
	assistants, foundAssistants, assistantsErr := schema.GetAssistantsFromDB(userData.Assistants, adb)
	if assistantsErr != nil || !foundAssistants {
		log.Error.Printf("Error in AssistantsInfo, could not get assistants from DB. foundAssistants: %v, error: %v", foundAssistants, assistantsErr)
		responses.SendRes(w, responses.DB_Get_Failure, assistants, assistantsErr.Error())
		return
	}
	getAssistantJsonString, getAssistantJsonStringErr := responses.JSON(assistants)
	if getAssistantJsonStringErr != nil {
		log.Error.Printf("Error in AssistantsInfo, could not format assistants as JSON. assistants: %v, error: %v", assistants, getAssistantJsonStringErr)
		responses.SendRes(w, responses.JSON_Marshal_Error, assistants, getAssistantJsonStringErr.Error())
		return
	}
	log.Debug.Printf("Sending response for AssistantsInfo:\n%v", getAssistantJsonString)
	responses.SendRes(w, responses.Generic_Success, assistants, "")
	log.Debug.Println(log.Cyan("-- End AssistantsInfo --"))
}

// Handler function for the secure route: /api/my/assistants/{uuid}
type AssistantInfo struct {
	Dbs *map[string]rdb.Database
}
func (h *AssistantInfo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- AssistantInfo --"))
	// Get symbol from route
	id := GetVarEntries(r, "assistant-id", AllCaps)
	// Get userinfoContext from validation middleware
	userInfo, userInfoErr := GetValidationFromCtx(r)
	if userInfoErr != nil {
		// Fail state getting context
		log.Error.Printf("Could not get validationpair in AssistantInfo")
		userInfoErrMsg := fmt.Sprintf("userInfo is nil, check auth validation context %v:\n%v", auth.ValidationContext, r.Context().Value(auth.ValidationContext))
		responses.SendRes(w, responses.No_AuthPair_Context, nil, userInfoErrMsg)
		return
	}
	uuid := userInfo.Username + "|Assistant-" + id
	log.Debug.Printf("AssistantInfo Requested for: %s", uuid)
	adb := (*h.Dbs)["assistants"]
	assistant, foundAssistant, assistantsErr := schema.GetAssistantFromDB(uuid, adb)
	if assistantsErr != nil || !foundAssistant {
		log.Debug.Printf("in AssistantInfo, could not get assistant from DB. foundAssistant: %v, error: %v", foundAssistant, assistantsErr)
		responses.SendRes(w, responses.DB_Get_Failure, assistant, assistantsErr.Error())
		return
	}
	getAssistantJsonString, getAssistantJsonStringErr := responses.JSON(assistant)
	if getAssistantJsonStringErr != nil {
		log.Error.Printf("Error in AssistantInfo, could not format assistants as JSON. assistants: %v, error: %v", assistant, getAssistantJsonStringErr)
		responses.SendRes(w, responses.JSON_Marshal_Error, assistant, getAssistantJsonStringErr.Error())
		return
	}
	log.Debug.Printf("Sending response for AssistantInfo:\n%v", getAssistantJsonString)
	responses.SendRes(w, responses.Generic_Success, assistant, "")
	log.Debug.Println(log.Cyan("-- End AssistantInfo --"))
}

// Handler function for the secure route: /api/my/caravans
type CaravansInfo struct {
	Dbs *map[string]rdb.Database
}
func (h *CaravansInfo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- CaravansInfo --"))
	udb := (*h.Dbs)["users"]
	OK, userData, _ := secureGetUser(w, r, udb)
	if !OK {
		return // Failure states handled by secureGetUser, simply return
	}
	adb := (*h.Dbs)["caravans"]
	caravans, foundCaravans, caravansErr := schema.GetCaravansFromDB(userData.Caravans, adb)
	if caravansErr != nil {
		log.Error.Printf("Error in CaravansInfo, could not get caravans from DB. foundCaravans: %v, error: %v", foundCaravans, caravansErr)
		responses.SendRes(w, responses.DB_Get_Failure, caravans, "Could not get caravans, no err?")
		return
	}
	if !foundCaravans {
		log.Debug.Printf("in CaravansInfo, could not get caravans from DB. foundCaravans: %v, error: %v, probably just none exist", foundCaravans, caravansErr)
		responses.SendRes(w, responses.Generic_Success, caravans, " None Found")
		return
	}
	getCaravanJsonString, getCaravanJsonStringErr := responses.JSON(caravans)
	if getCaravanJsonStringErr != nil {
		log.Error.Printf("Error in CaravansInfo, could not format caravans as JSON. caravans: %v, error: %v", caravans, getCaravanJsonStringErr)
		responses.SendRes(w, responses.JSON_Marshal_Error, caravans, getCaravanJsonStringErr.Error())
		return
	}
	log.Debug.Printf("Sending response for CaravansInfo:\n%v", getCaravanJsonString)
	responses.SendRes(w, responses.Generic_Success, caravans, "")
	log.Debug.Println(log.Cyan("-- End CaravansInfo --"))
}

// Handler function for the secure route: /api/my/caravans/{uuid}
type CaravanInfo struct {
	Dbs *map[string]rdb.Database
}
func (h *CaravanInfo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- CaravanInfo --"))
	// Get symbol from route
	id := GetVarEntries(r, "caravan-id", AllCaps)
	// Get userinfoContext from validation middleware
	userInfo, userInfoErr := GetValidationFromCtx(r)
	if userInfoErr != nil {
		// Fail state getting context
		log.Error.Printf("Could not get validationpair in CaravanInfo")
		userInfoErrMsg := fmt.Sprintf("userInfo is nil, check auth validation context %v:\n%v", auth.ValidationContext, r.Context().Value(auth.ValidationContext))
		responses.SendRes(w, responses.No_AuthPair_Context, nil, userInfoErrMsg)
		return
	}
	uuid := userInfo.Username + "|Caravan-" + id
	log.Debug.Printf("CaravanInfo Requested for: %s", uuid)
	adb := (*h.Dbs)["caravans"]
	caravan, foundCaravan, caravansErr := schema.GetCaravanFromDB(uuid, adb)
	if caravansErr != nil || !foundCaravan {
		log.Debug.Printf("in CaravanInfo, could not get caravan from DB. foundCaravan: %v, error: %v", foundCaravan, caravansErr)
		responses.SendRes(w, responses.DB_Get_Failure, caravan, caravansErr.Error())
		return
	}
	getCaravanJsonString, getCaravanJsonStringErr := responses.JSON(caravan)
	if getCaravanJsonStringErr != nil {
		log.Error.Printf("Error in CaravanInfo, could not format caravans as JSON. caravans: %v, error: %v", caravan, getCaravanJsonStringErr)
		responses.SendRes(w, responses.JSON_Marshal_Error, caravan, getCaravanJsonStringErr.Error())
		return
	}
	log.Debug.Printf("Sending response for CaravanInfo:\n%v", getCaravanJsonString)
	responses.SendRes(w, responses.Generic_Success, caravan, "")
	log.Debug.Println(log.Cyan("-- End CaravanInfo --"))
}

// Handler function for the secure route: /api/my/plots/{uuid}/interact
type CharterCaravan struct {
	Dbs *map[string]rdb.Database
	MainDictionary *schema.MainDictionary
}
func (h *CharterCaravan) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- CharterCaravan --"))
	// Get userinfoContext from validation middleware
	userInfo, userInfoErr := GetValidationFromCtx(r)
	if userInfoErr != nil {
		// Fail state getting context
		log.Error.Printf("Could not get validationpair in AssistantInfo")
		userInfoErrMsg := fmt.Sprintf("userInfo is nil, check auth validation context %v:\n%v", auth.ValidationContext, r.Context().Value(auth.ValidationContext))
		responses.SendRes(w, responses.No_AuthPair_Context, nil, userInfoErrMsg)
		return
	}

	// unmarshall request body to get charter including wares if applicable
	var body schema.CaravanCharter
	decoder := json.NewDecoder(r.Body)
	if decodeErr := decoder.Decode(&body); decodeErr != nil {
		// Fail case, could not decode
		errmsg := fmt.Sprintf("Decode Error in CharterCaravan: %v", decodeErr)
		log.Debug.Printf(errmsg)
		responses.SendRes(w, responses.Bad_Request, nil, "Could not decode request body, ensure it conforms to expected format.")
		return
	}

	// Get Assistants


	// Get warehouse
	wdb := (*h.Dbs)["warehouses"]
	locationSymbol := strings.Split(symbolSlice[1], "-")
	warehouseLocationSymbol := symbolSlice[0] + "|Warehouse-" + strings.Join(locationSymbol[1:],"-")
	warehouse, foundWarehouse, warehousesErr := schema.GetWarehouseFromDB(warehouseLocationSymbol, wdb)
	if warehousesErr != nil || !foundWarehouse {
		errmsg := fmt.Sprintf("Error in CharterCaravan, could not get warehouse from DB. foundWarehouse: %v, error: %v", foundWarehouse, warehousesErr)
		log.Error.Printf(errmsg)
		responses.SendRes(w, responses.DB_Get_Failure, nil, errmsg)
		return
	}

	consumableQuantityAvailable := uint64(0)
	consumableName := strings.Title(strings.ToLower(body.Consumable))
	// If consumables included, validate them
	if consumableName != string("") {
		// Validate specified consumable is a good
		goodsDict := (*h.MainDictionary).Goods
		if _, ok := goodsDict[consumableName]; !ok {
			// Fail, seed is not good
			errmsg := fmt.Sprintf("in CharterCaravan, consumable item does not exist in good dictionary. received consumable name: %v", consumableName)
			log.Debug.Printf(errmsg)
			responses.SendRes(w, responses.Item_Does_Not_Exist, nil, errmsg)
			return
		}
		// Validate consumable specified in given warehouse
		temp, ownedConsumableOk := warehouse.Goods[consumableName]; 
		if !ownedConsumableOk {
			// Fail, consumable good not in warehouse
			errmsg := fmt.Sprintf("in CharterCaravan, consumable item not found in local warehouse. received good name: %v, warehouse goods: %v", consumableName, warehouse.Goods)
			log.Debug.Printf(errmsg)
			responses.SendRes(w, responses.Not_Enough_Items_In_Warehouse, nil, errmsg)
			return
		}
		consumableQuantityAvailable = temp
	}

	// Validate plot available for interaction and body meets internal plot validation
	plantDef, plantDefOk := h.MainDictionary.Plants[plot.PlantedPlant.PlantType]

	if !plantDefOk {
		plotDefErrMsg := fmt.Sprintf("Error in CharterCaravan, Plot [%s] has PlantedPlant type [%s] not found in main dictionary!", plot.UUID, plot.PlantedPlant.PlantType)
		log.Error.Println(plotDefErrMsg)
		responses.SendRes(w, responses.Internal_Server_Error, plot, plotDefErrMsg)
	}
	plotValidationResponse, addedYield, usedConsumableQuantity, growthHarvest, growthTime, errInfoMsg := plot.IsInteractable(body, plantDef, consumableQuantityAvailable, warehouse.Tools)
	switch plotValidationResponse {
	case responses.Invalid_Plot_Action:
		log.Debug.Printf("in PlotInteract, Invalid_Plot_Action")
		responses.SendRes(w, responses.Invalid_Plot_Action, plot, errInfoMsg)
		return
	case responses.Tool_Not_Found:
		log.Debug.Printf("in PlotInteract, Tool_Not_Found")
		responses.SendRes(w, responses.Tool_Not_Found, plot, errInfoMsg)
		return
	case responses.Missing_Consumable_Selection:
		log.Debug.Printf("in PlotInteract, Missing_Consumable_Selection")
		responses.SendRes(w, responses.Missing_Consumable_Selection, plot, errInfoMsg)
		return
	case responses.Internal_Server_Error:
		log.Debug.Printf("in PlotInteract, Internal_Server_Error")
		responses.SendRes(w, responses.Internal_Server_Error, plot, errInfoMsg)
		return
	case responses.Not_Enough_Items_In_Warehouse:
		log.Debug.Printf("in PlotInteract, Not_Enough_Items_In_Warehouse")
		responses.SendRes(w, responses.Not_Enough_Items_In_Warehouse, plot, errInfoMsg)
		return
	case responses.Consumable_Not_In_Options:
		log.Debug.Printf("in PlotInteract, Consumable_Not_In_Options")
		responses.SendRes(w, responses.Consumable_Not_In_Options, plot, errInfoMsg)
		return
	case responses.Generic_Success:
		log.Debug.Printf("Plot growth action validated successfully: %s, action: %s", plot.UUID, body.Action)
	default:
		log.Error.Fatalf("Received unexpected response type from plot.IsPlantable. plot: %v body: %v", plot, body)
	}

	// Update objects with results of interaction

	plot.PlantedPlant.Yield += addedYield
	log.Debug.Printf("Interact Plot, Growth Time: %d", growthTime)
	plot.GrowthCompleteTimestamp = time.Now().Unix() + growthTime
	farm.Plots[uuid] = plot

	if consumableName != string("") {
		// if consumables used
		warehouse.RemoveGoods(consumableName, usedConsumableQuantity)
	}


	// Handle updating stage and potential harvesting
	var nextStage *schema.GrowthStage

	if growthHarvest != nil {
		// if was a harvest action
		harvest := plot.CalculateProduce(growthHarvest)
		log.Debug.Println("Harvest Calculated:")
		log.Debug.Println(harvest)
		for _, produce := range harvest.Produce {
			log.Debug.Printf("Add produce quantity: %d", produce.Quantity)
			warehouse.AddProduce(produce.Name, produce.Size, produce.Quantity)
			log.Debug.Println(warehouse.Produce)
		}
		for seedname, seedquantity := range harvest.Seeds {
			log.Debug.Printf("Add seed quantity: %d", seedquantity)
			warehouse.AddSeeds(seedname, seedquantity)
		}
		for goodname, goodquantity := range harvest.Goods {
			log.Debug.Printf("Add good quantity: %d", goodquantity)
			warehouse.AddGoods(goodname, goodquantity)
		}

		metrics.TrackHarvest(plantDef.Name)
		
		// check if final harvest
		if growthHarvest.FinalHarvest {
			// is final harvest
			// clear
			plot.PlantedPlant = nil
			plot.Quantity = 0
			farm.Plots[uuid] = plot
		} else {
			// not final harvest
			// move up current stage, initialize nextStage for response
			plot.PlantedPlant.CurrentStage++
			nextStage = &h.MainDictionary.Plants[plot.PlantedPlant.PlantType].GrowthStages[plot.PlantedPlant.CurrentStage]
		}
	} else {
		// if not harvest
		// move up current stage, initialize nextStage for response
		plot.PlantedPlant.CurrentStage++
		nextStage = &h.MainDictionary.Plants[plot.PlantedPlant.PlantType].GrowthStages[plot.PlantedPlant.CurrentStage]
	}

	// Save to DBs
		
	saveFarmErr := schema.SaveFarmDataAtPathToDB(fdb, farmLocationSymbol, "plots", farm.Plots)
	if saveFarmErr != nil {
		log.Error.Printf("Error in PlotInteract, could not save farm. error: %v", saveFarmErr)
		responses.SendRes(w, responses.DB_Save_Failure, nil, saveFarmErr.Error())
		return
	}

	saveWarehouseErr := schema.SaveWarehouseToDB(wdb, &warehouse)
	if saveWarehouseErr != nil {
		log.Error.Printf("Error in PlotInteract, could not save warehouse. error: %v", saveWarehouseErr)
		responses.SendRes(w, responses.DB_Save_Failure, nil, saveWarehouseErr.Error())
		return
	}

	// Construct and Send response
	response := schema.PlotActionResponse{Warehouse: &warehouse, Plot: &plot, NextStage: nextStage}
	getPlotPlantResponseJsonString, getPlotPlantResponseJsonStringErr := responses.JSON(response)
	if getPlotPlantResponseJsonStringErr != nil {
		log.Error.Printf("Error in PlotInfo, could not format interact response as JSON. response: %v, error: %v", response, getPlotPlantResponseJsonStringErr)
		responses.SendRes(w, responses.JSON_Marshal_Error, response, getPlotPlantResponseJsonStringErr.Error())
		return
	}
	log.Debug.Printf("Sending response for CharterCaravan:\n%v", getPlotPlantResponseJsonString)
	responses.SendRes(w, responses.Generic_Success, response, "")
	
	log.Debug.Println(log.Cyan("-- End CharterCaravan --"))
}


// Handler function for the secure route: /api/my/locations
// Returns a list of locations 
type LocationsInfo struct {
	Dbs *map[string]rdb.Database
	World *schema.World
}
func (h *LocationsInfo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- LocationsInfo --"))
	udb := (*h.Dbs)["users"]
	OK, userData, _ := secureGetUser(w, r, udb)
	if !OK {
		return // Failure states handled by secureGetUser, simply return
	}
	// Get assistant locations to determine fog of war
	adb := (*h.Dbs)["assistants"]
	assistants, foundAssistants, assistantsErr := schema.GetAssistantsFromDB(userData.Assistants, adb)
	if assistantsErr != nil || !foundAssistants {
		log.Error.Printf("Error in LocationsInfo, could not get assistants from DB. foundAssistants: %v, error: %v", foundAssistants, assistantsErr)
		responses.SendRes(w, responses.DB_Get_Failure, assistants, assistantsErr.Error())
		return
	}
	// Get owned farm locations cause these always have vision
	fdb := (*h.Dbs)["farms"]
	farms, foundFarms, farmsErr := schema.GetFarmsFromDB(userData.Farms, fdb)
	if farmsErr != nil || !foundFarms {
		log.Error.Printf("Error in LocationsInfo, could not get farms from DB. foundFarms: %v, error: %v", foundFarms, farmsErr)
		responses.SendRes(w, responses.DB_Get_Failure, farms, farmsErr.Error())
		return
	}
	// use myLocs as a set to get all unique locations visible in fow
	myLocs := make(map[string]bool)
	for _, assistant := range assistants {
		myLocs[assistant.Location] = true
	}
	for _, farm := range farms {
		myLocs[farm.LocationSymbol] = true
	}
	// finally get all locations in each region
	resLocations := make([]schema.Location, 0)
	for location := range myLocs {
		resLocations = append(resLocations, h.World.Locations[location])
	}
	responses.SendRes(w, responses.Generic_Success, resLocations, "")
	log.Debug.Println(log.Cyan("-- End LocationsInfo --"))
}

// Handler function for the secure route: /api/my/nearby-locations
// Returns a list of locations 
type NearbyLocationsInfo struct {
	Dbs *map[string]rdb.Database
	World *schema.World
}
func (h *NearbyLocationsInfo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- NearbyLocationsInfo --"))
	udb := (*h.Dbs)["users"]
	OK, userData, _ := secureGetUser(w, r, udb)
	if !OK {
		return // Failure states handled by secureGetUser, simply return
	}
	// Get assistant locations to determine fog of war
	adb := (*h.Dbs)["assistants"]
	assistants, foundAssistants, assistantsErr := schema.GetAssistantsFromDB(userData.Assistants, adb)
	if assistantsErr != nil || !foundAssistants {
		log.Error.Printf("Error in NearbyLocationsInfo, could not get assistants from DB. foundAssistants: %v, error: %v", foundAssistants, assistantsErr)
		responses.SendRes(w, responses.DB_Get_Failure, assistants, assistantsErr.Error())
		return
	}
	// Get owned farm locations cause these always have vision
	fdb := (*h.Dbs)["farms"]
	farms, foundFarms, farmsErr := schema.GetFarmsFromDB(userData.Farms, fdb)
	if farmsErr != nil || !foundFarms {
		log.Error.Printf("Error in LocationsInfo, could not get farms from DB. foundFarms: %v, error: %v", foundFarms, farmsErr)
		responses.SendRes(w, responses.DB_Get_Failure, farms, farmsErr.Error())
		return
	}
	// use myLocs as a set
	myLocs := make(map[string]bool)
	for _, assistant := range assistants {
		myLocs[assistant.Location] = true
	}
	for _, farm := range farms {
		myLocs[farm.LocationSymbol] = true
	}
	// get every region revealed based on location using myRegions as a set
	myRegions := make(map[string]bool)
	for locSymb := range myLocs {
		lastInd := strings.LastIndex(locSymb, "-")
		regionSymbol := locSymb[:lastInd]
		myRegions[regionSymbol] = true
	}
	// finally get all locations in each region
	resLocations := make(map[string]map[string]string, 0)
	i := 0
	for regionSymbol := range myRegions {
		resLocations[regionSymbol] = make(map[string]string)
		for _, loc := range h.World.Locations {
			lastInd := strings.LastIndex(loc.Symbol, "-")
			locRegionSymb := loc.Symbol[:lastInd]
			if regionSymbol == locRegionSymb {
				resLocations[regionSymbol][loc.Symbol] = loc.Name
			}
		}
		i++
	}
	responses.SendRes(w, responses.Generic_Success, resLocations, "")
	log.Debug.Println(log.Cyan("-- End NearbyLocationsInfo --"))
}

// Handler function for the secure route: /api/my/locations/location-symbol
// Returns a locations 
type LocationInfo struct {
	Dbs *map[string]rdb.Database
	World *schema.World
}
func (h *LocationInfo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- LocationInfo --"))
	udb := (*h.Dbs)["users"]
	OK, userData, _ := secureGetUser(w, r, udb)
	if !OK {
		return // Failure states handled by secureGetUser, simply return
	}
	// Get assistant locations to determine fog of war
	adb := (*h.Dbs)["assistants"]
	assistants, foundAssistants, assistantsErr := schema.GetAssistantsFromDB(userData.Assistants, adb)
	if assistantsErr != nil || !foundAssistants {
		log.Error.Printf("Error in LocationInfo, could not get assistants from DB. foundAssistants: %v, error: %v", foundAssistants, assistantsErr)
		responses.SendRes(w, responses.DB_Get_Failure, assistants, assistantsErr.Error())
		return
	}
	// Get owned farm locations cause these always have vision
	fdb := (*h.Dbs)["farms"]
	farms, foundFarms, farmsErr := schema.GetFarmsFromDB(userData.Farms, fdb)
	if farmsErr != nil || !foundFarms {
		log.Error.Printf("Error in LocationsInfo, could not get farms from DB. foundFarms: %v, error: %v", foundFarms, farmsErr)
		responses.SendRes(w, responses.DB_Get_Failure, farms, farmsErr.Error())
		return
	}
	// use myLocs as a set to get all unique locations visible in fow
	myLocs := make(map[string]bool)
	for _, assistant := range assistants {
		myLocs[assistant.Location] = true
	}
	for _, farm := range farms {
		myLocs[farm.LocationSymbol] = true
	}
	// Get symbol from route
	symbol := GetVarEntries(r, "location-symbol", UUID)
	// finally get specified location if available
	var resLocation schema.Location
	found := false
	for location := range myLocs {
		if strings.ToUpper(location) == symbol {
			resLocation = h.World.Locations[location]
			found = true
		}
	}
	if !found {
		log.Debug.Printf("Not found %s in locations %v", symbol, myLocs)
		responses.SendRes(w, responses.No_Assitant_At_Location, nil, "")
		return
	}
	responses.SendRes(w, responses.Generic_Success, resLocation, "")
	log.Debug.Println(log.Cyan("-- End LocationInfo --"))
}

// Handler function for the secure route: /api/my/markets
// Returns a list of markets 
type MarketsInfo struct {
	Dbs *map[string]rdb.Database
	MainDictionary *schema.MainDictionary
}
func (h *MarketsInfo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- MarketsInfo --"))
	udb := (*h.Dbs)["users"]
	OK, userData, _ := secureGetUser(w, r, udb)
	if !OK {
		return // Failure states handled by secureGetUser, simply return
	}
	// Get assistant markets to determine fog of war
	adb := (*h.Dbs)["assistants"]
	assistants, foundAssistants, assistantsErr := schema.GetAssistantsFromDB(userData.Assistants, adb)
	if assistantsErr != nil || !foundAssistants {
		log.Error.Printf("Error in MarketsInfo, could not get assistants from DB. foundAssistants: %v, error: %v", foundAssistants, assistantsErr)
		responses.SendRes(w, responses.DB_Get_Failure, assistants, assistantsErr.Error())
		return
	}
	// use myLocs as a set to get all unique markets visible in fow
	myLocs := make(map[string]bool)
	for _, assistant := range assistants {
		myLocs[assistant.Location] = true
	}
	// finally get all markets in each region
	resMarkets := make([]schema.Market, 0)
	for market := range myLocs {
		resMarkets = append(resMarkets, h.MainDictionary.Markets[market])
	}
	responses.SendRes(w, responses.Generic_Success, resMarkets, "")
	log.Debug.Println(log.Cyan("-- End MarketsInfo --"))
}

// Handler function for the secure route: /api/my/markets/{symbol}
// Returns a markets 
type MarketInfo struct {
	Dbs *map[string]rdb.Database
	MainDictionary *schema.MainDictionary
}
func (h *MarketInfo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- MarketInfo --"))
	udb := (*h.Dbs)["users"]
	OK, userData, _ := secureGetUser(w, r, udb)
	if !OK {
		return // Failure states handled by secureGetUser, simply return
	}
	// Get assistant markets to determine fog of war
	adb := (*h.Dbs)["assistants"]
	assistants, foundAssistants, assistantsErr := schema.GetAssistantsFromDB(userData.Assistants, adb)
	if assistantsErr != nil || !foundAssistants {
		log.Error.Printf("Error in MarketInfo, could not get assistants from DB. foundAssistants: %v, error: %v", foundAssistants, assistantsErr)
		responses.SendRes(w, responses.DB_Get_Failure, assistants, assistantsErr.Error())
		return
	}
	// use myLocs as a set to get all unique markets visible in fow
	myLocs := make(map[string]bool)
	for _, assistant := range assistants {
		myLocs[assistant.Location] = true
	}
	// Get symbol from route
	symbol := GetVarEntries(r, "location-symbol", UUID)
	// finally get specified market if available
	var resMarket schema.Market
	found := false
	for market := range myLocs {
		if strings.ToUpper(market) == symbol {
			resMarket = h.MainDictionary.Markets[market]
			found = true
		}
	}
	if !found {
		log.Debug.Printf("Not found %s in markets %v", symbol, myLocs)
		responses.SendRes(w, responses.No_Assitant_At_Location, nil, "")
		return
	}
	responses.SendRes(w, responses.Generic_Success, resMarket, "")
	log.Debug.Println(log.Cyan("-- End MarketInfo --"))
}

// Handler function for the secure route: /api/my/farms
type FarmsInfo struct {
	Dbs *map[string]rdb.Database
}
func (h *FarmsInfo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- FarmsInfo --"))
	udb := (*h.Dbs)["users"]
	OK, userData, _ := secureGetUser(w, r, udb)
	if !OK {
		return // Failure states handled by secureGetUser, simply return
	}
	adb := (*h.Dbs)["farms"]
	farms, foundFarms, farmsErr := schema.GetFarmsFromDB(userData.Farms, adb)
	if farmsErr != nil || !foundFarms {
		log.Error.Printf("Error in FarmsInfo, could not get farms from DB. foundFarms: %v, error: %v", foundFarms, farmsErr)
		responses.SendRes(w, responses.DB_Get_Failure, farms, farmsErr.Error())
		return
	}
	getFarmJsonString, getFarmJsonStringErr := responses.JSON(farms)
	if getFarmJsonStringErr != nil {
		log.Error.Printf("Error in FarmsInfo, could not format farms as JSON. farms: %v, error: %v", farms, getFarmJsonStringErr)
		responses.SendRes(w, responses.JSON_Marshal_Error, farms, getFarmJsonStringErr.Error())
		return
	}
	log.Debug.Printf("Sending response for FarmsInfo:\n%v", getFarmJsonString)
	responses.SendRes(w, responses.Generic_Success, farms, "")
	log.Debug.Println(log.Cyan("-- End FarmsInfo --"))
}

// Handler function for the secure route: /api/my/farms/{uuid}
type FarmInfo struct {
	Dbs *map[string]rdb.Database
}
func (h *FarmInfo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- FarmInfo --"))
	// Get symbol from route
	symbol := GetVarEntries(r, "location-symbol", AllCaps)
	// Get userinfoContext from validation middleware
	userInfo, userInfoErr := GetValidationFromCtx(r)
	if userInfoErr != nil {
		// Fail state getting context
		log.Error.Printf("Could not get validationpair in AssistantInfo")
		userInfoErrMsg := fmt.Sprintf("userInfo is nil, check auth validation context %v:\n%v", auth.ValidationContext, r.Context().Value(auth.ValidationContext))
		responses.SendRes(w, responses.No_AuthPair_Context, nil, userInfoErrMsg)
		return
	}
	uuid := userInfo.Username + "|Farm-" + symbol
	log.Debug.Printf("FarmInfo Requested for: %s", uuid)
	adb := (*h.Dbs)["farms"]
	farm, foundFarm, farmsErr := schema.GetFarmFromDB(uuid, adb)
	if farmsErr != nil || !foundFarm {
		log.Error.Printf("Error in FarmInfo, could not get farm from DB. foundFarm: %v, error: %v", foundFarm, farmsErr)
		responses.SendRes(w, responses.DB_Get_Failure, farm, farmsErr.Error())
		return
	}
	getFarmJsonString, getFarmJsonStringErr := responses.JSON(farm)
	if getFarmJsonStringErr != nil {
		log.Error.Printf("Error in FarmInfo, could not format farms as JSON. farms: %v, error: %v", farm, getFarmJsonStringErr)
		responses.SendRes(w, responses.JSON_Marshal_Error, farm, getFarmJsonStringErr.Error())
		return
	}
	log.Debug.Printf("Sending response for FarmInfo:\n%v", getFarmJsonString)
	responses.SendRes(w, responses.Generic_Success, farm, "")
	log.Debug.Println(log.Cyan("-- End FarmInfo --"))
}

// Handler function for the secure route: /api/my/contracts
type ContractsInfo struct {
	Dbs *map[string]rdb.Database
}
func (h *ContractsInfo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- ContractsInfo --"))
	udb := (*h.Dbs)["users"]
	OK, userData, _ := secureGetUser(w, r, udb)
	if !OK {
		return // Failure states handled by secureGetUser, simply return
	}
	adb := (*h.Dbs)["contracts"]
	contracts, foundContracts, contractsErr := schema.GetContractsFromDB(userData.Contracts, adb)
	if contractsErr != nil || !foundContracts {
		log.Error.Printf("Error in ContractsInfo, could not get contracts from DB. foundContracts: %v, error: %v", foundContracts, contractsErr)
		responses.SendRes(w, responses.DB_Get_Failure, contracts, contractsErr.Error())
		return
	}
	getContractJsonString, getContractJsonStringErr := responses.JSON(contracts)
	if getContractJsonStringErr != nil {
		log.Error.Printf("Error in ContractsInfo, could not format contracts as JSON. contracts: %v, error: %v", contracts, getContractJsonStringErr)
		responses.SendRes(w, responses.JSON_Marshal_Error, contracts, getContractJsonStringErr.Error())
		return
	}
	log.Debug.Printf("Sending response for ContractsInfo:\n%v", getContractJsonString)
	responses.SendRes(w, responses.Generic_Success, contracts, "")
	log.Debug.Println(log.Cyan("-- End ContractsInfo --"))
}

// Handler function for the secure route: /api/my/contracts/{uuid}
type ContractInfo struct {
	Dbs *map[string]rdb.Database
}
func (h *ContractInfo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- ContractInfo --"))
	// Get symbol from route
	id := GetVarEntries(r, "contract-id", AllCaps)
	// Get userinfoContext from validation middleware
	userInfo, userInfoErr := GetValidationFromCtx(r)
	if userInfoErr != nil {
		// Fail state getting context
		log.Error.Printf("Could not get validationpair in AssistantInfo")
		userInfoErrMsg := fmt.Sprintf("userInfo is nil, check auth validation context %v:\n%v", auth.ValidationContext, r.Context().Value(auth.ValidationContext))
		responses.SendRes(w, responses.No_AuthPair_Context, nil, userInfoErrMsg)
		return
	}
	uuid := userInfo.Username + "|Contract-" + id
	log.Debug.Printf("ContractInfo Requested for: %s", uuid)
	adb := (*h.Dbs)["contracts"]
	contract, foundContract, contractsErr := schema.GetContractFromDB(uuid, adb)
	if contractsErr != nil || !foundContract {
		log.Error.Printf("Error in ContractInfo, could not get contract from DB. foundContract: %v, error: %v", foundContract, contractsErr)
		responses.SendRes(w, responses.DB_Get_Failure, contract, contractsErr.Error())
		return
	}
	getContractJsonString, getContractJsonStringErr := responses.JSON(contract)
	if getContractJsonStringErr != nil {
		log.Error.Printf("Error in ContractInfo, could not format contracts as JSON. contracts: %v, error: %v", contract, getContractJsonStringErr)
		responses.SendRes(w, responses.JSON_Marshal_Error, contract, getContractJsonStringErr.Error())
		return
	}
	log.Debug.Printf("Sending response for ContractInfo:\n%v", getContractJsonString)
	responses.SendRes(w, responses.Generic_Success, contract, "")
	log.Debug.Println(log.Cyan("-- End ContractInfo --"))
}

// Handler function for the secure route: /api/my/warehouses
type WarehousesInfo struct {
	Dbs *map[string]rdb.Database
}
func (h *WarehousesInfo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- WarehousesInfo --"))
	udb := (*h.Dbs)["users"]
	OK, userData, _ := secureGetUser(w, r, udb)
	if !OK {
		return // Failure states handled by secureGetUser, simply return
	}
	adb := (*h.Dbs)["warehouses"]
	warehouses, foundWarehouses, warehousesErr := schema.GetWarehousesFromDB(userData.Warehouses, adb)
	if warehousesErr != nil || !foundWarehouses {
		log.Error.Printf("Error in WarehousesInfo, could not get warehouses from DB. foundWarehouses: %v, error: %v", foundWarehouses, warehousesErr)
		responses.SendRes(w, responses.DB_Get_Failure, warehouses, warehousesErr.Error())
		return
	}
	getWarehousesJsonString, getWarehousesJsonStringErr := responses.JSON(warehouses)
	if getWarehousesJsonStringErr != nil {
		log.Error.Printf("Error in WarehousesInfo, could not format warehouses as JSON. warehouses: %v, error: %v", warehouses, getWarehousesJsonStringErr)
		responses.SendRes(w, responses.JSON_Marshal_Error, warehouses, getWarehousesJsonStringErr.Error())
		return
	}
	log.Debug.Printf("Sending response for WarehousesInfo:\n%v", getWarehousesJsonString)
	responses.SendRes(w, responses.Generic_Success, warehouses, "")
	log.Debug.Println(log.Cyan("-- End WarehousesInfo --"))
}

// Handler function for the secure route: /api/my/warehouses/{uuid}
type WarehouseInfo struct {
	Dbs *map[string]rdb.Database
}
func (h *WarehouseInfo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- WarehouseInfo --"))
	// Get symbol from route
	symbol := GetVarEntries(r, "location-symbol", AllCaps)
	// Get userinfoContext from validation middleware
	userInfo, userInfoErr := GetValidationFromCtx(r)
	if userInfoErr != nil {
		// Fail state getting context
		log.Error.Printf("Could not get validationpair in AssistantInfo")
		userInfoErrMsg := fmt.Sprintf("userInfo is nil, check auth validation context %v:\n%v", auth.ValidationContext, r.Context().Value(auth.ValidationContext))
		responses.SendRes(w, responses.No_AuthPair_Context, nil, userInfoErrMsg)
		return
	}
	uuid := userInfo.Username + "|Warehouse-" + symbol
	log.Debug.Printf("WarehouseInfo Requested for: %s", uuid)
	adb := (*h.Dbs)["warehouses"]
	warehouse, foundWarehouse, warehousesErr := schema.GetWarehouseFromDB(uuid, adb)
	if warehousesErr != nil || !foundWarehouse {
		log.Error.Printf("Error in WarehouseInfo, could not get warehouse from DB. foundWarehouse: %v, error: %v", foundWarehouse, warehousesErr)
		responses.SendRes(w, responses.DB_Get_Failure, warehouse, warehousesErr.Error())
		return
	}
	getWarehouseJsonString, getWarehouseJsonStringErr := responses.JSON(warehouse)
	if getWarehouseJsonStringErr != nil {
		log.Error.Printf("Error in WarehouseInfo, could not format warehouses as JSON. warehouses: %v, error: %v", warehouse, getWarehouseJsonStringErr)
		responses.SendRes(w, responses.JSON_Marshal_Error, warehouse, getWarehouseJsonStringErr.Error())
		return
	}
	log.Debug.Printf("Sending response for WarehouseInfo:\n%v", getWarehouseJsonString)
	responses.SendRes(w, responses.Generic_Success, warehouse, "")
	log.Debug.Println(log.Cyan("-- End WarehouseInfo --"))
}

// Handler function for the secure route: /api/my/plots
type PlotsInfo struct {
	Dbs *map[string]rdb.Database
}
func (h *PlotsInfo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- PlotsInfo --"))
	udb := (*h.Dbs)["users"]
	OK, userData, _ := secureGetUser(w, r, udb)
	if !OK {
		return // Failure states handled by secureGetUser, simply return
	}
	// Get farms from db
	adb := (*h.Dbs)["farms"]
	farms, foundFarms, farmsErr := schema.GetFarmsFromDB(userData.Farms, adb)
	if farmsErr != nil || !foundFarms {
		log.Error.Printf("Error in FarmsInfo, could not get farms from DB. foundFarms: %v, error: %v", foundFarms, farmsErr)
		responses.SendRes(w, responses.DB_Get_Failure, farms, farmsErr.Error())
		return
	}
	// Get plots from farms
	plots := make(map[string]schema.Plot, 0)
	for _, farm := range farms {
		for plotsymbol, plot := range farm.Plots {
			plots[plotsymbol] = plot
		}
	}
	getPlotsJsonString, getPlotsJsonStringErr := responses.JSON(plots)
	if getPlotsJsonStringErr != nil {
		log.Error.Printf("Error in PlotsInfo, could not format plots as JSON. plots: %v, error: %v", plots, getPlotsJsonStringErr)
		responses.SendRes(w, responses.JSON_Marshal_Error, plots, getPlotsJsonStringErr.Error())
		return
	}
	log.Debug.Printf("Sending response for PlotsInfo:\n%v", getPlotsJsonString)
	responses.SendRes(w, responses.Generic_Success, plots, "")
	log.Debug.Println(log.Cyan("-- End PlotsInfo --"))
}

// Handler function for the secure route: /api/my/plots/{uuid}
type PlotInfo struct {
	Dbs *map[string]rdb.Database
}
func (h *PlotInfo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- PlotInfo --"))
	// Get symbol from route
	id := GetVarEntries(r, "plot-id", None)
	// Get userinfoContext from validation middleware
	userInfo, userInfoErr := GetValidationFromCtx(r)
	if userInfoErr != nil {
		// Fail state getting context
		log.Error.Printf("Could not get validationpair in AssistantInfo")
		userInfoErrMsg := fmt.Sprintf("userInfo is nil, check auth validation context %v:\n%v", auth.ValidationContext, r.Context().Value(auth.ValidationContext))
		responses.SendRes(w, responses.No_AuthPair_Context, nil, userInfoErrMsg)
		return
	}
	idSlice := strings.Split(id, "!")
	if len(idSlice) < 2 {
		// Fail, malformed plot id
		errmsg := fmt.Sprintf("Malformed plot id, format must be '[farm-location-symbol]!Plot-[id-number]' received: %v", id)
		log.Debug.Printf(errmsg)
		responses.SendRes(w, responses.Bad_Request, nil, errmsg)
		return
	}
	uuid := userInfo.Username + "|Farm-" + idSlice[0] + "|" + idSlice[1]
	log.Debug.Printf("PlotInfo Requested for: %s", uuid)
	// Get farm
	symbolSlice := strings.Split(uuid, "|")
	farmSymbol := strings.Join(symbolSlice[:len(symbolSlice)-1], "|")
	adb := (*h.Dbs)["farms"]
	farm, foundFarm, farmsErr := schema.GetFarmFromDB(farmSymbol, adb)
	if farmsErr != nil || !foundFarm {
		log.Error.Printf("Error in FarmInfo, could not get farm from DB. foundFarm: %v, error: %v", foundFarm, farmsErr)
		responses.SendRes(w, responses.DB_Get_Failure, farm, farmsErr.Error())
		return
	}
	plot := farm.Plots[uuid]
	getPlotJsonString, getPlotJsonStringErr := responses.JSON(plot)
	if getPlotJsonStringErr != nil {
		log.Error.Printf("Error in PlotInfo, could not format plots as JSON. plots: %v, error: %v", plot, getPlotJsonStringErr)
		responses.SendRes(w, responses.JSON_Marshal_Error, plot, getPlotJsonStringErr.Error())
		return
	}
	log.Debug.Printf("Sending response for PlotInfo:\n%v", getPlotJsonString)
	responses.SendRes(w, responses.Generic_Success, plot, "")
	log.Debug.Println(log.Cyan("-- End PlotInfo --"))
}

// Handler function for the secure route: /api/my/plots/{uuid}/plant
type PlantPlot struct {
	Dbs *map[string]rdb.Database
	MainDictionary *schema.MainDictionary
}
func (h *PlantPlot) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- PlantPlot --"))
	// Get symbol from route
	id := GetVarEntries(r, "plot-id", None)
	// Get userinfoContext from validation middleware
	userInfo, userInfoErr := GetValidationFromCtx(r)
	if userInfoErr != nil {
		// Fail state getting context
		log.Error.Printf("Could not get validationpair in AssistantInfo")
		userInfoErrMsg := fmt.Sprintf("userInfo is nil, check auth validation context %v:\n%v", auth.ValidationContext, r.Context().Value(auth.ValidationContext))
		responses.SendRes(w, responses.No_AuthPair_Context, nil, userInfoErrMsg)
		return
	}
	idSlice := strings.Split(id, "!")
	if len(idSlice) < 2 {
		// Fail, malformed plot id
		errmsg := fmt.Sprintf("Malformed plot id, format must be '[farm-location-symbol]!Plot-[id-number]' received: %v", id)
		log.Debug.Printf(errmsg)
		responses.SendRes(w, responses.Bad_Request, nil, errmsg)
		return
	}
	uuid := userInfo.Username + "|Farm-" + idSlice[0] + "|" + idSlice[1]
	symbolSlice := strings.Split(uuid, "|")
	if len(symbolSlice) < 3 {
		// Fail, malformed plot id
		errmsg := fmt.Sprintf("Malformed plot id, symbolSlice less than 3: %v", symbolSlice)
		log.Debug.Printf(errmsg)
		responses.SendRes(w, responses.Bad_Request, nil, errmsg)
		return
	}
	farmLocationSymbol := strings.Join(symbolSlice[:len(symbolSlice)-1], "|")
	locationSymbol := strings.Split(symbolSlice[1], "-")
	warehouseLocationSymbol := symbolSlice[0] + "|Warehouse-" + strings.Join(locationSymbol[1:],"-")
	log.Debug.Printf("PlantPlot Requested for: %s", uuid)
	// unmarshall request body to get action and consumables if applicable
	var body schema.PlotPlantBody
	decoder := json.NewDecoder(r.Body)
	if decodeErr := decoder.Decode(&body); decodeErr != nil {
		// Fail case, could not decode
		errmsg := fmt.Sprintf("Decode Error in PlantPlot: %v", decodeErr)
		log.Debug.Printf(errmsg)
		responses.SendRes(w, responses.Bad_Request, nil, "Could not decode request body, ensure it conforms to expected format.")
		return
	}
	// Validate specified quantity is above 0
	if body.SeedQuantity <= 0 {
		// Fail, quantity must be > 0
		errmsg := fmt.Sprintf("in PlantPlot, SeedQuantity MUST be greater than 0. received seed quantity: %v", body.SeedQuantity)
		log.Debug.Printf(errmsg)
		responses.SendRes(w, responses.Bad_Request, nil, errmsg)
		return
	}
	// Validate specified seed is a seed
	seedsDict := (*h.MainDictionary).Seeds
	plantName, plantNameOk := seedsDict[body.SeedName]
	if !plantNameOk {
		// Fail, seed name specified does not match a known seed
		errmsg := fmt.Sprintf("in PlantPlot, SeedName does not map to seed in seed dictionary. received seed name: %v", body.SeedName)
		log.Debug.Printf(errmsg)
		responses.SendRes(w, responses.Item_Is_Not_Seed, nil, errmsg)
		return
	}
	// Get warehouses
	wdb := (*h.Dbs)["warehouses"]
	warehouse, foundWarehouse, warehousesErr := schema.GetWarehouseFromDB(warehouseLocationSymbol, wdb)
	if warehousesErr != nil || !foundWarehouse {
		errmsg := fmt.Sprintf("Error in PlantPlot, could not get warehouse from DB. foundWarehouse: %v, error: %v", foundWarehouse, warehousesErr)
		log.Error.Printf(errmsg)
		responses.SendRes(w, responses.DB_Get_Failure, nil, errmsg)
		return
	}
	// Validate seeds specified in given warehouse
	numLocalSeeds, ownedSeedsOk := warehouse.Seeds[body.SeedName]; 
	if !ownedSeedsOk {
		// Fail, seed good not in warehouse
		errmsg := fmt.Sprintf("Error in PlantPlot, Seed item not found in local warehouse. received good name: %v, warehouse goods: %v", body.SeedName, warehouse.Seeds)
		log.Error.Printf(errmsg)
		responses.SendRes(w, responses.Not_Enough_Items_In_Warehouse, nil, errmsg)
		return
	}
	if numLocalSeeds < uint64(body.SeedQuantity) {
		// Fail, not enough seeds in local inventory
		errmsg := fmt.Sprintf("Error in PlantPlot, not enough seed item found in local warehouse. received good name: %v, # local goods: %v", body.SeedName, numLocalSeeds)
		log.Error.Printf(errmsg)
		responses.SendRes(w, responses.Not_Enough_Items_In_Warehouse, map[string]interface{}{"name": body.SeedName, "number_in_warehouse": numLocalSeeds}, errmsg)
		return
	}
	// Get farms
	fdb := (*h.Dbs)["farms"]
	farm, foundFarm, farmsErr := schema.GetFarmFromDB(farmLocationSymbol, fdb)
	if farmsErr != nil || !foundFarm {
		log.Error.Printf("Error in PlantPlot, could not get farm from DB. foundFarm: %v, error: %v", foundFarm, farmsErr)
		responses.SendRes(w, responses.DB_Get_Failure, farm, farmsErr.Error())
		return
	}
	// Validate plot available for planting and body meets internal plot validation
	plot := farm.Plots[uuid]
	switch plot.IsPlantable(body) {
	case responses.Plot_Already_Planted:
		log.Error.Printf("Error in PlantPlot, plot already planted")
		responses.SendRes(w, responses.Plot_Already_Planted, plot, "")
		return
	case responses.Bad_Request:
		log.Error.Printf("Error in PlantPlot, seed size invalid")
		responses.SendRes(w, responses.Bad_Request, plot, "Seed Size specified in request body was invalid")
		return
	case responses.Plot_Too_Small:
		log.Error.Printf("Error in PlantPlot, plot too small")
		responses.SendRes(w, responses.Plot_Too_Small, plot, "")
		return
	case responses.Generic_Success:
		log.Debug.Printf("Plot ready for planting: %s", plot.UUID)
	default:
		log.Error.Fatalf("Received unexpected response type from plot.IsPlantable. plot: %v body: %v", plot, body)
	}

	plot.PlantedPlant = schema.NewPlant(plantName, body.SeedSize)
	plot.Quantity = body.SeedQuantity
	plot.GrowthCompleteTimestamp = time.Now().Unix()
	warehouse.RemoveSeeds(body.SeedName, uint64(body.SeedQuantity))
	farm.Plots[uuid] = plot

	// Save to DBs
	saveWarehouseErr := schema.SaveWarehouseDataAtPathToDB(wdb, warehouseLocationSymbol, "seeds", warehouse.Seeds)
	if saveWarehouseErr != nil {
		log.Error.Printf("Error in PlotInfo, could not save warehouse. error: %v", saveWarehouseErr)
		responses.SendRes(w, responses.DB_Save_Failure, nil, saveWarehouseErr.Error())
		return
	}
	saveFarmErr := schema.SaveFarmDataAtPathToDB(fdb, farmLocationSymbol, "plots", farm.Plots)
	if saveFarmErr != nil {
		log.Error.Printf("Error in PlantPlot, could not save farm. error: %v", saveFarmErr)
		responses.SendRes(w, responses.DB_Save_Failure, nil, saveFarmErr.Error())
		return
	}

	// Construct and Send response
	response := schema.PlotPlantResponse{Warehouse: &warehouse, Plot: &plot, NextStage: &h.MainDictionary.Plants[plantName].GrowthStages[plot.PlantedPlant.CurrentStage]}
	getPlotPlantResponseJsonString, getPlotPlantResponseJsonStringErr := responses.JSON(response)
	if getPlotPlantResponseJsonStringErr != nil {
		log.Error.Printf("Error in PlotInfo, could not format plant response as JSON. response: %v, error: %v", response, getPlotPlantResponseJsonStringErr)
		responses.SendRes(w, responses.JSON_Marshal_Error, response, getPlotPlantResponseJsonStringErr.Error())
		return
	}
	log.Debug.Printf("Sending response for PlantPlot:\n%v", getPlotPlantResponseJsonString)
	responses.SendRes(w, responses.Generic_Success, response, "")
	
	log.Debug.Println(log.Cyan("-- End PlantPlot --"))
}


// Handler function for the secure route: /api/my/plots/{uuid}/clear
type ClearPlot struct {
	Dbs *map[string]rdb.Database
}
func (h *ClearPlot) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- ClearPlot --"))
	// Get symbol from route
	id := GetVarEntries(r, "plot-id", None)
	// Get userinfoContext from validation middleware
	userInfo, userInfoErr := GetValidationFromCtx(r)
	if userInfoErr != nil {
		// Fail state getting context
		log.Error.Printf("Could not get validationpair in AssistantInfo")
		userInfoErrMsg := fmt.Sprintf("userInfo is nil, check auth validation context %v:\n%v", auth.ValidationContext, r.Context().Value(auth.ValidationContext))
		responses.SendRes(w, responses.No_AuthPair_Context, nil, userInfoErrMsg)
		return
	}
	idSlice := strings.Split(id, "!")
	if len(idSlice) < 2 {
		// Fail, malformed plot id
		errmsg := fmt.Sprintf("Malformed plot id, format must be '[farm-location-symbol]!Plot-[id-number]' received: %v", id)
		log.Debug.Printf(errmsg)
		responses.SendRes(w, responses.Bad_Request, nil, errmsg)
		return
	}
	uuid := userInfo.Username + "|Farm-" + idSlice[0] + "|" + idSlice[1]
	symbolSlice := strings.Split(uuid, "|")
	if len(symbolSlice) < 3 {
		// Fail, malformed plot id
		errmsg := fmt.Sprintf("Malformed plot id, symbolSlice less than 3: %v", symbolSlice)
		log.Debug.Printf(errmsg)
		responses.SendRes(w, responses.Bad_Request, nil, errmsg)
		return
	}
	farmLocationSymbol := strings.Join(symbolSlice[:len(symbolSlice)-1], "|")
	log.Debug.Printf("ClearPlot Requested for: %s", uuid)

	// Get farms
	fdb := (*h.Dbs)["farms"]
	farm, foundFarm, farmsErr := schema.GetFarmFromDB(farmLocationSymbol, fdb)
	if farmsErr != nil || !foundFarm {
		log.Error.Printf("Error in ClearPlot, could not get farm from DB. foundFarm: %v, error: %v", foundFarm, farmsErr)
		responses.SendRes(w, responses.DB_Get_Failure, farm, farmsErr.Error())
		return
	}

	// Validate plot available for clearing and body meets internal plot validation
	plot := farm.Plots[uuid]
	if plot.PlantedPlant == nil {
		log.Error.Printf("Error in ClearPlot, plot already empty")
		responses.SendRes(w, responses.Plot_Already_Empty, plot, "")
		return
	}

	plot.PlantedPlant = nil
	plot.Quantity = 0
	plot.GrowthCompleteTimestamp = time.Now().Unix()
	farm.Plots[uuid] = plot

	// Save to DB
	saveFarmErr := schema.SaveFarmDataAtPathToDB(fdb, farmLocationSymbol, "plots", farm.Plots)
	if saveFarmErr != nil {
		log.Error.Printf("Error in ClearPlot, could not save farm. error: %v", saveFarmErr)
		responses.SendRes(w, responses.DB_Save_Failure, nil, saveFarmErr.Error())
		return
	}

	// Construct and Send response
	resmsg := fmt.Sprintf("Successfully cleared plot: %s", uuid)
	log.Debug.Printf("Sending successful response for PlantPlot",)
	responses.SendRes(w, responses.Generic_Success, plot, resmsg)
	
	log.Debug.Println(log.Cyan("-- End PlantPlot --"))
}

// Handler function for the secure route: /api/my/plots/{uuid}/interact
type InteractPlot struct {
	Dbs *map[string]rdb.Database
	MainDictionary *schema.MainDictionary
}
func (h *InteractPlot) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- InteractPlot --"))
	// Get symbol from route
	id := GetVarEntries(r, "plot-id", None)
	// Get userinfoContext from validation middleware
	userInfo, userInfoErr := GetValidationFromCtx(r)
	if userInfoErr != nil {
		// Fail state getting context
		log.Error.Printf("Could not get validationpair in AssistantInfo")
		userInfoErrMsg := fmt.Sprintf("userInfo is nil, check auth validation context %v:\n%v", auth.ValidationContext, r.Context().Value(auth.ValidationContext))
		responses.SendRes(w, responses.No_AuthPair_Context, nil, userInfoErrMsg)
		return
	}
	idSlice := strings.Split(id, "!")
	if len(idSlice) < 2 {
		// Fail, malformed plot id
		errmsg := fmt.Sprintf("Malformed plot id, format must be '[farm-location-symbol]!Plot-[id-number]' received: %v", id)
		log.Debug.Printf(errmsg)
		responses.SendRes(w, responses.Bad_Request, nil, errmsg)
		return
	}
	uuid := userInfo.Username + "|Farm-" + idSlice[0] + "|" + idSlice[1]
	symbolSlice := strings.Split(uuid, "|")
	if len(symbolSlice) < 3 {
		// Fail, malformed plot id
		errmsg := fmt.Sprintf("Malformed plot id, symbolSlice less than 3: %v", symbolSlice)
		log.Debug.Printf(errmsg)
		responses.SendRes(w, responses.Bad_Request, nil, errmsg)
		return
	}
	farmLocationSymbol := strings.Join(symbolSlice[:len(symbolSlice)-1], "|")
	log.Debug.Printf("InteractPlot Requested for: %s", uuid)

	// Get farm, plot
	fdb := (*h.Dbs)["farms"]
	farm, foundFarm, farmsErr := schema.GetFarmFromDB(farmLocationSymbol, fdb)
	if farmsErr != nil || !foundFarm {
		log.Error.Printf("Error in InteractPlot, could not get farm from DB. foundFarm: %v, error: %v", foundFarm, farmsErr)
		responses.SendRes(w, responses.DB_Get_Failure, farm, farmsErr.Error())
		return
	}
	plot := farm.Plots[uuid]

	// Validate Planted
	if plot.PlantedPlant == nil {
		// no plant, cannot interact
		log.Debug.Printf("in InteractPlot, plot not planted. foundPlot: %v", plot)
		responses.SendRes(w, responses.Plot_Not_Planted, plot, "")
		return
	}

	// Validate Timestamp
	if plot.GrowthCompleteTimestamp > time.Now().Unix() {
		// too soon, reject
		timestampMsg := fmt.Sprintf("Ready in %d seconds", plot.GrowthCompleteTimestamp - time.Now().Unix())
		responses.SendRes(w, responses.Plants_Still_Growing, plot, timestampMsg)
		return
	}

	// unmarshall request body to get action and consumables if applicable
	var body schema.PlotInteractBody
	decoder := json.NewDecoder(r.Body)
	if decodeErr := decoder.Decode(&body); decodeErr != nil {
		// Fail case, could not decode
		errmsg := fmt.Sprintf("Decode Error in InteractPlot: %v", decodeErr)
		log.Debug.Printf(errmsg)
		responses.SendRes(w, responses.Bad_Request, nil, "Could not decode request body, ensure it conforms to expected format.")
		return
	}

	// Get warehouse
	wdb := (*h.Dbs)["warehouses"]
	locationSymbol := strings.Split(symbolSlice[1], "-")
	warehouseLocationSymbol := symbolSlice[0] + "|Warehouse-" + strings.Join(locationSymbol[1:],"-")
	warehouse, foundWarehouse, warehousesErr := schema.GetWarehouseFromDB(warehouseLocationSymbol, wdb)
	if warehousesErr != nil || !foundWarehouse {
		errmsg := fmt.Sprintf("Error in InteractPlot, could not get warehouse from DB. foundWarehouse: %v, error: %v", foundWarehouse, warehousesErr)
		log.Error.Printf(errmsg)
		responses.SendRes(w, responses.DB_Get_Failure, nil, errmsg)
		return
	}

	consumableQuantityAvailable := uint64(0)
	consumableName := strings.Title(strings.ToLower(body.Consumable))
	// If consumables included, validate them
	if consumableName != string("") {
		// Validate specified consumable is a good
		goodsDict := (*h.MainDictionary).Goods
		if _, ok := goodsDict[consumableName]; !ok {
			// Fail, seed is not good
			errmsg := fmt.Sprintf("in InteractPlot, consumable item does not exist in good dictionary. received consumable name: %v", consumableName)
			log.Debug.Printf(errmsg)
			responses.SendRes(w, responses.Item_Does_Not_Exist, nil, errmsg)
			return
		}
		// Validate consumable specified in given warehouse
		temp, ownedConsumableOk := warehouse.Goods[consumableName]; 
		if !ownedConsumableOk {
			// Fail, consumable good not in warehouse
			errmsg := fmt.Sprintf("in InteractPlot, consumable item not found in local warehouse. received good name: %v, warehouse goods: %v", consumableName, warehouse.Goods)
			log.Debug.Printf(errmsg)
			responses.SendRes(w, responses.Not_Enough_Items_In_Warehouse, nil, errmsg)
			return
		}
		consumableQuantityAvailable = temp
	}

	// Validate plot available for interaction and body meets internal plot validation
	plantDef, plantDefOk := h.MainDictionary.Plants[plot.PlantedPlant.PlantType]

	if !plantDefOk {
		plotDefErrMsg := fmt.Sprintf("Error in InteractPlot, Plot [%s] has PlantedPlant type [%s] not found in main dictionary!", plot.UUID, plot.PlantedPlant.PlantType)
		log.Error.Println(plotDefErrMsg)
		responses.SendRes(w, responses.Internal_Server_Error, plot, plotDefErrMsg)
	}
	plotValidationResponse, addedYield, usedConsumableQuantity, growthHarvest, growthTime, errInfoMsg := plot.IsInteractable(body, plantDef, consumableQuantityAvailable, warehouse.Tools)
	switch plotValidationResponse {
	case responses.Invalid_Plot_Action:
		log.Debug.Printf("in PlotInteract, Invalid_Plot_Action")
		responses.SendRes(w, responses.Invalid_Plot_Action, plot, errInfoMsg)
		return
	case responses.Tool_Not_Found:
		log.Debug.Printf("in PlotInteract, Tool_Not_Found")
		responses.SendRes(w, responses.Tool_Not_Found, plot, errInfoMsg)
		return
	case responses.Missing_Consumable_Selection:
		log.Debug.Printf("in PlotInteract, Missing_Consumable_Selection")
		responses.SendRes(w, responses.Missing_Consumable_Selection, plot, errInfoMsg)
		return
	case responses.Internal_Server_Error:
		log.Debug.Printf("in PlotInteract, Internal_Server_Error")
		responses.SendRes(w, responses.Internal_Server_Error, plot, errInfoMsg)
		return
	case responses.Not_Enough_Items_In_Warehouse:
		log.Debug.Printf("in PlotInteract, Not_Enough_Items_In_Warehouse")
		responses.SendRes(w, responses.Not_Enough_Items_In_Warehouse, plot, errInfoMsg)
		return
	case responses.Consumable_Not_In_Options:
		log.Debug.Printf("in PlotInteract, Consumable_Not_In_Options")
		responses.SendRes(w, responses.Consumable_Not_In_Options, plot, errInfoMsg)
		return
	case responses.Generic_Success:
		log.Debug.Printf("Plot growth action validated successfully: %s, action: %s", plot.UUID, body.Action)
	default:
		log.Error.Fatalf("Received unexpected response type from plot.IsPlantable. plot: %v body: %v", plot, body)
	}

	// Update objects with results of interaction

	plot.PlantedPlant.Yield += addedYield
	log.Debug.Printf("Interact Plot, Growth Time: %d", growthTime)
	plot.GrowthCompleteTimestamp = time.Now().Unix() + growthTime
	farm.Plots[uuid] = plot

	if consumableName != string("") {
		// if consumables used
		warehouse.RemoveGoods(consumableName, usedConsumableQuantity)
	}


	// Handle updating stage and potential harvesting
	var nextStage *schema.GrowthStage

	if growthHarvest != nil {
		// if was a harvest action
		harvest := plot.CalculateProduce(growthHarvest)
		log.Debug.Println("Harvest Calculated:")
		log.Debug.Println(harvest)
		for _, produce := range harvest.Produce {
			log.Debug.Printf("Add produce quantity: %d", produce.Quantity)
			warehouse.AddProduce(produce.Name, produce.Size, produce.Quantity)
			log.Debug.Println(warehouse.Produce)
		}
		for seedname, seedquantity := range harvest.Seeds {
			log.Debug.Printf("Add seed quantity: %d", seedquantity)
			warehouse.AddSeeds(seedname, seedquantity)
		}
		for goodname, goodquantity := range harvest.Goods {
			log.Debug.Printf("Add good quantity: %d", goodquantity)
			warehouse.AddGoods(goodname, goodquantity)
		}

		metrics.TrackHarvest(plantDef.Name)
		
		// check if final harvest
		if growthHarvest.FinalHarvest {
			// is final harvest
			// clear
			plot.PlantedPlant = nil
			plot.Quantity = 0
			farm.Plots[uuid] = plot
		} else {
			// not final harvest
			// move up current stage, initialize nextStage for response
			plot.PlantedPlant.CurrentStage++
			nextStage = &h.MainDictionary.Plants[plot.PlantedPlant.PlantType].GrowthStages[plot.PlantedPlant.CurrentStage]
		}
	} else {
		// if not harvest
		// move up current stage, initialize nextStage for response
		plot.PlantedPlant.CurrentStage++
		nextStage = &h.MainDictionary.Plants[plot.PlantedPlant.PlantType].GrowthStages[plot.PlantedPlant.CurrentStage]
	}

	// Save to DBs
		
	saveFarmErr := schema.SaveFarmDataAtPathToDB(fdb, farmLocationSymbol, "plots", farm.Plots)
	if saveFarmErr != nil {
		log.Error.Printf("Error in PlotInteract, could not save farm. error: %v", saveFarmErr)
		responses.SendRes(w, responses.DB_Save_Failure, nil, saveFarmErr.Error())
		return
	}

	saveWarehouseErr := schema.SaveWarehouseToDB(wdb, &warehouse)
	if saveWarehouseErr != nil {
		log.Error.Printf("Error in PlotInteract, could not save warehouse. error: %v", saveWarehouseErr)
		responses.SendRes(w, responses.DB_Save_Failure, nil, saveWarehouseErr.Error())
		return
	}

	// Construct and Send response
	response := schema.PlotActionResponse{Warehouse: &warehouse, Plot: &plot, NextStage: nextStage}
	getPlotPlantResponseJsonString, getPlotPlantResponseJsonStringErr := responses.JSON(response)
	if getPlotPlantResponseJsonStringErr != nil {
		log.Error.Printf("Error in PlotInfo, could not format interact response as JSON. response: %v, error: %v", response, getPlotPlantResponseJsonStringErr)
		responses.SendRes(w, responses.JSON_Marshal_Error, response, getPlotPlantResponseJsonStringErr.Error())
		return
	}
	log.Debug.Printf("Sending response for InteractPlot:\n%v", getPlotPlantResponseJsonString)
	responses.SendRes(w, responses.Generic_Success, response, "")
	
	log.Debug.Println(log.Cyan("-- End InteractPlot --"))
}

// Handler function for the secure route: /api/my/markets/{symbol}/order
// Returns a list of markets 
type MarketOrder struct {
	Dbs *map[string]rdb.Database
	MainDictionary *schema.MainDictionary
}
func (h *MarketOrder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- MarketOrder --"))
	udb := (*h.Dbs)["users"]
	OK, userData, _ := secureGetUser(w, r, udb)
	if !OK {
		return // Failure states handled by secureGetUser, simply return
	}
	// Get assistant markets to determine fog of war
	adb := (*h.Dbs)["assistants"]
	assistants, foundAssistants, assistantsErr := schema.GetAssistantsFromDB(userData.Assistants, adb)
	if assistantsErr != nil || !foundAssistants {
		log.Error.Printf("Error in MarketOrder, could not get assistants from DB. foundAssistants: %v, error: %v", foundAssistants, assistantsErr)
		responses.SendRes(w, responses.DB_Get_Failure, assistants, assistantsErr.Error())
		return
	}
	// use myLocs as a set to get all unique markets visible in fow
	myLocs := make(map[string]bool)
	for _, assistant := range assistants {
		myLocs[assistant.Location] = true
	}
	// Get symbol from route
	symbol := GetVarEntries(r, "location-symbol", UUID)
	// finally get specified market if available
	var resMarket schema.Market
	found := false
	for market := range myLocs {
		if strings.ToUpper(market) == symbol {
			resMarket = h.MainDictionary.Markets[market]
			found = true
		}
	}
	if !found {
		log.Debug.Printf("Not found %s in markets %v", symbol, myLocs)
		responses.SendRes(w, responses.No_Assitant_At_Location, nil, "")
		return
	}

	// NOW handle setting up the order and executing it if type is MARKET order
	// unmarshall request order to get action and consumables if applicable
	var order schema.MarketOrder
	decoder := json.NewDecoder(r.Body)
	if decodeErr := decoder.Decode(&order); decodeErr != nil {
		// Fail case, could not decode
		errmsg := fmt.Sprintf("Decode Error in MarketOrder: %v", decodeErr)
		log.Debug.Printf(errmsg)
		responses.SendRes(w, responses.Bad_Request, nil, "Could not decode request body, ensure it conforms to expected format.")
		return
	}

	// Get warehouse
	wdb := (*h.Dbs)["warehouses"]
	warehouseLocationSymbol := userData.Username + "|Warehouse-" + symbol
	warehouse, foundWarehouse, warehousesErr := schema.GetWarehouseFromDB(warehouseLocationSymbol, wdb)
	if warehousesErr != nil || !foundWarehouse {
		errmsg := fmt.Sprintf("Error in MarketOrder, could not get warehouse from DB. foundWarehouse: %v, error: %v", foundWarehouse, warehousesErr)
		log.Error.Printf(errmsg)
		responses.SendRes(w, responses.DB_Get_Failure, nil, errmsg)
		return
	}

	// Validate order parameters
	if order.Quantity <= 0 {
		// fail, quantity too low
		errmsg := fmt.Sprintf("in MarketOrder, invalid order quantity, must be > 0, got %d.", order.Quantity)
		log.Debug.Printf(errmsg)
		responses.SendRes(w, responses.Bad_Request, nil, errmsg)
		return
	}
	// Find order item in specified market imports/exports
	var ioField schema.MarketIOField
	if order.TXType == schema.BUY {
		ioField = resMarket.Exports
	} else {
		ioField = resMarket.Imports
	}
	var itemDict map[string]uint64
	var warehouseDict map[string]uint64
	var itemName string
	var simpleItemName string
	var sizeMod uint64
	switch order.ItemCategory {
	case schema.GOOD:
		itemDict = ioField.Goods
		warehouseDict = warehouse.Goods
		itemName = order.ItemName
		simpleItemName = order.ItemName
		sizeMod = 1
	case schema.SEED:
		itemDict = ioField.Seeds
		warehouseDict = warehouse.Seeds
		itemName = order.ItemName
		simpleItemName = order.ItemName
		sizeMod = 1
	case schema.TOOL:
		itemDict = ioField.Tools
		warehouseDict = warehouse.Tools
		itemName = order.ItemName
		simpleItemName = order.ItemName
		sizeMod = 1
	case schema.PRODUCE:
		itemDict = ioField.Produce
		warehouseDict = warehouse.GetSimpleProduceDict()
		splitSlice := strings.Split(order.ItemName, "|")
		if len(splitSlice) <= 1 {
			// fail, no size specified
			errmsg := "in MarketOrder, order.ItemName for PRODUCE category MUST have size specified e.g. 'Potato|Large', ensure correct ItemCategory specified and item name contains size"
			log.Debug.Printf("%s, %s: %v", errmsg, itemName, splitSlice)
			responses.SendRes(w, responses.Market_Order_Failed_Validation, nil, errmsg)
			return
		}
		size, ok := schema.SizeToID[strings.Title(strings.ToLower(splitSlice[1]))]
		if !ok {
			// fail nonsensical size
			errmsg := fmt.Sprintf("in MarketOrder, order.ItemName for PRODUCE category MUST have valid size specified e.g. 'Potato|Large', ensure correct ItemCategory specified and item name contains a valid size (received: %s)", splitSlice[1])
			log.Debug.Printf("%s, %s: %v", errmsg, itemName, splitSlice)
			responses.SendRes(w, responses.Market_Order_Failed_Validation, nil, errmsg)
			return
		}
		itemName = splitSlice[0] + "|" + strings.Title(strings.ToLower(splitSlice[1]))
		simpleItemName = splitSlice[0]
		sizeMod = uint64(size)
	}

	// get market value
	marketValue, mvOk := itemDict[simpleItemName]
	if !mvOk {
		// fail, item not in specified market list
		errmsg := "in MarketOrder, order.ItemName not in itemDict, ensure correct ItemCategory specified and item is available for trade in specified market for given transaction type"
		log.Debug.Printf("%s, %s: %v", errmsg, simpleItemName, itemDict)
		responses.SendRes(w, responses.Market_Order_Failed_Validation, nil, errmsg)
		return
	}

	// execute buy or sell is have enough in warehouse/ledger
	log.Debug.Printf("Execute Market Order: %s %s x%d for %d each * %d sizeMod", order.OrderType, itemName, order.Quantity, itemDict[simpleItemName], sizeMod)
	coins := userData.Ledger.Currencies["Coins"]
	if order.TXType == schema.BUY {
		orderCost := order.Quantity * marketValue * sizeMod
		// Validate currency in ledger in sufficient quantity
		if orderCost > coins {
			// fail, not enough currency for specified item and quantity
			errmsg := fmt.Sprintf("in MarketOrder, order cost %d > coins: %d", orderCost, coins)
			log.Debug.Printf(errmsg)
			responses.SendRes(w, responses.Market_Order_Failed_Validation, nil, errmsg)
			return
		}
		// Execute buy
		coins -= orderCost
		warehouseDict[itemName] += order.Quantity
		metrics.TrackMarketBuySell(itemName, true, order.Quantity)
	} else {
		orderProfit := order.Quantity * marketValue * sizeMod
		// Validate in warehouse in sufficient quantity
		warehouseQuantity, wqOk := warehouseDict[itemName]
		if !wqOk {
			// fail, item not in location's warehouse
			errmsg := "in MarketOrder, specified item not found in local warehouse, ensure order.ItemCategory is specified correctly"
			log.Debug.Printf("%s, %s: %v", errmsg, itemName, warehouseDict)
			responses.SendRes(w, responses.Market_Order_Failed_Validation, nil, errmsg)
			return
		}
		if warehouseQuantity < order.Quantity {
			// fail, not enough in location's warehouse
			errmsg := fmt.Sprintf("in MarketOrder, order specifies greater quantity %d than available in local warehouse %d.", order.Quantity, warehouseQuantity)
			log.Debug.Printf(errmsg)
			responses.SendRes(w, responses.Market_Order_Failed_Validation, nil, errmsg)
			return
		}
		// Execute sell
		coins += orderProfit
		warehouseDict[itemName] -= order.Quantity
		// Handle deleting if necessary
		if warehouseDict[itemName] <= 0 {
			delete(warehouseDict, itemName)
		}
		metrics.TrackMarketBuySell(itemName, false, order.Quantity)
	}
	
	// Apply results to original objects
	switch order.ItemCategory {
	case schema.GOOD:
		warehouse.Goods = warehouseDict
	case schema.SEED:
		warehouse.Goods = warehouseDict
	case schema.TOOL:
		warehouse.Goods = warehouseDict
	case schema.PRODUCE:
		warehouse.SetSimpleProduceDict(warehouseDict)
	}
	userData.Ledger.Currencies["Coins"] = coins

	// Save warehouse
	saveWarehouseErr := schema.SaveWarehouseToDB(wdb, &warehouse)
	if saveWarehouseErr != nil {
		log.Error.Printf("Error in MarketOrder, could not save warehouse. error: %v", saveWarehouseErr)
		responses.SendRes(w, responses.DB_Save_Failure, nil, saveWarehouseErr.Error())
		return
	}

	// Save user
	saveUserErr := schema.SaveUserToDB((*h.Dbs)["users"], &userData)
	if saveUserErr != nil {
		log.Error.Printf("Error in MarketOrder, could not save user. error: %v", saveUserErr)
		responses.SendRes(w, responses.DB_Save_Failure, nil, saveUserErr.Error())
		return
	}

	responses.SendRes(w, responses.Generic_Success, map[string]interface{}{"warehouse": warehouse, "ledger": userData.Ledger}, "")
	log.Debug.Println(log.Cyan("-- End MarketOrder --"))
}