// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"apricate/filemngr"
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
	Caravans []string `json:"caravans" binding:"required"`
	Farms []string `json:"farms" binding:"required"`
	Plots []string `json:"plots" binding:"required"`
	Warehouses []string `json:"warehouses" binding:"required"`
	LatticeInterferenceRejectionEnd int64 `json:"lattice_interference_rejection_end" binding:"required"`
}

// Defines the public User info for the /users/{username} endpoint
type PublicInfo struct {
	Username string `json:"username" binding:"required"`
	Title Achievement `json:"title" binding:"required"`
	Ledger Ledger `json:"ledger" binding:"required"`
	ArcaneFlux float64 `json:"arcane_flux" binding:"required"`
	DistortionTier float64 `json:"distortion_tier" binding:"required"`
	UserSince int64 `json:"user-since" binding:"required"`
	Achievements []Achievement `json:"achievements" binding:"required"`
}

func ConvertFluxToDistortion(flux float64) float64 {
	return math.Floor(math.Log10(flux) * 100) / 100
}

func NewUser(token string, username string, dbs map[string]rdb.Database, devUser bool) *User {
	// starting location
	startLocation := "TS-PR-HF"
	// generate starting assistant
	assistant := NewAssistant(username, 0, Imp, startLocation)
	assistant2 := NewAssistant(username, 1, Familiar, startLocation)
	SaveAssistantToDB(dbs["assistants"], assistant)
	SaveAssistantToDB(dbs["assistants"], assistant2)
	// generate starting farm
	farm := NewFarm(dbs["plots"], 0, username, startLocation)
	SaveFarmToDB(dbs["farms"], farm)
	// generate starting contract
	contract := NewContract(username, 0, startLocation, ContractType_Talk, "Viridis", []ContractTerms{{NPC: "Reldor"}}, []ContractReward{{RewardType: RewardType_Currency, Item: "Coins", Quantity: 100}})
	SaveContractToDB(dbs["contracts"], contract)
	// generate starting warehouse
	var warehouse *Warehouse
	if devUser {
		warehouse = NewWarehouse(username, startLocation, map[string]uint64{"Spade": 1, "Hoe": 1, "Rake": 1, "Pitchfork": 1, "Shears": 1, "Water Wand": 1, "Knife": 1, "Pestle and Mortar": 1, "Drying Rack": 1, "Sprouting Pot": 1, "Scroll of Hyperspecific Cloud Cover": 1, "Sickle": 1, "Spirit Flute": 1, "Scroll of Bind Evil": 1}, map[string]uint64{"Potato|Tiny": uint64(1000)}, map[string]uint64{"Cabbage Seeds":1000,"Shelvis Fig Seeds":1000,"Potato Chunk":1000,"Spectral Grass Seeds":1000,"Gulb Bulb":1000,"Spinosus Vas Seeds":1000,"Convocare Bulb":1000,"Uona Spore":1000,"Grape Seeds":1000}, map[string]uint64{"Fertilizer":1000, "Enchanted Fertilizer": 1000, "Dragon Fertilizer": 1000, "Enchanted Dragon Fertilizer": 1000, "Water": 1000, "Enchanted Water": 1000, "Vocatus Blossom": 1000, "Vocatus Blossom In Perfect Bloom": 1000, "Wagyu Fungus Steak": 1000})
	} else {
		warehouse = NewWarehouse(username, startLocation, map[string]uint64{"Shears": 1, "Sickle": 1}, map[string]uint64{}, map[string]uint64{"Cabbage Seeds":8,"Potato Chunk":4,"Spectral Grass Seeds":16}, map[string]uint64{})
	}
	SaveWarehouseToDB(dbs["warehouses"], warehouse)
	//TODO: generate each of these
	var starting_farm_id string = farm.UUID
	var starting_farm_warehouse_id string = warehouse.UUID
	var starting_contract_id string = contract.UUID
	var starting_assistant_id string = assistant.UUID
	var starting_assistant_2_id string = assistant2.UUID

	var starting_currencies map[string]uint64
	var starting_favor map[string]int8

	var startingFlux float64
	if devUser {
		startingFlux = 2200
		starting_currencies = map[string]uint64{"Coins": 1000}
		starting_favor = map[string]int8{"Vince Kosuga": 50}
	} else {
		startingFlux = 10
		starting_currencies = map[string]uint64{"Coins": 100}
		starting_favor = map[string]int8{"Vince Kosuga": 50}
	}

	plotIds := make([]string, 0)
	for _, plot := range farm.Plots {
		plotIds = append(plotIds, plot.UUID)
	}

	TrackUserMagic(username, startingFlux, ConvertFluxToDistortion(startingFlux))

	return &User{
		Token: token,
		PublicInfo: PublicInfo{
			Username: username,
			Title: Achievement_Noob,
			Ledger: Ledger{
				Currencies: starting_currencies,
				Favor: starting_favor,
				Escrow: make(map[string]uint64),
			},
			ArcaneFlux: startingFlux,
			DistortionTier: ConvertFluxToDistortion(startingFlux),
			Achievements: []Achievement{Achievement_Noob},
			UserSince: time.Now().Unix(),
		},
		LatticeInterferenceRejectionEnd: 0,
		Contracts: []string{starting_contract_id},
		Farms: []string{starting_farm_id},
		Plots: plotIds,
		Warehouses: []string{starting_farm_warehouse_id},
		Assistants: []string{starting_assistant_id, starting_assistant_2_id},
		Caravans: make([]string, 0),
	}
}

