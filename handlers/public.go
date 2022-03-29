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

// Handler function for the route: /api/about
func AboutSummary(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- AboutSummary --"))
	res := schema.AboutInfo{
		Name: "About Summary",
		Description: "Describes the various available /about/... endpoints available",
		Information: map[string]interface{}{
			"about/sizes": "Describes the Size enum used by the game for plots, plants, and produce. The information field is the map of string to integer size.",
			"about/magic": "Describes various parts of the magic system including lore, arcane flux, and distortion. The information field is a map of string topic to string information.",
			"about/plants": "Describes plants and the stages they go through during growth, as well as the actions used to advance them through these stages",
			"about/world": "Describes the lore of the world, the player, and the mechanical grouping of locations",
		},
	}
	responses.SendRes(w, responses.Generic_Success, res, "")
	log.Debug.Println(log.Cyan("-- End AboutSummary --"))
}

// Handler function for the route: /api/about/sizes
func AboutSizes(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- AboutSizes --"))
	res := schema.AboutInfo{
		Name: "About Sizes",
		Description: "Describes the Size enum used by the game for plots, plants, and produce. The information field is the map of string size to integer size.",
		Information: map[string]interface{}{
			"Miniature": 1,
			"Tiny": 2,
			"Small": 4,
			"Modest": 8,
			"Average": 16,
			"Large": 32,
			"Huge": 64,
			"Gigantic": 256,
			"Colossal": 1024,
			"Titanic": 4096,
		},
	}
	responses.SendRes(w, responses.Generic_Success, res, "")
	log.Debug.Println(log.Cyan("-- End AboutSizes --"))
}

// Handler function for the route: /api/about/magic
func AboutMagic(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- AboutMagic --"))
	res := schema.AboutInfo{
		Name: "About Magic",
		Description: "Describes various parts of the magic system including lore, arcane flux, and distortion. The information field is a map of string topic to string information.",
		Information: map[string]interface{}{
			"Lore": "The Lattice is the magical fabric of the universe. Mages like the player are like pins that can distort the lattice around them in nth dimensional space to cause various material effects.",
			"Rituals": "Rituals are the spells of this world. The requirements for their casting are defined by a corresponding Rite.",
			"Rites": "Rites are the instructions for conducting rituals. These instructions include the required buildings, any necessary component materials/currencies, and describe the change in arcane flux effected by the ritual, as well as the required min/max distortion tiers for casting. A mage of too high a distortion tier (i.e. with too much arcane flux) cannot modulate their power low enough to cast low tier spells, and similarly cannot cast certain spells without meeting their minimum tier.",
			"Arcane Flux": "Arcane Flux is a measure of magic power available to a given mage. Rituals may add or remove various amounts of flux. Arcane Flux is bounded between 1 and 1 billion inclusive.",
			"Distortion (Tier)": "Distortion, or Distortion Tier is a metric for the power level of a given mage or spell. It is equal to Log10(flux). Distortion is bounded between 0 and 9 inclusive.",
			"Lattice Interference Rejection": "Casting rituals causes interference in the Lattice that prevents the mage from further magic for a given time depending on the ritual.",
		},
	}
	responses.SendRes(w, responses.Generic_Success, res, "")
	log.Debug.Println(log.Cyan("-- End AboutMagic --"))
}

