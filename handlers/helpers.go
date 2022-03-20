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

type StringFormat uint16
const (
	None StringFormat = 0
	UnderscoresToSpaces StringFormat = 1
	TrueTitle StringFormat = 2
	AllCaps StringFormat = 3
	UUID StringFormat = 4
	SpacedName StringFormat = 5
)

// Get entries from mux vars in correct format
func GetVarEntries(r *http.Request, varName string, format StringFormat) string {
	route_vars := mux.Vars(r)
	entry := route_vars[varName]
	switch format {
	case None:
		// none
	case UnderscoresToSpaces:
		entry = strings.Replace(entry, "_", " ", -1)
	case TrueTitle:
		entry = strings.Title(strings.ToLower(entry))
	case AllCaps:
		entry = strings.ToUpper(entry)
	case UUID:
		entry = strings.ToUpper(entry)
	case SpacedName:
		entry = strings.Title(strings.ToLower(strings.Replace(entry, "_", " ", -1)))
	default:
		// none
	}
	return entry
}

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
		responses.SendRes(w, responses.DB_Get_Failure, nil, getErrorMsg)
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

	return true, thisUser, userInfo
}