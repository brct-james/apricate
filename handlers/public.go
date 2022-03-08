// Package handlers provides functions for handling web routes
package handlers

import (
	"fmt"
	"net/http"

	"apricate/auth"
	"apricate/log"
	"apricate/metrics"
	"apricate/rdb"
	"apricate/responses"
	"apricate/schema"
	"apricate/tokengen"

	"github.com/gorilla/mux"
)

// Helper Functions

// // Attempt to get user from db
// func publicGetUser(w http.ResponseWriter, r *http.Request, username string, token string) (bool, schema.User, rdb.Database) {
// 	// Get udb from context
// 	udb, udbErr := GetUdbFromCtx(r)
// 	if udbErr != nil {
// 		// Fail state getting context
// 		log.Error.Printf("Could not get UserDBContext in publicGetUser")
// 		responses.SendRes(w, responses.No_UDB_Context, nil, "in publicGetUser")
// 		return false, schema.User{}, rdb.Database{}
// 	}
// 	// Check db for user
// 	thisUser, userFound, getUserErr := schema.GetUserFromDB(token, udb)
// 	if getUserErr != nil {
// 		// fail state
// 		getErrorMsg := fmt.Sprintf("in publicGetUser, could not get from DB for username: %s, error: %v", username, getUserErr)
// 		responses.SendRes(w, responses.UDB_Get_Failure, nil, getErrorMsg)
// 		return false, schema.User{}, rdb.Database{}
// 	}
// 	if !userFound {
// 		// fail state - user not found
// 		userNotFoundMsg := fmt.Sprintf("in publicGetUser, no user found in DB with username: %s", username)
// 		responses.SendRes(w, responses.User_Not_Found, nil, userNotFoundMsg)
// 		return false, schema.User{}, rdb.Database{}
// 	}
// 	// Success case
// 	return true, thisUser, udb
// }

// Handler Functions

// Handler function for the route: /
func Homepage(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- Homepage --"))
	responses.SendRes(w, responses.Unimplemented, nil, "Homepage")
	log.Debug.Println(log.Cyan("-- End Homepage --"))
}

// Handler function for the route: /api/users/{username}/claim
type UsernameClaim struct {
	Dbs *map[string]rdb.Database
}
func (h *UsernameClaim) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- usernameClaim --"))
	log.Debug.Println("Recover udb from context")
	udb := (*h.Dbs)["users"]
	// Get username from route
	route_vars := mux.Vars(r)
	username := route_vars["username"]
	log.Debug.Printf("Username Requested: %s", username)
	// Validate username (length & content, plus characters)
	usernameValidationStatus := auth.ValidateUsername(username)
	if usernameValidationStatus != "OK" {
		// fail state
		validationErrorMessage := fmt.Sprintf("in UsernameClaim: Username: %v | ValidationResponse: %v", username, usernameValidationStatus)
		log.Debug.Println(validationErrorMessage)
		responses.SendRes(w, responses.Username_Validation_Failure, nil, validationErrorMessage)
		return
	}
	// generate token
	token, genTokenErr := tokengen.GenerateToken(username)
	if genTokenErr != nil {
		// fail state
		log.Important.Printf("in UsernameClaim: Attempted to generate token using username %s but was unsuccessful with error: %v", username, genTokenErr)
		genErrorMsg := fmt.Sprintf("Username: %v | GenerateTokenErr: %v", username, genTokenErr)
		responses.SendRes(w, responses.Generate_Token_Failure, nil, genErrorMsg)
		return
	}
	// check DB for existing user
	userExists, dbGetError := schema.CheckForExistingUser(token, udb)
	if dbGetError != nil {
		// fail state - db error
		dbGetErrorMsg := fmt.Sprintf("in UsernameClaim | Username: %v | UDB Get Error: %v", username, dbGetError)
		log.Debug.Println(dbGetErrorMsg)
		responses.SendRes(w, responses.UDB_Get_Failure, nil, dbGetErrorMsg)
		return
	}
	if userExists {
		// fail state - user already exists
		validationFailMsg := fmt.Sprintf("in UsernameClaim | Username: %v | Reason: USER_ALREADY_EXISTS", username)
		log.Debug.Println(validationFailMsg)
		responses.SendRes(w, responses.Username_Validation_Failure, nil, validationFailMsg)
		return
	}
	// create new user in DB
	newUser := schema.NewUser(token, username)
	saveUserErr := schema.SaveUserToDB(udb, newUser)
	if saveUserErr != nil {
		// fail state - could not save
		saveUserErrMsg := fmt.Sprintf("in UsernameClaim | Username: %v | CreateNewUserInDB failed, dbSaveResult: %v", username, saveUserErr)
		log.Debug.Println(saveUserErrMsg)
		responses.SendRes(w, responses.DB_Save_Failure, nil, saveUserErrMsg)
		return
	}
	// Created successfully
	// Track in user metrics
	metrics.TrackNewUser(username)
	log.Debug.Printf("Generated token %s and claimed username %s", token, username)
	responses.SendRes(w, responses.Generic_Success, newUser, "")
	log.Debug.Println(log.Cyan("-- End usernameClaim --"))
}