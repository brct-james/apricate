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
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// HELPER FUNCTIONS

// Attempt to get validation context
func GetValidationFromCtx(r *http.Request) (auth.ValidationPair, error) {
	log.Debug.Println("Recover validationpair from context")
	userInfo, ok := r.Context().Value(auth.ValidationContext).(auth.ValidationPair)
	if !ok {
		return auth.ValidationPair{}, errors.New("could not get ValidationPair")
	}
	return userInfo, nil
}

// Get User from Middleware and DB
// Returns: OK, userData, userAuthPair
func secureGetUser(w http.ResponseWriter, r *http.Request, udb rdb.Database) (bool, schema.User, auth.ValidationPair) {
	// Get userinfoContext from validation middleware
	userInfo, userInfoErr := GetValidationFromCtx(r)
	if userInfoErr != nil {
		// Fail state getting context
		log.Error.Printf("Could not get validationpair in secureGetUser")
		userInfoErrMsg := fmt.Sprintf("userInfo is nil, check auth validation context %v:\n%v", auth.ValidationContext, r.Context().Value(auth.ValidationContext))
		responses.SendRes(w, responses.No_AuthPair_Context, nil, userInfoErrMsg)
		return false, schema.User{}, userInfo
	}
	log.Debug.Printf("Validated with username: %s and token %s", userInfo.Username, userInfo.Token)
	// Check db for user
	thisUser, userFound, getUserErr := schema.GetUserFromDB(userInfo.Token, udb)
	if getUserErr != nil {
		// fail state
		getErrorMsg := fmt.Sprintf("in secureGetUser, could not get from DB for username: %s, error: %v", userInfo.Username, getUserErr)
		responses.SendRes(w, responses.UDB_Get_Failure, nil, getErrorMsg)
		return false, schema.User{}, auth.ValidationPair{}
	}
	if !userFound {
		// fail state - user not found
		userNotFoundMsg := fmt.Sprintf("in secureGetUser, no user found in DB with username: %s", userInfo.Username)
		responses.SendRes(w, responses.User_Not_Found, nil, userNotFoundMsg)
		return false, schema.User{}, auth.ValidationPair{}
	}

	// Any time a user hits a secure endpoint, track a call from their account
	metrics.TrackUserCall(userInfo.Username)

	// // Get wdb
	// wdbSuccess, wdb := GetWdbFromCtx(w, r)
	// if !wdbSuccess {
	// 	log.Debug.Printf("Could not get wdb from ctx")
	// 	return false, schema.User{}, auth.ValidationPair{} // Fail state, could not get wdb, handled by func - simply return
	// }
	// // Success state, got wdb

	// // Success case
	// thisUser, calcErr := gamelogic.CalculateUserUpdates(thisUser, wdb)
	// if calcErr != nil {
	// 	// Fail state could not calculate user updates
	// 	resMsg := fmt.Sprintf("calcErr: %v", calcErr)
	// 	responses.SendRes(w, responses.Generic_Failure, nil, resMsg)
	// 	return false, thisUser, userInfo
	// }

	// // Lastly, GetGolemMapWithPublicInfo
	// thisUser.Golems = schema.UpdateGolemMapLinkedData(thisUser, thisUser.Golems) 
	return true, thisUser, userInfo
}

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
		log.Error.Printf("Error in AccountInfo, could not format thisUser as JSON. userData: %v, error: %v", userData, getUserJsonStringErr)
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
	// Get uuid from route
	route_vars := mux.Vars(r)
	uuid := route_vars["uuid"]
	log.Debug.Printf("AssistantInfo Requested for: %s", uuid)
	adb := (*h.Dbs)["assistants"]
	assistant, foundAssistant, assistantsErr := schema.GetAssistantFromDB(uuid, adb)
	if assistantsErr != nil || !foundAssistant {
		log.Error.Printf("Error in AssistantInfo, could not get assistant from DB. foundAssistant: %v, error: %v", foundAssistant, assistantsErr)
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
		myLocs[assistant.LocationSymbol] = true
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
		myLocs[assistant.LocationSymbol] = true
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

// Handler function for the secure route: /api/my/locations
// Returns a list of locations 
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
		myLocs[assistant.LocationSymbol] = true
	}
	for _, farm := range farms {
		myLocs[farm.LocationSymbol] = true
	}
	// Get symbol from route
	route_vars := mux.Vars(r)
	symbol := strings.ToUpper(route_vars["symbol"])
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
	// Get uuid from route
	route_vars := mux.Vars(r)
	uuid := route_vars["uuid"]
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
	// Get uuid from route
	route_vars := mux.Vars(r)
	uuid := route_vars["uuid"]
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
	// Get uuid from route
	route_vars := mux.Vars(r)
	uuid := route_vars["uuid"]
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
	// Get uuid from route
	route_vars := mux.Vars(r)
	uuid := route_vars["uuid"]
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
	// Get uuid from route
	route_vars := mux.Vars(r)
	uuid := route_vars["uuid"]
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
	// Validate specified seed is a good
	goodsDict := (*h.MainDictionary).Goods
	if _, ok := goodsDict[body.SeedName]; !ok {
		// Fail, seed is not good
		errmsg := fmt.Sprintf("Error in PlantPlot, Seed item does not exist in good dictionary. received seed name: %v", body.SeedName)
		log.Error.Printf(errmsg)
		responses.SendRes(w, responses.Item_Does_Not_Exist, nil, errmsg)
		return
	}
	// Validate specified seed is a seed
	seedsDict := (*h.MainDictionary).Seeds
	plantName, plantNameOk := seedsDict[body.SeedName]
	if !plantNameOk {
		// Fail, seed name specified does not match a known seed
		errmsg := fmt.Sprintf("Error in PlantPlot, SeedName does not map to seed in seed dictionary. received seed name: %v", body.SeedName)
		log.Error.Printf(errmsg)
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
	numLocalSeeds, ownedSeedsOk := warehouse.Goods[body.SeedName]; 
	if !ownedSeedsOk {
		// Fail, seed good not in warehouse
		errmsg := fmt.Sprintf("Error in PlantPlot, Seed item not found in local warehouse. received good name: %v, warehouse goods: %v", body.SeedName, warehouse.Goods)
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
	case responses.Plot_Too_Small:
		log.Error.Printf("Error in PlantPlot, plot too small")
		responses.SendRes(w, responses.Plot_Too_Small, plot, "")
		return
	case responses.Generic_Success:
		log.Debug.Printf("Plot ready for planting: %s", plot.UUID)
	default:
		log.Error.Fatalf("Received unexpected response type from plot.IsPlantable. plot: %v body: %v", plot, body)
	}

	plot.PlantedPlant = schema.NewPlant(schema.PlantTypeFromString(plantName), body.SeedSize)
	plot.Quantity = body.SeedQuantity
	warehouse.RemoveGoods(body.SeedName, uint64(body.SeedQuantity))
	farm.Plots[uuid] = plot

	// Save to DBs
	saveWarehouseErr := schema.SaveWarehouseDataAtPathToDB(wdb, warehouseLocationSymbol, "goods", warehouse.Goods)
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
	// Get uuid from route
	route_vars := mux.Vars(r)
	uuid := route_vars["uuid"]
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
	// Get uuid from route
	route_vars := mux.Vars(r)
	uuid := route_vars["uuid"]
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
		log.Error.Printf("Error in InteractPlot, plot not planted. foundPlot: %v", plot)
		responses.SendRes(w, responses.Plot_Not_Planted, plot, "")
		return
	}

	// // Validate Timestamp
	// if plot.GrowthCompleteTimestamp > time.Now().Unix() {
	// 	// too soon, reject
	// 	timestampMsg := fmt.Sprintf("Ready in %d seconds", plot.GrowthCompleteTimestamp - time.Now().Unix())
	// 	responses.SendRes(w, responses.Plants_Still_Growing, plot, timestampMsg)
	// 	return
	// }

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
	// If consumables included, validate them
	if body.Consumable != string("") {
		// Validate specified consumable is a good
		goodsDict := (*h.MainDictionary).Goods
		if _, ok := goodsDict[body.Consumable]; !ok {
			// Fail, seed is not good
			errmsg := fmt.Sprintf("Error in InteractPlot, consumable item does not exist in good dictionary. received consumable name: %v", body.Consumable)
			log.Error.Printf(errmsg)
			responses.SendRes(w, responses.Item_Does_Not_Exist, nil, errmsg)
			return
		}
		// Validate consumable specified in given warehouse
		temp, ownedConsumableOk := warehouse.Goods[body.Consumable]; 
		if !ownedConsumableOk {
			// Fail, consumable good not in warehouse
			errmsg := fmt.Sprintf("Error in InteractPlot, consumable item not found in local warehouse. received good name: %v, warehouse goods: %v", body.Consumable, warehouse.Goods)
			log.Error.Printf(errmsg)
			responses.SendRes(w, responses.Not_Enough_Items_In_Warehouse, nil, errmsg)
			return
		}
		consumableQuantityAvailable = temp
	}

	// Validate plot available for interaction and body meets internal plot validation
	plantDef, plantDefOk := h.MainDictionary.Plants[plot.PlantedPlant.PlantType.String()]

	if !plantDefOk {
		plotDefErrMsg := fmt.Sprintf("Error in InteractPlot, Plot [%s] has PlantedPlant type [%s] not found in main dictionary!", plot.UUID, plot.PlantedPlant.PlantType.String())
		log.Error.Println(plotDefErrMsg)
		responses.SendRes(w, responses.Internal_Server_Error, plot, plotDefErrMsg)
	}
	plotValidationResponse, addedYield, usedConsumableQuantity, growthHarvest, growthTime := plot.IsInteractable(body, plantDef, consumableQuantityAvailable, warehouse.Tools)
	switch plotValidationResponse {
	case responses.Invalid_Plot_Action:
		log.Error.Printf("Error in PlantPlot, Invalid_Plot_Action")
		responses.SendRes(w, responses.Invalid_Plot_Action, plot, "")
		return
	case responses.Tool_Not_Found:
		log.Error.Printf("Error in PlantPlot, Tool_Not_Found")
		responses.SendRes(w, responses.Tool_Not_Found, plot, "")
		return
	case responses.Missing_Consumable_Selection:
		log.Error.Printf("Error in PlantPlot, Missing_Consumable_Selection")
		responses.SendRes(w, responses.Missing_Consumable_Selection, plot, "")
		return
	case responses.Internal_Server_Error:
		log.Error.Printf("Error in PlantPlot, Internal_Server_Error")
		responses.SendRes(w, responses.Internal_Server_Error, plot, "Could not get scaled growth stage, contact Developer")
		return
	case responses.Not_Enough_Items_In_Warehouse:
		log.Error.Printf("Error in PlantPlot, Not_Enough_Items_In_Warehouse")
		responses.SendRes(w, responses.Not_Enough_Items_In_Warehouse, plot, "")
		return
	case responses.Consumable_Not_In_Options:
		log.Error.Printf("Error in PlantPlot, Consumable_Not_In_Options")
		responses.SendRes(w, responses.Consumable_Not_In_Options, plot, "")
		return
	case responses.Generic_Success:
		log.Debug.Printf("Plot growth action validated successfully: %s, action: %s", plot.UUID, body.Action)
	default:
		log.Error.Fatalf("Received unexpected response type from plot.IsPlantable. plot: %v body: %v", plot, body)
	}

	plot.PlantedPlant.CurrentStage++
	plot.PlantedPlant.Yield += addedYield
	plot.GrowthCompleteTimestamp = time.Now().Unix() + growthTime
	farm.Plots[uuid] = plot

	// Save to DBs
	
	saveFarmErr := schema.SaveFarmDataAtPathToDB(fdb, farmLocationSymbol, "plots", farm.Plots)
	if saveFarmErr != nil {
		log.Error.Printf("Error in PlotInteract, could not save farm. error: %v", saveFarmErr)
		responses.SendRes(w, responses.DB_Save_Failure, nil, saveFarmErr.Error())
		return
	}

	if body.Consumable != string("") {
		// if consumables used
		warehouse.RemoveGoods(body.Consumable, usedConsumableQuantity)
	}

	var nextStage *schema.GrowthStage

	if growthHarvest != nil {
		// if was a harvest action
		// TODO: write harvest logic
		log.Important.Println("harvest!")
		// if not final harvest
		if !growthHarvest.FinalHarvest {
			nextStage = &h.MainDictionary.Plants[plot.PlantedPlant.PlantType.String()].GrowthStages[plot.PlantedPlant.CurrentStage]
		}
	} else {
		nextStage = &h.MainDictionary.Plants[plot.PlantedPlant.PlantType.String()].GrowthStages[plot.PlantedPlant.CurrentStage]
	}
	saveWarehouseErr := schema.SaveWarehouseDataAtPathToDB(wdb, warehouseLocationSymbol, "goods", warehouse.Goods)
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