// Handler function for the route: /api/about/plants
func AboutPlants(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- AboutPlants --"))
	res := schema.AboutInfo{
		Name: "About Plants",
		Description: "Describes plants and the stages they go through during growth, as well as the actions used to advance them through these stages",
		Information: map[string]interface{}{
			"Plants": "Plants have very distinct stages of growth, each with a specific action that must be taken by the player to advance it to the next stage. Plant definitions can be queried which have all the info on individual plants, including min/max size and their growth stages.",
			"Sizes": "Plants can often be planted in several 'sizes' which multiply both any consumable costs to grow the plant, as well as the plant's Yield modifier, which improves harvests.",
			"Yield": "The Yield modifier starts at 1.0 and can typically be improved through optional growth actions or using better consumables like fertilizer. The Yield modifier is used with size to calculate harvest chance and quantity. Seeds are not affected by yield or size, goods are affected by both yield and size, produce is affected only by yield.",
			"Growth Actions": "Growth actions are how you advance a plant between growth stages. You must have the corresponding tool to use every action (except the Wait and Skip actions, which are tool-less).",
			"Growth Stages": "Every plant has several growth stages that must be advanced between using growth actions. The plant definition describes these, including in-lore name/description, action name, and whether the stage added to the Yield modifier.",
			"Growth Time": "Most growth stages specify a growth time, which is the cooldown on using another growth action.",
			"Consumables": "Some stages specify a list of consumable options. If consumable options are present, one of the options MUST be specified in the growth action (unless the stage is optional). Higher quality consumables typically add to the Yield modifier when used, and/or may be required in lesser quantities. Consumable quantity is multiplied by plant size, and is always defined with respect to the miniature (1) size in the plant definition.",
			"Harvestable Stages": "Some stages specify a harvestable object, which explains what can be harvested, and likelihood of harvesting the given item. As seeds are not affected by yield nor size, the number here is exactly the number of seeds you will get per plant (decimal dropped after multiplying by quantity, e.g. a plot with 3 plants and 0.5 seed chance is 1.5, but only gives 1 seed).",
			"Optional Stages": "Some stages have the 'skippable' property set to True, meaning they are optional. To skip an optional stage, simply specify the Skip action. Skipping a stage means no consumables are used, nothing is harvested, and there is no growth time cooldown (you may immediately send the next growth action).",
			"Repeatable Stages": "Some stages have the 'repeatable' property set to True, meaning the plant will not advance to the next stage when the growth action is sent. On each recurrence: any consumables are used, the plant is harvested if harvestable, the growth time cooldown applies, and any yield improvements are added to the Yield modifier. If the action is skippable, you may escape the recurrence loop with the Skip action. If the action is not skippable, you must clear the plot to get rid of the plant.",
		},
	}
	responses.SendRes(w, responses.Generic_Success, res, "")
	log.Debug.Println(log.Cyan("-- End AboutPlants --"))
}