func PregenerateUser(username string, dbs map[string]rdb.Database, devuser bool) {
	// generate token
	token, genTokenErr := tokengen.GenerateToken(username)
	if genTokenErr != nil {
		// fail state
		log.Important.Printf("in UsernameClaim: Attempted to generate token using username %s but was unsuccessful with error: %v", username, genTokenErr)
		genErrorMsg := fmt.Sprintf("Username: %v | GenerateTokenErr: %v", username, genTokenErr)
		panic(genErrorMsg)
	}
	// create new user in DB
	newUser := NewUser(token, username, dbs, devuser)
	newUser.Title = Achievement_Owner
	newUser.Achievements = []Achievement{Achievement_Owner, Achievement_Contributor, Achievement_Noob}
	saveUserErr := SaveUserToDB(dbs["users"], newUser)
	if saveUserErr != nil {
		// fail state - could not save
		saveUserErrMsg := fmt.Sprintf("in UsernameClaim | Username: %v | CreateNewUserInDB failed, dbSaveResult: %v", username, saveUserErr)
		log.Debug.Println(saveUserErrMsg)
		panic(saveUserErrMsg)
	}
	// Write out my token
	lines, readErr := filemngr.ReadFileToLineSlice("data/secrets.env")
	if readErr != nil {
		// Auth is mission-critical, using Fatal
		log.Error.Fatalf("Could not read lines from secrets.env. Err: %v", readErr)
	}
	secretIdentifier := strings.ToUpper(username) + "_TOKEN="
	secretString :=  secretIdentifier + string(token)
	// Search existing file for secret identifier
	found, i := filemngr.KeyInSliceOfLines(secretIdentifier, lines)
	if found {
		// Update existing secret
		lines [i] = secretString
	} else {
		// Create secret in env file since could not find one to update
		// If empty file then replace 1st line else append to end
		log.Debug.Printf("Creating new secret in env file. secrets.env[0] == ''? %v", lines[0] == "")
		if lines[0] == "" {
			log.Debug.Printf("Blank secrets.env, replacing line 0")
			lines[0] = secretString
		} else {
			log.Debug.Printf("Not blank secrets.env, appending to end")
			lines = append(lines, secretString)
		}
	}
	
	// Join and write out
	writeErr := filemngr.WriteLinesToFile("data/secrets.env", lines)
	if writeErr != nil {
		log.Error.Fatalf("Could not write secrets.env: %v", writeErr)
	}
	log.Info.Printf("Wrote token for user: %s to secrets.env", username)
	// Created successfully
	// Track in user metrics
	// metrics.TrackNewUser(username) // CANT IN SCHEMA MOD CAUSE IMPORT CYCLE
	log.Debug.Printf("Generated token %s and claimed username %s", token, username)
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
	TrackUserCoins(userData.Username, userData.Ledger.Currencies["Coins"])
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