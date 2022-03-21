// Package handlers provides functions for handling web routes
package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"apricate/auth"
	"apricate/log"
	"apricate/metrics"
	"apricate/rdb"
	"apricate/responses"
	"apricate/schema"
	"apricate/tokengen"
)

// Handler Functions

// Handler function for the route: /
func Homepage(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- Homepage --"))
	responses.SendRes(w, responses.Unimplemented, nil, "Homepage")
	log.Debug.Println(log.Cyan("-- End Homepage --"))
}

// Handler function for the route: /api/users
func UsersSummary(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- usersSummary --"))
	res := metrics.AssembleUsersMetrics()
	responses.SendRes(w, responses.Generic_Success, res, "")
	log.Debug.Println(log.Cyan("-- End usersSummary --"))
}

// Handler function for the route: /api/islands
type IslandsOverview struct {
	World *schema.World
}
func (h *IslandsOverview) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- IslandsOverview --"))
	res := h.World.Islands
	responses.SendRes(w, responses.Generic_Success, res, "")
	log.Debug.Println(log.Cyan("-- End IslandsOverview --"))
}

// Handler function for the route: /api/islands/{islandName}
type IslandOverview struct {
	World *schema.World
}
func (h *IslandOverview) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- IslandOverview --"))
	// Get islandName from route
	islandName := GetVarEntries(r, "islandName", SpacedName)
	res := h.World.Islands[islandName]
	responses.SendRes(w, responses.Generic_Success, res, "")
	log.Debug.Println(log.Cyan("-- End IslandOverview --"))
}

// Handler function for the route: /api/regions
type RegionsOverview struct {
	World *schema.World
}
func (h *RegionsOverview) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- RegionsOverview --"))
	res := h.World.Regions
	responses.SendRes(w, responses.Generic_Success, res, "")
	log.Debug.Println(log.Cyan("-- End RegionsOverview --"))
}

// Handler function for the route: /api/regions/{regionName}
type RegionOverview struct {
	World *schema.World
}
func (h *RegionOverview) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- RegionOverview --"))
	// Get regionName from route
	regionName := GetVarEntries(r, "regionName", SpacedName)
	res := h.World.Regions[regionName]
	responses.SendRes(w, responses.Generic_Success, res, "")
	log.Debug.Println(log.Cyan("-- End RegionOverview --"))
}

// Handler function for the route: /api/users/{username}
type UsernameInfo struct {
	Dbs *map[string]rdb.Database
}
func (h *UsernameInfo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- usernameInfo --"))
	// Get username from route
	username := GetVarEntries(r, "username", None)
	log.Debug.Printf("UsernameInfo Requested for: %s", username)
	// Get username info from DB
	token, genTokenErr := tokengen.GenerateToken(username)
	if genTokenErr != nil {
		// fail state
		log.Important.Printf("in UsernameInfo: Attempted to generate token using username %s but was unsuccessful with error: %v", username, genTokenErr)
		genErrorMsg := fmt.Sprintf("Could not get, failed to convert username to DB token. Username: %v | GenerateTokenErr: %v", username, genTokenErr)
		responses.SendRes(w, responses.Generate_Token_Failure, nil, genErrorMsg)
		return
	}
	udb := (*h.Dbs)["users"]
	// Check db for user
	userData, userFound, getUserErr := schema.GetUserFromDB(token, udb)
	if getUserErr != nil {
		// fail state
		getErrorMsg := fmt.Sprintf("in publicGetUser, could not get from DB for username: %s, error: %v", username, getUserErr)
		responses.SendRes(w, responses.DB_Get_Failure, nil, getErrorMsg)
		return
	}
	if !userFound {
		// fail state - user not found
		userNotFoundMsg := fmt.Sprintf("in publicGetUser, no user found in DB with username: %s", username)
		responses.SendRes(w, responses.User_Not_Found, nil, userNotFoundMsg)
		return
	}
	// success state
	resData := schema.PublicInfo{
		Username: userData.Username,
		Title: userData.Title,
		Ledger: userData.Ledger,
		Achievements: userData.Achievements,
		UserSince: userData.UserSince,
	}
	responses.SendRes(w, responses.Generic_Success, resData, "")
	log.Debug.Println(log.Cyan("-- End usernameInfo --"))
}