// Handler function for the route: /api/about/world
func AboutWorld(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- AboutWorld --"))
	res := schema.AboutInfo{
		Name: "About World",
		Description: "Describes the lore of the world, the player, and the mechanical grouping of locations",
		Information: map[string]interface{}{
			"Name": "Astrid",
			"Lore": "The world of Astrid is a high-magic fantasy world in the pre-industrial feudal era that has been constantly embroiled in large wars. In the recent past, a great magic catastrophy known as The Fracturing split apart the continents into thousands of small islands. Magical storms now blanket the seas, leaving few navigable routes.",
			"The Player": "The player is an artificial life form created by an Arch Mage to serve in the war. After winning it for your masters, you grew tired of fighting and retired to a quiet farm on a backwater island called Pria. Here you plan to spend you days bringing life into the world for a chance, growing plants of mundane and magical origins, crafting potions and other wares to sell in neighboring towns, and so forth.",
			"Locations": "Locations are the most local type of... location. They have markets and NPCs, as well as a warehouse for the player (once they need it).",
			"Islands": "Islands are the next step up from locations, holding several. These are connected to each other by Ports.",
			"Ports": "Ports connect one island to another in a 1-1 map. Port travel has a set travel time and a fare cost in Coins.",
			"Regions": "Regions are the next step up from islands, holding several. Regions are more of a conceptual designation, and are not explicitly separated. Travel between regions occurs via island ports as typical.",
			"Shatteres": "Shatteres are the next step up from regions, holding several. Shatteres are conceptual designations, and are not explicitly separated. Travel between shatters occurs via island port as typical.",
			"The Central Wheel": "A large shattere composed of several regions, most of the powerful nations have capitals here. The Central Wheel is so named because the islands it represents are connected in a large loop around the equator.",
			"The Nevish Extremities": "A small shattere composed of a few spur regions north of the Central Wheel shattere. In contrast to those in the Central Wheel, the islands here are generally self-governing without overarching nations. The only thing keeping them relatively independent of the fighting down south is the Treaty of Neversia, which binds all islands in this shattere in mutual defense.",
			"The Tyldian Spur": "The Region the player starts in. It is part of the Nevish Extremities shattere, connected by Tyldia to the Neversian Bulwark region.",
		},
	}
	responses.SendRes(w, responses.Generic_Success, res, "")
	log.Debug.Println(log.Cyan("-- End AboutWorld --"))
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

// Handler function for the route: /api/islands/{island-symbol}
type IslandOverview struct {
	World *schema.World
}
func (h *IslandOverview) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- IslandOverview --"))
	// Get island_symbol from route
	island_symbol := GetVarEntries(r, "island-symbol", AllCaps)
	log.Debug.Printf("Island Overview For: %s", island_symbol)
	res, ok := h.World.Islands[island_symbol]
	if !ok {
		responses.SendRes(w, responses.Location_Not_Found, nil, "")
		log.Debug.Println(log.Cyan("-- End IslandOverview --"))
		return
	}
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

// Handler function for the route: /api/regions/{region-symbol}
type RegionOverview struct {
	World *schema.World
}
func (h *RegionOverview) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- RegionOverview --"))
	// Get region-symbol from route
	region_symbol := GetVarEntries(r, "region-symbol", AllCaps)
	log.Debug.Printf("Region Overview For: %s", region_symbol)
	res, ok := h.World.Regions[region_symbol]
	if !ok {
		responses.SendRes(w, responses.Location_Not_Found, nil, "")
		log.Debug.Println(log.Cyan("-- End RegionOverview --"))
		return
	}
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
		ArcaneFlux: userData.ArcaneFlux,
		DistortionTier: userData.DistortionTier,
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

// Handler function for the route: /api/plants/{plant-name}
type PlantOverview struct {
	MainDictionary *schema.MainDictionary
}
func (h *PlantOverview) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- PlantOverview --"))
	// Get username from route
	plant_name := GetVarEntries(r, "plant-name", SpacedName)
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

// Handler function for the route: /api/plants/{plant-name}/stage/{stageNum}
type PlantStageOverview struct {
	MainDictionary *schema.MainDictionary
}
func (h *PlantStageOverview) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- PlantStageOverview --"))
	// Get plant_name from route
	plant_name := GetVarEntries(r, "plant-name", SpacedName)
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
	res := metrics.GetMetricsResponse()
	responses.SendRes(w, responses.Generic_Success, res, "")
	log.Debug.Println(log.Cyan("-- End MetricsOverview --"))
}

// Handler function for the route: /api/rites
type RitesOverview struct {
	MainDictionary *schema.MainDictionary
}
func (h *RitesOverview) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- RitesOverview --"))
	res := h.MainDictionary.Rites
	responses.SendRes(w, responses.Generic_Success, res, "")
	log.Debug.Println(log.Cyan("-- End RitesOverview --"))
}

// Handler function for the route: /api/rites/{runic-symbol}
type RiteOverview struct {
	MainDictionary *schema.MainDictionary
}
func (h *RiteOverview) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug.Println(log.Yellow("-- RiteOverview --"))
	// Get username from route
	rite_name := GetVarEntries(r, "runic-symbol", AllCaps)
	log.Debug.Printf("RiteOverview Requested for: %s", rite_name)
	// Get rite
	if rite, ok := (*h.MainDictionary).Rites[rite_name]; ok {
		res := rite
		responses.SendRes(w, responses.Generic_Success, res, "")
	} else {
		responses.SendRes(w, responses.Specified_Rite_Not_Found, nil, "")
	}
	log.Debug.Println(log.Cyan("-- End RiteOverview --"))
}