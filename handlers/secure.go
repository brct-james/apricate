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