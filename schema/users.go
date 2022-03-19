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
	Contracts []string `json:"contracts" binding:"required"`
	Assistants []string `json:"assistants" binding:"required"`
	Farms []string `json:"farms" binding:"required"`
	Plots []string `json:"plots" binding:"required"`
	Warehouses []string `json:"warehouses" binding:"required"`
}

// Defines the public User info for the /users/{username} endpoint
type PublicInfo struct {
	Username string `json:"username" binding:"required"`
	Title Achievement `json:"title" binding:"required"`
	Ledger Ledger `json:"ledger" binding:"required"`
	UserSince int64 `json:"user-since" binding:"required"`
	Achievements []Achievement `json:"achievements" binding:"required"`
}

func NewUser(token string, username string, dbs map[string]rdb.Database) *User {
	// starting location
	startLocation := "TS-PR-HF"
	// generate starting assistant
	assistant := NewAssistant(username, 0, Hireling, startLocation)
	SaveAssistantToDB(dbs["assistants"], assistant)
	// generate starting farm
	farm := NewFarm(dbs["plots"], 0, username, startLocation)
	SaveFarmToDB(dbs["farms"], farm)
	// generate starting contract
	contract := NewContract(username, 0, startLocation, ContractType_Talk, "Viridis", []ContractTerms{{NPC: "Reldor"}}, []ContractReward{{RewardType: RewardType_Currency, Item: "Coins", Quantity: 100}})
	SaveContractToDB(dbs["contracts"], contract)
	// generate starting warehouse
	warehouse := NewWarehouse(username, startLocation, map[ToolTypes]uint8{Tool_Spade: 1, Tool_Shears: 1, Tool_Sickle: 1, Tool_Hoe: 1, Tool_SproutingPot: 1, Tool_Pitchfork: 1, Tool_WaterWand: 1, Tool_Rake: 1}, map[string]Produce{"Potato|Miniature": *NewProduce("Potato", Miniature, uint64(10))}, map[string]uint64{"Cabbage Seeds":10,"Shelvis Fig Seeds":10,"Potato Chunk":10,"Spectral Grass Seeds": 16,"Gulb Bulb": 4}, map[string]uint64{"Fertilizer":100, "Water": 100, "Enchanted Water": 100})
	SaveWarehouseToDB(dbs["warehouses"], warehouse)
	//TODO: generate each of these
	var starting_farm_id string = farm.UUID
	var starting_farm_warehouse_id string = warehouse.UUID
	var starting_contract_id string = contract.UUID
	var starting_assistant_id string = assistant.UUID

	plotIds := make([]string, 0)
	for _, plot := range farm.Plots {
		plotIds = append(plotIds, plot.UUID)
	}

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
			Achievements: []Achievement{Achievement_Noob},
			UserSince: time.Now().Unix(),
		},
		Contracts: []string{starting_contract_id},
		Farms: []string{starting_farm_id},
		Plots: plotIds,
		Warehouses: []string{starting_farm_warehouse_id},
		Assistants: []string{starting_assistant_id},
	}
}

// Check DB for existing user with given token and return bool for if exists, and error if error encountered
func CheckForExistingUser (token string, tdb rdb.Database) (bool, error) {
	// Get user
	_, getError := tdb.GetJsonData(token, ".")
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
func GetUserFromDB (token string, tdb rdb.Database) (User, bool, error) {
	// Get user json
	someJson, getError := tdb.GetJsonData(token, ".")
	if getError != nil {
		if fmt.Sprint(getError) == "redis: nil" {
			// user not found
			return User{}, false, nil
		}
		// error
		return User{}, false, getError
	}
	// Got successfully, unmarshal
	someData := User{}
	unmarshalErr := json.Unmarshal(someJson, &someData)
	if unmarshalErr != nil {
		log.Error.Fatalf("Could not unmarshal user json from DB: %v", unmarshalErr)
		return User{}, false, unmarshalErr
	}
	return someData, true, nil
}

// Get userdata at path from DB, bool is user found
func GetUserDataAtPathFromDB (token string, path string, tdb rdb.Database) (interface{}, bool, error) {
	// Get user json
	someJson, getError := tdb.GetJsonData(token, path)
	if getError != nil {
		if fmt.Sprint(getError) == "redis: nil" {
			// user not found
			return nil, false, nil
		}
		// error
		return nil, false, getError
	}
	// Got successfully, unmarshal
	var someData interface{}
	unmarshalErr := json.Unmarshal(someJson, &someData)
	if unmarshalErr != nil {
		log.Error.Fatalf("Could not unmarshal user json from DB: %v", unmarshalErr)
		return nil, false, unmarshalErr
	}
	return someData, true, nil
}

// Get user from DB by username, bool is user found
func GetUserByUsernameFromDB(username string, tdb rdb.Database) (User, bool, error) {
	token, tokenErr := tokengen.GenerateToken(username)
	if tokenErr != nil {
		return User{}, false, tokenErr
	}
	return GetUserFromDB(token, tdb)
}

// Attempt to save user, returns error or nil if successful
func SaveUserToDB(tdb rdb.Database, userData *User) error {
	log.Debug.Printf("Saving user %s to DB", userData.Username)
	err := tdb.SetJsonData(userData.Token, ".", userData)
	// creationSuccess := rdb.CreateUser(tdb, username, token, 0)
	return err
}

// Attempt to save user data at path, returns error or nil if successful
func SaveUserDataAtPathToDB(tdb rdb.Database, token string, path string, newValue interface{}) error {
	log.Debug.Printf("Saving user data at path %s to DB for token %s", path, token)
	err := tdb.SetJsonData(token, path, newValue)
	return err
}