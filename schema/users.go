// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

import (
	"encoding/json"
	"fmt"
	"time"

	"apricate/log"
	"apricate/rdb"
	"apricate/tokengen"
)

// Defines a user
type User struct {
	Token string `json:"token" binding:"required"`
	PublicInfo
}

// Defines the public User info for the /users/{username} endpoint
type PublicInfo struct {
	Username string `json:"username" binding:"required"`
	Title Achievement `json:"title" binding:"required"`
	Ledger Ledger `json:"ledger" binding:"required"`
	UserSince int64 `json:"user-since" binding:"required"`
	Achievements []Achievement `json:"achievements" binding:"required"`
	Contracts []uint64 `json:"contracts" binding:"required"`
	Assistants []uint64 `json:"assistants" binding:"required"`
	Farms []uint64 `json:"farms" binding:"required"`
	Inventories []uint64 `json:"inventories" binding:"required"`
}

func NewUser(token string, username string) *User {
	//TODO: generate each of these
	var starting_farm_id uint64 = 0
	var starting_farm_inventory_id uint64 = 0
	var starting_contract_id uint64 = 0
	var starting_assistant_id uint64 = 0

	return &User{
		Token: token,
		PublicInfo: PublicInfo{
			Username: username,
			Title: Achievement_Noob,
			Ledger: Ledger{
				Currencies: make(map[string]uint64),
				Favor: make(map[string]int8),
				Escrow: make(map[string]uint64),
			},
			Contracts: []uint64{starting_contract_id},
			Farms: []uint64{starting_farm_id},
			Inventories: []uint64{starting_farm_inventory_id},
			Assistants: []uint64{starting_assistant_id},
			Achievements: []Achievement{Achievement_Noob},
			UserSince: time.Now().Unix(),
		},
	}
}

// Check DB for existing user with given token and return bool for if exists, and error if error encountered
func CheckForExistingUser (token string, udb rdb.Database) (bool, error) {
	// Get user
	_, getError := udb.GetJsonData(token, ".")
	if getError != nil {
		if fmt.Sprint(getError) != "redis: nil" {
			// error
			return false, getError
		}
		// user not found
		return false, nil
	}
	// Got successfully
	return true, nil
}

// Get user from DB, bool is user found
func GetUserFromDB (token string, udb rdb.Database) (User, bool, error) {
	// Get user json
	uJson, getError := udb.GetJsonData(token, ".")
	if getError != nil {
		if fmt.Sprint(getError) != "redis: nil" {
			// user not found
			return User{}, false, nil
		}
		// error
		return User{}, false, getError
	}
	// Got successfully, unmarshal
	uData := User{}
	unmarshalErr := json.Unmarshal(uJson, &uData)
	if unmarshalErr != nil {
		log.Error.Fatalf("Could not unmarshal user json from DB: %v", unmarshalErr)
		return User{}, false, unmarshalErr
	}
	return uData, true, nil
}

// Get userdata at path from DB, bool is user found
func GetUserDataAtPathFromDB (token string, path string, udb rdb.Database) (interface{}, bool, error) {
	// Get user json
	uJson, getError := udb.GetJsonData(token, path)
	if getError != nil {
		if fmt.Sprint(getError) != "redis: nil" {
			// user not found
			return nil, false, nil
		}
		// error
		return nil, false, getError
	}
	// Got successfully, unmarshal
	var uData interface{}
	unmarshalErr := json.Unmarshal(uJson, &uData)
	if unmarshalErr != nil {
		log.Error.Fatalf("Could not unmarshal user json from DB: %v", unmarshalErr)
		return nil, false, unmarshalErr
	}
	return uData, true, nil
}

// Get user from DB by username, bool is user found
func GetUserByUsernameFromDB(username string, udb rdb.Database) (User, bool, error) {
	token, tokenErr := tokengen.GenerateToken(username)
	if tokenErr != nil {
		return User{}, false, tokenErr
	}
	return GetUserFromDB(token, udb)
}

// Attempt to save user, returns error or nil if successful
func SaveUserToDB(udb rdb.Database, userData *User) error {
	log.Debug.Printf("Saving user %s to DB", userData.Username)
	err := udb.SetJsonData(userData.Token, ".", userData)
	// creationSuccess := rdb.CreateUser(udb, username, token, 0)
	return err
}

// Attempt to save user data at path, returns error or nil if successful
func SaveUserDataAtPathToDB(udb rdb.Database, token string, path string, newValue interface{}) error {
	log.Debug.Printf("Saving user data at path %s to DB for token %s", path, token)
	err := udb.SetJsonData(token, path, newValue)
	return err
}