// Handler function for the route: /api/users/{username}/claim
type UsernameClaim struct {
	Dbs *map[string]rdb.Database
	SlurFilter *[]string
}
func (h *UsernameClaim) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- usernameClaim --"))
	log.Debug.Println("Recover udb from context")
	udb := (*h.Dbs)["users"]
	// Get username from route
	username := GetVarEntries(r, "username", None)
	log.Debug.Printf("Username Requested: %s", username)
	// Validate username (length & content, plus characters)
	usernameValidationStatus := auth.ValidateUsername(username, h.SlurFilter)
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
		responses.SendRes(w, responses.DB_Get_Failure, nil, dbGetErrorMsg)
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
	newUser := schema.NewUser(token, username, *h.Dbs, false)
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

// Handler function for the route: /api/plants
type PlantsOverview struct {
	MainDictionary *schema.MainDictionary
}
func (h *PlantsOverview) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- PlantsOverview --"))
	res := h.MainDictionary.Plants
	responses.SendRes(w, responses.Generic_Success, res, "")
	log.Debug.Println(log.Cyan("-- End PlantsOverview --"))
}

// Handler function for the route: /api/plants/{plantName}
type PlantOverview struct {
	MainDictionary *schema.MainDictionary
}
func (h *PlantOverview) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- PlantOverview --"))
	// Get username from route
	plant_name := GetVarEntries(r, "plantName", SpacedName)
	log.Debug.Printf("PlantOverview Requested for: %s", plant_name)
	// Get plant
	if plant, ok := (*h.MainDictionary).Plants[plant_name]; ok {
		res := plant
		responses.SendRes(w, responses.Generic_Success, res, "")
	} else {
		responses.SendRes(w, responses.Specified_Plant_Not_Found, nil, "")
	}
	log.Debug.Println(log.Cyan("-- End PlantOverview --"))
}

// Handler function for the route: /api/plants/{plantName}/stage/{stageNum}
type PlantStageOverview struct {
	MainDictionary *schema.MainDictionary
}
func (h *PlantStageOverview) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- PlantStageOverview --"))
	// Get plant_name from route
	plant_name := GetVarEntries(r, "plantName", SpacedName)
	stageNumRaw := GetVarEntries(r, "stageNum", None)
	stage_num, err := strconv.Atoi(stageNumRaw)
	if err != nil {
		errmsg := fmt.Sprintf("PlantStageOverview Requested for: %s, stagenum: %d but failed to parse stageNum to Int for reason: %v", plant_name, stage_num, err)
		log.Debug.Printf(errmsg)
		responses.SendRes(w, responses.Could_Not_Parse_URI_Param, nil, errmsg)
		return
	}
	log.Debug.Printf("PlantStageOverview Requested for: %s, stagenum: %d", plant_name, stage_num)
	// Get plant
	if plant, ok := (*h.MainDictionary).Plants[plant_name]; ok {
		res := plant
		responses.SendRes(w, responses.Generic_Success, res.GrowthStages[stage_num], "")
	} else {
		responses.SendRes(w, responses.Specified_Plant_Not_Found, nil, "")
	}
	log.Debug.Println(log.Cyan("-- End PlantStageOverview --"))
}

// Handler function for the route: /api/metrics
func MetricsOverview(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- MetricsOverview --"))
	res := map[string]interface{}{"Global Market Buy/Sell": metrics.TrackingMarket, "User Coins": metrics.TrackingUserCoins}
	responses.SendRes(w, responses.Generic_Success, res, "")
	log.Debug.Println(log.Cyan("-- End MetricsOverview --"))
}