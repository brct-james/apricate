// Package handlers provides functions for handling web routes
package handlers

import (
	"apricate/auth"
	"apricate/log"
	"apricate/metrics"
	"apricate/rdb"
	"apricate/responses"
	"apricate/schema"
	"errors"
	"fmt"
	"net/http"
	"strings"

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
	// TODO: IMPLEMENT FARMS, for now only worry about starting farm
	farms := []schema.Farm{{
		UUID: "",
		RegionLocation: "Pria|Homestead Farm",
		Bonuses: []schema.FarmBonuses{},
		Tools: make(map[string]uint8),
		Buildings: make(map[string]uint8),
		Plots: []string{""},
	}}
	// use myRLs as a set to get all unique locations visible in fow
	myRLs := make(map[string]bool)
	for _, assistant := range assistants {
		myRLs[assistant.RegionLocation] = true
	}
	for _, farm := range farms {
		myRLs[farm.RegionLocation] = true
	}
	// finally get all locations in each region
	resLocations := make([]schema.Location, 0)
	for rl := range myRLs {
		split := strings.Split(rl, "|")
		region, location := split[0], split[1]
		resLocations = append(resLocations, h.World.Locations[region][location])
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
	// TODO: IMPLEMENT FARMS, for now only worry about starting farm
	farms := []schema.Farm{{
		UUID: "",
		RegionLocation: "Homestead Farm",
		Bonuses: []schema.FarmBonuses{},
		Tools: make(map[string]uint8),
		Buildings: make(map[string]uint8),
		Plots: []string{""},
	}}
	// use myRLs as a set
	myRLs := make(map[string]bool)
	for _, assistant := range assistants {
		myRLs[assistant.RegionLocation] = true
	}
	for _, farm := range farms {
		myRLs[farm.RegionLocation] = true
	}
	// get every region revealed based on location using myRegions as a set
	myRegions := make(map[string]bool)
	for regionloc := range myRLs {
		split := strings.Split(regionloc, "|")
		region := split[0]
		myRegions[region] = true
	}
	// finally get all locations in each region
	resLocations := make(map[string][]string, 0)
	i := 0
	for regionName := range myRegions {
		for _, loc := range h.World.Locations[regionName] {
			resLocations[regionName] = append(resLocations[regionName], loc.Name)
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
	// TODO: IMPLEMENT FARMS, for now only worry about starting farm
	farms := []schema.Farm{{
		UUID: "",
		RegionLocation: "Pria|Homestead Farm",
		Bonuses: []schema.FarmBonuses{},
		Tools: make(map[string]uint8),
		Buildings: make(map[string]uint8),
		Plots: []string{""},
	}}
	// use myRLs as a set to get all unique locations visible in fow
	myRLs := make(map[string]bool)
	for _, assistant := range assistants {
		myRLs[assistant.RegionLocation] = true
	}
	for _, farm := range farms {
		myRLs[farm.RegionLocation] = true
	}
	// Get name from route
	route_vars := mux.Vars(r)
	name := strings.ToUpper(strings.ReplaceAll(route_vars["name"], "_", " "))
	// finally get specified location if available
	var resLocation schema.Location
	found := false
	for rl := range myRLs {
		split := strings.Split(rl, "|")
		region, location := split[0], split[1]
		if strings.ToUpper(location) == name {
			resLocation = h.World.Locations[region][location]
			found = true
		}
	}
	if !found {
		responses.SendRes(w, responses.No_Assitant_At_Location, nil, "")
		return
	}
	responses.SendRes(w, responses.Generic_Success, resLocation, "")
	log.Debug.Println(log.Cyan("-- End LocationInfo --"))
}