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
	"math"
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
		responses.SendRes(w, responses.Generic_Success, caravans, "None Found")
		return
	}
	// Modify Caravan SecondsTillArrival
	for i, caravan := range caravans {
		caravan.SecondsTillArrival = caravan.ArrivalTime - time.Now().Unix() 
		if caravan.SecondsTillArrival < 0 {
			caravan.SecondsTillArrival = 0
		}
		caravans[i] = caravan
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
	// Modify Caravan SecondsTillArrival
	caravan.SecondsTillArrival = caravan.ArrivalTime - time.Now().Unix() 
	if caravan.SecondsTillArrival < 0 {
		caravan.SecondsTillArrival = 0
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

// Handler function for the secure route: PATCH: /api/my/caravans/
type CharterCaravan struct {
	Dbs *map[string]rdb.Database
	World *schema.World
}
func (h *CharterCaravan) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- CharterCaravan --"))
	// Get user info
	udb := (*h.Dbs)["users"]
	OK, userData, _ := secureGetUser(w, r, udb)
	if !OK {
		return // Failure states handled by secureGetUser, simply return
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

	// Validate Caravan Charter Request Body
	validationMap := schema.ValidateCaravanCharter(body)
	if len(validationMap) > 0 {
		// Failed basic validation
		errmsg := fmt.Sprintf("Validation Error in CharterCaravan: %v", validationMap)
		log.Debug.Printf(errmsg)
		responses.SendRes(w, responses.Bad_Request, validationMap, "Request body did not pass validation, see data for specifics.")
		return
	}

	// Get Assistants
	adb := (*h.Dbs)["assistants"]
	assistantLocationSymbols := make([]string, len(body.Assistants))
	for i, assistantID := range body.Assistants {
		assistantLocationSymbols[i] = userData.Username + "|Assistant-" + assistantID
	}
	assistants, foundAssistants, assistantsErr := schema.GetAssistantsFromDB(assistantLocationSymbols, adb)
	if assistantsErr != nil || !foundAssistants {
		errmsg := fmt.Sprintf("Error in CharterCaravan, could not get assistants from DB. foundWarehouse: %v, error: %v", foundAssistants, assistantsErr)
		log.Error.Printf(errmsg)
		responses.SendRes(w, responses.DB_Get_Failure, nil, errmsg)
		return
	}

	// Generate a timestamp for caravan id. Get slowest assistant speeds, total carry cap, set their location to caravan UUID
	caravanTimestamp := time.Now()
	caravanUUID := userData.Username + "|Caravan-" + fmt.Sprintf("%d", caravanTimestamp.UnixNano())
	slowestSpeed := int(1000000)
	carryCap := int(0)
	assistantOriginValidation := make(map[string]string)
	for i, assistant := range assistants {
		if assistant.Location != body.Origin {
			// Fail, assistant not at origin, prepare validation response will error later
			assistantOriginValidation[fmt.Sprintf("Assistant-%s", i)] = fmt.Sprintf("Assistant not at origin (%s) specified in request", body.Origin)
		}
		if assistant.Speed < slowestSpeed {
			slowestSpeed = assistant.Speed
			log.Debug.Printf("Setting slowest speed to %d based on %d: %s", assistant.Speed, assistant.ID, assistant.Archetype.String())
		}
		carryCap += assistant.CarryCap
		assistant.Location = caravanUUID
		assistants[i] = assistant
	}

	if len(assistantOriginValidation) > 0 {
		// Fail, found assitant not at origin
		errmsg := fmt.Sprintf("Validation Error in CharterCaravan: %v", assistantOriginValidation)
		log.Debug.Printf(errmsg)
		responses.SendRes(w, responses.Bad_Request, assistantOriginValidation, "Request body did not pass validation, see data for specifics.")
		return
	}

	// Calculate travel time and construct caravan
	travelTimeValidationMap, caravanTravelTime := schema.CalculateTravelTime((*h.World), body.Origin, body.Destination, slowestSpeed)
	if len(travelTimeValidationMap) > 0 {
		// Fail, origin and destination could not be routed
		errmsg := fmt.Sprintf("Validation Error in CharterCaravan: %v", assistantOriginValidation)
		log.Debug.Printf(errmsg)
		responses.SendRes(w, responses.Bad_Request, assistantOriginValidation, "Request body did not pass validation, see data for specifics.")
		return
	}
	caravan := schema.NewCaravan(caravanUUID, caravanTimestamp, body.Origin, body.Destination, body.Assistants, body.Wares, caravanTravelTime)
	log.Debug.Printf("Prepared caravan, now to validate wares. Caravan: %v", caravan)

	// Validate wares if present
	wareCategories := make(map[string]int)
	if len(body.Wares.Goods) > 0 {
		wareCategories["Goods"] = len(body.Wares.Goods)
	}
	if len(body.Wares.Produce) > 0 {
		wareCategories["Produce"] = len(body.Wares.Produce)
	}
	if len(body.Wares.Seeds) > 0 {
		wareCategories["Seeds"] = len(body.Wares.Seeds)
	}
	if len(body.Wares.Tools) > 0 {
		wareCategories["Tools"] = len(body.Wares.Tools)
	}
	if len(wareCategories) >= 1 {
		// Wares found
		log.Debug.Printf("Wares found: %v", wareCategories)

		// Get carry cap
		log.Debug.Printf("carryCap before team_factor: %d", carryCap)
		teamFactor := 1 + (float64(0.1) * float64(len(assistants) - 1))
		carryCap := uint64(math.Ceil(float64(carryCap) * teamFactor))
		log.Debug.Printf("carryCap after team_factor (%f): %d", teamFactor, carryCap)

		// Get warehouse
		wdb := (*h.Dbs)["warehouses"]
		warehouseLocationSymbol := userData.Username + "|Warehouse-" + body.Origin
		warehouse, foundWarehouse, warehousesErr := schema.GetWarehouseFromDB(warehouseLocationSymbol, wdb)
		if warehousesErr != nil || !foundWarehouse {
			errmsg := fmt.Sprintf("Error in CharterCaravan, could not get warehouse from DB. foundWarehouse: %v, error: %v", foundWarehouse, warehousesErr)
			log.Error.Printf(errmsg)
			responses.SendRes(w, responses.DB_Get_Failure, nil, errmsg)
			return
		}

		// Validate warehouse contains wares in specified quantities and that there is enough carry capacity for specified goods (after multiplying produce quantity by size)
		waresValidationMap := make(map[string]map[string]string)
		carryCapNeeded := uint64(0)
		for category := range wareCategories {
			switch category {
			case "Goods":
				for item, quantity := range body.Wares.Goods {
					if _, ok := warehouse.Goods[item]; !ok {
						// FAIL item not in warehouse
						if len(waresValidationMap["goods"]) < 1 {
							waresValidationMap["goods"] = make(map[string]string)
						}
						waresValidationMap["goods"][item] = fmt.Sprintf("Item not found in local warehouse (received name: %s)", item)
						continue
					}
					if warehouse.Goods[item] < quantity {
						// FAIL not enough item
						if len(waresValidationMap["goods"]) < 1 {
							waresValidationMap["goods"] = make(map[string]string)
						}
						waresValidationMap["goods"][item] = fmt.Sprintf("Not enough in local warehouse (requested: %d, have: %d)", quantity, warehouse.Goods[item])
						continue
					}
					log.Debug.Printf("Passed Validation, Remove Goods: %s %d", item, quantity)
					carryCapNeeded += quantity
					warehouse.RemoveGoods(item, quantity)
				}
			case "Produce":
				for item, quantity := range body.Wares.Produce {
					_, psize, splittable := warehouse.GetProduceNameSizeSlice(item)
					if !splittable {
						// FAIL produce must have size component
						if len(waresValidationMap["produce"]) < 1 {
							waresValidationMap["produce"] = make(map[string]string)
						}
						waresValidationMap["produce"][item] = fmt.Sprintf("Could not split to find produce size, ensure follows example format 'Potato|Tiny' (received name: %s)", item)
						continue
					}
					if _, ok := warehouse.Produce[item]; !ok {
						// FAIL item not in warehouse
						if len(waresValidationMap["produce"]) < 1 {
							waresValidationMap["produce"] = make(map[string]string)
						}
						waresValidationMap["produce"][item] = fmt.Sprintf("Item not found in local warehouse (received name: %s)", item)
						continue
					}
					if warehouse.Produce[item] < quantity {
						// FAIL not enough item
						if len(waresValidationMap["produce"]) < 1 {
							waresValidationMap["produce"] = make(map[string]string)
						}
						waresValidationMap["produce"][item] = fmt.Sprintf("Not enough in local warehouse (requested: %d, have: %d)", quantity, warehouse.Produce[item])
						continue
					}
					log.Debug.Printf("Passed Validation, Remove Produce: %s %d", item, quantity)
					
					carryCapNeeded += quantity * uint64(schema.SizeToID[psize])
					warehouse.RemoveProduce(item, quantity)
				}
			case "Seeds":
				for item, quantity := range body.Wares.Seeds {
					if _, ok := warehouse.Seeds[item]; !ok {
						// FAIL item not in warehouse
						if len(waresValidationMap["seeds"]) < 1 {
							waresValidationMap["seeds"] = make(map[string]string)
						}
						waresValidationMap["seeds"][item] = fmt.Sprintf("Item not found in local warehouse (received name: %s)", item)
						continue
					}
					if warehouse.Seeds[item] < quantity {
						// FAIL not enough item
						if len(waresValidationMap["seeds"]) < 1 {
							waresValidationMap["seeds"] = make(map[string]string)
						}
						waresValidationMap["seeds"][item] = fmt.Sprintf("Not enough in local warehouse (requested: %d, have: %d)", quantity, warehouse.Seeds[item])
						continue
					}
					log.Debug.Printf("Passed Validation, Remove Seeds: %s %d", item, quantity)
					carryCapNeeded += quantity
					warehouse.RemoveSeeds(item, quantity)
				}
			case "Tools":
				for item, quantity := range body.Wares.Tools {
					if _, ok := warehouse.Tools[item]; !ok {
						// FAIL item not in warehouse
						if len(waresValidationMap["tools"]) < 1 {
							waresValidationMap["tools"] = make(map[string]string)
						}
						waresValidationMap["tools"][item] = fmt.Sprintf("Item not found in local warehouse (received name: %s)", item)
						continue
					}
					if warehouse.Tools[item] < quantity {
						// FAIL not enough item
						if len(waresValidationMap["tools"]) < 1 {
							waresValidationMap["tools"] = make(map[string]string)
						}
						waresValidationMap["tools"][item] = fmt.Sprintf("Not enough in local warehouse (requested: %d, have: %d)", quantity, warehouse.Tools[item])
						continue
					}
					log.Debug.Printf("Passed Validation, Remove Tools: %s %d", item, quantity)
					carryCapNeeded += quantity
					warehouse.RemoveTools(item, quantity)
				}
			default:
				errmsg := fmt.Sprintf("Error in CharterCaravan, unexpected ware category after validating wares: %s", category)
				log.Error.Printf(errmsg)
				responses.SendRes(w, responses.Internal_Server_Error, nil, errmsg)
				return
			}
		}
		if carryCapNeeded > carryCap {
			// Too many wares or not enough assistants
			if len(waresValidationMap["meta"]) < 1 {
				waresValidationMap["meta"] = make(map[string]string)
			}
			waresValidationMap["meta"]["carrying_capacity"] = fmt.Sprintf("Not enough carrying capacity, have: %d, need: %d. Add more assistants or split into smaller loads", carryCap, carryCapNeeded)
		}
		if len(waresValidationMap) > 0 {
			// Wares failed validation
			resWareVMap := map[string]map[string]map[string]string{"wares": waresValidationMap}
			errmsg := fmt.Sprintf("Validation Error in CharterCaravan: %v", resWareVMap)
			log.Debug.Printf(errmsg)
			responses.SendRes(w, responses.Bad_Request, resWareVMap, "Request body did not pass validation, see data for specifics.")
			return
		}

		// If found all, remove from local warehouse & save (dont need to add anywhere cause charter already specifies wares)
		saveWarehouseErr := schema.SaveWarehouseToDB(wdb, &warehouse)
		if saveWarehouseErr != nil {
			log.Error.Printf("Error in CharterCaravan, could not save warehouse. error: %v", saveWarehouseErr)
			responses.SendRes(w, responses.DB_Save_Failure, nil, saveWarehouseErr.Error())
			return
		}
	}

	// Caravan valid and prepared, save assistants, user, and save caravan
	cdb := (*h.Dbs)["caravans"]
	saveCaravanErr := schema.SaveCaravanToDB(cdb, caravan)
	if saveCaravanErr != nil {
		log.Error.Printf("Error in PlotInteract, could not save caravan. error: %v", saveCaravanErr)
		responses.SendRes(w, responses.DB_Save_Failure, nil, saveCaravanErr.Error())
		return
	}

	caravanList := append(userData.Caravans, caravanUUID)
	saveUserErr := schema.SaveUserDataAtPathToDB(udb, userData.Token, "caravans", caravanList)
	if saveUserErr != nil {
		log.Error.Printf("Error in PlotInteract, could not save user. error: %v", saveUserErr)
		responses.SendRes(w, responses.DB_Save_Failure, nil, saveUserErr.Error())
		return
	}

	for _, assistant := range assistants {
		saveAssistantErr := schema.SaveAssistantDataAtPathToDB(adb, assistant.UUID, "location", assistant.Location)
		if saveAssistantErr != nil {
				log.Error.Printf("Error in CharterCaravan, could not save assistant. error: %v", saveAssistantErr)
				responses.SendRes(w, responses.DB_Save_Failure, nil, saveAssistantErr.Error())
				return
		}
	}

	// Construct and Send response
	getCaravanJsonString, getCaravanJsonStringErr := responses.JSON(caravan)
	if getCaravanJsonStringErr != nil {
		log.Error.Printf("Error in CharterCaravan, could not format caravan response as JSON. response: %v, error: %v", caravan, getCaravanJsonStringErr)
		responses.SendRes(w, responses.JSON_Marshal_Error, caravan, getCaravanJsonStringErr.Error())
		return
	}
	log.Debug.Printf("Sending response for CharterCaravan:\n%v", getCaravanJsonString)
	responses.SendRes(w, responses.Generic_Success, caravan, "")
	
	log.Debug.Println(log.Cyan("-- End CharterCaravan --"))
}

// Handler function for the secure route: DELETE: /api/my/caravans/{caravan-id}
type UnpackCaravan struct {
	Dbs *map[string]rdb.Database
}
func (h *UnpackCaravan) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- UnpackCaravan --"))
	// Get symbol from route
	id := GetVarEntries(r, "caravan-id", None)

	// Get user info
	udb := (*h.Dbs)["users"]
	OK, userData, _ := secureGetUser(w, r, udb)
	if !OK {
		return // Failure states handled by secureGetUser, simply return
	}

	// Parse ID
	cIdSlice := strings.Split(id, "-")
	var cId string
	if len(cIdSlice) > 1 {
		cId = cIdSlice[:len(cIdSlice)-1][0]
	} else {
		cId = cIdSlice[0]
	}
	
	cUUID := userData.Username + "|Caravan-" + cId
	log.Debug.Printf("UnpackCaravan UUID: %s", cUUID)

	// Get Caravan
	cdb := (*h.Dbs)["caravans"]
	caravan, foundCaravan, caravanErr := schema.GetCaravanFromDB(cUUID, cdb)
	if caravanErr != nil {
		errmsg := fmt.Sprintf("Error in UnpackCaravan, could not get caravan from DB. foundCaravan: %v, error: %v", foundCaravan, caravanErr)
		log.Error.Printf(errmsg)
		responses.SendRes(w, responses.DB_Get_Failure, nil, errmsg)
		return
	}
	if !foundCaravan {
		// Validation error
		errmsg := fmt.Sprintf("caravan of given id not found")
		log.Debug.Printf(errmsg)
		responses.SendRes(w, responses.Bad_Request, nil, errmsg)
		return
	}

	// Validate Timestamp
	now := time.Now()
	if caravan.ArrivalTime > now.Unix() {
		// Too early
		caravan.SecondsTillArrival = caravan.ArrivalTime - now.Unix()
		errmsg := fmt.Sprintf("caravan has not arrived yet, arrives in %d seconds", caravan.SecondsTillArrival)
		log.Debug.Printf(errmsg)
		responses.SendRes(w, responses.Caravan_Not_Arrived, caravan, errmsg)
		return
	}

	// VALID: Update user, warehouses, assistants, caravans DBs

	// Get warehouse
	wdb := (*h.Dbs)["warehouses"]
	warehouseLocationSymbol := userData.Username + "|Warehouse-" + caravan.Destination
	warehouse, foundWarehouse, warehousesErr := schema.GetWarehouseFromDB(warehouseLocationSymbol, wdb)
	if warehousesErr != nil {
		errmsg := fmt.Sprintf("Error in UnpackCaravan, could not get warehouse from DB. foundWarehouse: %v, error: %v", foundWarehouse, warehousesErr)
		log.Error.Printf(errmsg)
		responses.SendRes(w, responses.DB_Get_Failure, nil, errmsg)
		return
	}
	if !foundWarehouse {
		// Need to add new warehouse
		warehouse = *schema.NewEmptyWarehouse(userData.Username, caravan.Destination)
		userData.Warehouses = append(userData.Warehouses, warehouse.UUID)
	}

	// Update User
	if len(userData.Caravans) <= 1 {
		userData.Caravans = make([]string, 0)
	} else {
		userData.Caravans = remove(userData.Caravans, cUUID)
	}
	// Update Warehouse
	if len(caravan.Wares.Goods) > 0 {
		for g, q := range caravan.Wares.Goods {
			warehouse.AddGoods(g, q)
		}
	}
	if len(caravan.Wares.Seeds) > 0 {
		for s, q := range caravan.Wares.Seeds {
			warehouse.AddSeeds(s, q)
		}
	}
	if len(caravan.Wares.Tools) > 0 {
		for t, q := range caravan.Wares.Tools {
			warehouse.AddTools(t, q)
		}
	}
	if len(caravan.Wares.Produce) > 0 {
		for p, q := range caravan.Wares.Produce {
			warehouse.AddProduce(p, q)
		}
	}
	// Update and Save Assistants
	adb := (*h.Dbs)["assistants"]
	for _, aID := range caravan.Assistants {
		aUUID := userData.Username + "|Assistant-" + aID
		saveAssistantErr := schema.SaveAssistantDataAtPathToDB(adb, aUUID, "location", caravan.Destination)
		if saveAssistantErr != nil {
			log.Error.Printf("Error in UnpackCaravan, could not save assistant. error: %v", saveAssistantErr)
			responses.SendRes(w, responses.DB_Save_Failure, nil, saveAssistantErr.Error())
			return
		}
	}

	// Save User
	saveUserErr := schema.SaveUserToDB(udb, &userData)
	if saveUserErr != nil {
		log.Error.Printf("Error in UnpackCaravan, could not save user. error: %v", saveUserErr)
		responses.SendRes(w, responses.DB_Save_Failure, nil, saveUserErr.Error())
		return
	}
	// Save Warehouse
	saveWarehouseErr := schema.SaveWarehouseToDB(wdb, &warehouse)
	if saveWarehouseErr != nil {
		log.Error.Printf("Error in UnpackCaravan, could not save warehouse. error: %v", saveWarehouseErr)
		responses.SendRes(w, responses.DB_Save_Failure, nil, saveWarehouseErr.Error())
		return
	}
	// Delete Caravan
	delCaravanErr := schema.DeleteCaravanFromDB(cdb, cUUID)
	if delCaravanErr != nil {
		log.Error.Printf("Error in UnpackCaravan, could not delete caravan. error: %v", delCaravanErr)
		responses.SendRes(w, responses.Internal_Server_Error, nil, delCaravanErr.Error())
		return
	}

	// Construct and Send response
	response := map[string]interface{}{"local_warehouse": &warehouse, "assistants_released": &caravan.Assistants}
	getPlotPlantResponseJsonString, getPlotPlantResponseJsonStringErr := responses.JSON(response)
	if getPlotPlantResponseJsonStringErr != nil {
		log.Error.Printf("Error in PlotInfo, could not format interact response as JSON. response: %v, error: %v", response, getPlotPlantResponseJsonStringErr)
		responses.SendRes(w, responses.JSON_Marshal_Error, response, getPlotPlantResponseJsonStringErr.Error())
		return
	}
	log.Debug.Printf("Sending response for UnpackCaravan:\n%v", getPlotPlantResponseJsonString)
	responses.SendRes(w, responses.Generic_Success, response, fmt.Sprintf("Caravan unpacked at %s", caravan.Destination))
	
	log.Debug.Println(log.Cyan("-- End UnpackCaravan --"))
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

// Handler function for the secure route: POST: /api/my/farms/{location-symbol}/ritual
type ConductRitual struct {
	Dbs *map[string]rdb.Database
	MainDictionary *schema.MainDictionary
}
func (h *ConductRitual) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- ConductRitual --"))
	udb := (*h.Dbs)["users"]
	OK, userData, _ := secureGetUser(w, r, udb)
	if !OK {
		return // Failure states handled by secureGetUser, simply return
	}
	// unmarshall request body to get ritual
	var body schema.CaravanCharter
	decoder := json.NewDecoder(r.Body)
	if decodeErr := decoder.Decode(&body); decodeErr != nil {
		// Fail case, could not decode
		errmsg := fmt.Sprintf("Decode Error in CharterCaravan: %v", decodeErr)
		log.Debug.Printf(errmsg)
		responses.SendRes(w, responses.Bad_Request, nil, "Could not decode request body, ensure it conforms to expected format.")
		return
	}
	adb := (*h.Dbs)["assistants"]
	assistants, foundAssistants, assistantsErr := schema.GetAssistantsFromDB(userData.Assistants, adb)
	if assistantsErr != nil || !foundAssistants {
		log.Error.Printf("Error in ConductRitual, could not get assistants from DB. foundAssistants: %v, error: %v", foundAssistants, assistantsErr)
		responses.SendRes(w, responses.DB_Get_Failure, assistants, assistantsErr.Error())
		return
	}
	getAssistantJsonString, getAssistantJsonStringErr := responses.JSON(assistants)
	if getAssistantJsonStringErr != nil {
		log.Error.Printf("Error in ConductRitual, could not format assistants as JSON. assistants: %v, error: %v", assistants, getAssistantJsonStringErr)
		responses.SendRes(w, responses.JSON_Marshal_Error, assistants, getAssistantJsonStringErr.Error())
		return
	}
	log.Debug.Printf("Sending response for ConductRitual:\n%v", getAssistantJsonString)
	responses.SendRes(w, responses.Generic_Success, assistants, "")
	log.Debug.Println(log.Cyan("-- End ConductRitual --"))
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
	// Validate specified seed meets min/max size requirements
	plantDict := (*h.MainDictionary).Plants
	plantDef, plantDefOk := plantDict[plantName]
	if !plantDefOk {
		// Fail, could not lookup plant with matching seed name, internal error
		errmsg := fmt.Sprintf("Error in PlantPlot, found seed but not plant in master dictionary... received seed name: %v", body.SeedName)
		log.Error.Printf(errmsg)
		responses.SendRes(w, responses.Internal_Server_Error, nil, errmsg)
		return
	}
	if body.SeedSize.String() == "" {
		// FAIL invalid size
		errmsg := fmt.Sprintf("in PlantPlot, invalid size")
		log.Debug.Printf(errmsg)
		responses.SendRes(w, responses.Bad_Request, nil, errmsg)
		return
	}
	if uint16(body.SeedSize) < uint16(schema.SizeToID[plantDef.MinSize]) {
		// FAIL, this plant won't grow that small
		errmsg := fmt.Sprintf("in PlantPlot, this plant won't grow that small (%s), minsize: %v", body.SeedSize, plantDef.MinSize)
		log.Debug.Printf(errmsg)
		responses.SendRes(w, responses.Bad_Request, nil, errmsg)
		return
	}
	if uint16(body.SeedSize) > uint16(schema.SizeToID[plantDef.MaxSize]) {
		// FAIL, this plant won't grow that large
		errmsg := fmt.Sprintf("in PlantPlot, this plant won't grow that large (%s), maxsize: %v", body.SeedSize, plantDef.MaxSize)
		log.Debug.Printf(errmsg)
		responses.SendRes(w, responses.Bad_Request, nil, errmsg)
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
		log.Debug.Printf("in PlantPlot, plot already planted")
		responses.SendRes(w, responses.Plot_Already_Planted, plot, "")
		return
	case responses.Bad_Request:
		log.Debug.Printf("in PlantPlot, seed size invalid")
		responses.SendRes(w, responses.Bad_Request, plot, "Seed Size specified in request body was invalid")
		return
	case responses.Plot_Too_Small:
		log.Debug.Printf("in PlantPlot, plot too small")
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
	plotValidationResponse, addedYield, usedConsumableQuantity, growthHarvest, growthTime, repeatStage, errInfoMsg := plot.IsInteractable(body, plantDef, consumableQuantityAvailable, warehouse.Tools)
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
		for producename, producequantity := range harvest.Produce {
			log.Debug.Printf("Add produce quantity: %d", producequantity)
			warehouse.AddProduce(producename, producequantity)
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
			// move up current stage if not repeat, initialize nextStage for response
			if !repeatStage {
				plot.PlantedPlant.CurrentStage++
			}
			nextStage = &h.MainDictionary.Plants[plot.PlantedPlant.PlantType].GrowthStages[plot.PlantedPlant.CurrentStage]
		}
	} else {
		// if not harvest
		// move up current stage if not repeat, initialize nextStage for response
		if !repeatStage {
			plot.PlantedPlant.CurrentStage++
		}
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
		warehouseDict = warehouse.Produce
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
		warehouse.Produce = warehouseDict
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