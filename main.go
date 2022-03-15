package main

import (
	"fmt"
	"net/http"
	"strings"

	"apricate/auth"
	"apricate/filemngr"
	"apricate/handlers"
	"apricate/log"
	"apricate/metrics"
	"apricate/rdb"
	"apricate/schema"
	"apricate/tokengen"

	"github.com/gorilla/mux"
)

// Global Vars
var (
	ListenPort = ":50250"
	RedisAddr = "localhost:6382"
	apiVersion = "0.2.0"
	// Define relationship between string database name and redis db
	dbs = make(map[string]rdb.Database)
	world schema.World
	flush_DBs = true
	regenerate_auth_secret = false
)

func initialize_dbs() {
	dbs["users"] = rdb.NewDatabase(RedisAddr, 0)
	dbs["assistants"] = rdb.NewDatabase(RedisAddr, 1)
	dbs["farms"] = rdb.NewDatabase(RedisAddr, 2)
	dbs["contracts"] = rdb.NewDatabase(RedisAddr, 3)
	dbs["warehouses"] = rdb.NewDatabase(RedisAddr, 4)
	dbs["clearinghouse"] = rdb.NewDatabase(RedisAddr, 5)
	dbs["plots"] = rdb.NewDatabase(RedisAddr, 6)

	if flush_DBs || regenerate_auth_secret {
		for _, db := range dbs {
			db.Flush()
		}
	}
}

func setup_my_character() {
	if flush_DBs || regenerate_auth_secret {
		username := "Greenitthe"
		// generate token
		token, genTokenErr := tokengen.GenerateToken(username)
		if genTokenErr != nil {
			// fail state
			log.Important.Printf("in UsernameClaim: Attempted to generate token using username %s but was unsuccessful with error: %v", username, genTokenErr)
			genErrorMsg := fmt.Sprintf("Username: %v | GenerateTokenErr: %v", username, genTokenErr)
			panic(genErrorMsg)
		}
		// create new user in DB
		newUser := schema.NewUser(token, username, dbs)
		newUser.Title = schema.Achievement_Owner
		newUser.Achievements = []schema.Achievement{schema.Achievement_Owner, schema.Achievement_Contributor, schema.Achievement_Noob}
		saveUserErr := schema.SaveUserToDB(dbs["users"], newUser)
		if saveUserErr != nil {
			// fail state - could not save
			saveUserErrMsg := fmt.Sprintf("in UsernameClaim | Username: %v | CreateNewUserInDB failed, dbSaveResult: %v", username, saveUserErr)
			log.Debug.Println(saveUserErrMsg)
			panic(saveUserErrMsg)
		}
		// Write out my token
		lines, readErr := filemngr.ReadFileToLineSlice("secrets.env")
		if readErr != nil {
			// Auth is mission-critical, using Fatal
			log.Error.Fatalf("Could not read lines from secrets.env. Err: %v", readErr)
		}
		secretString :=  "GREENITTHE_TOKEN=" + string(token)
		// Search existing file for secret identifier
		found, i := filemngr.KeyInSliceOfLines("GREENITTHE_TOKEN=", lines)
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
		writeErr := filemngr.WriteLinesToFile("secrets.env", lines)
		if writeErr != nil {
			log.Error.Fatalf("Could not write secrets.net: %v", writeErr)
		}
		log.Info.Println("Wrote token for user: Greenitthe to secrets.env")
		// Created successfully
		// Track in user metrics
		metrics.TrackNewUser(username)
		log.Debug.Printf("Generated token %s and claimed username %s", token, username)
	}
	log.Info.Println("Neither flushing DBs, nor regenerating auth secret. Token for user: Greenitthe should already exist in secrets.env. Skipping creation")
}

func main() {
	log.Info.Printf("Guild-Golems Rest API Server %s", apiVersion)
	log.Info.Printf("Connecting to Redis DB")

	// Setup redis databases for each namespace
	initialize_dbs()

	// Handle auth secret generation if requested
	if regenerate_auth_secret {
		log.Important.Printf("(Re)Generating Auth Secret")
		auth.CreateOrUpdateAuthSecretInFile()
	}

	log.Info.Println("Loading secrets from envfile")
	auth.LoadSecretsToEnv()

	setup_my_character()

	// Preload 
	// Ensure exists
	filemngr.Touch("slur_filter.txt")
	// Read file to lines array splitting by newline
	read_slur_filter, readErr := filemngr.ReadFileToLineSlice("slur_filter.txt")
	if readErr != nil {
		// Auth is mission-critical, using Fatal
		log.Error.Fatalf("Could not read lines from slur_filter.txt. Err: %v", readErr)
	}
	slur_filter := make([]string, len(read_slur_filter))
	for i, word := range read_slur_filter {
		slur_filter[i] = strings.ToUpper(word)
	}
	log.Info.Printf("Created/Loaded Username Slur Filter")

	// Load World from YAML
	world = schema.World_load("./yaml/regions.yaml", "./yaml/locations")
	log.Debug.Println(world)
	log.Info.Printf("Loaded world")

	// Begin Serving
	handle_requests(slur_filter)
}

func handle_requests(slur_filter []string) {
	// Define Routes
	//mux router
	mxr := mux.NewRouter().StrictSlash(true)
	// mxr.Use(handlers.GenerateHandlerMiddlewareFunc(userDatabase,worldDatabase))
	mxr.HandleFunc("/", handlers.Homepage).Methods("GET")
	// mxr.HandleFunc("/api", handlers.ApiSelection).Methods("GET")
	// mxr.HandleFunc("/api/v0", handlers.V0Status).Methods("GET")
	// mxr.HandleFunc("/api/v0/leaderboards", handlers.LeaderboardDescriptions).Methods("GET")
	// mxr.HandleFunc("/api/v0/leaderboards/{board}", handlers.GetLeaderboards).Methods("GET")
	mxr.HandleFunc("/api/users", handlers.UsersSummary).Methods("GET")
	mxr.Handle("/api/users/{username}", &handlers.UsernameInfo{Dbs: &dbs}).Methods("GET")
	mxr.Handle("/api/users/{username}/claim", &handlers.UsernameClaim{Dbs: &dbs, SlurFilter: &slur_filter}).Methods("POST")
	mxr.Handle("/api/regions", &handlers.RegionsOverview{World: &world}).Methods("GET")

	// secure subrouter for account-specific routes
	secure := mxr.PathPrefix("/api/my").Subrouter()
	secure.Use(auth.GenerateTokenValidationMiddlewareFunc(dbs["users"]))
	secure.Handle("/account", &handlers.AccountInfo{Dbs: &dbs}).Methods("GET")
	secure.Handle("/assistants", &handlers.AssistantsInfo{Dbs: &dbs}).Methods("GET")
	secure.Handle("/assistants/{uuid}", &handlers.AssistantInfo{Dbs: &dbs}).Methods("GET")
	secure.Handle("/farms", &handlers.FarmsInfo{Dbs: &dbs}).Methods("GET")
	secure.Handle("/farms/{uuid}", &handlers.FarmInfo{Dbs: &dbs}).Methods("GET")
	secure.Handle("/contracts", &handlers.ContractsInfo{Dbs: &dbs}).Methods("GET")
	secure.Handle("/contracts/{uuid}", &handlers.ContractsInfo{Dbs: &dbs}).Methods("GET")
	secure.Handle("/warehouses", &handlers.WarehousesInfo{Dbs: &dbs}).Methods("GET")
	secure.Handle("/warehouses/{uuid}", &handlers.WarehouseInfo{Dbs: &dbs}).Methods("GET")
	secure.Handle("/nearby-locations", &handlers.NearbyLocationsInfo{Dbs: &dbs, World: &world}).Methods("GET")
	secure.Handle("/locations", &handlers.LocationsInfo{Dbs: &dbs, World: &world}).Methods("GET")
	secure.Handle("/locations/{name}", &handlers.LocationInfo{Dbs: &dbs, World: &world}).Methods("GET")
	// secure.HandleFunc("/inventories", handlers.InventoryInfo).Methods("GET")
	// secure.HandleFunc("/itineraries", handlers.ItineraryInfo).Methods("GET")
	// secure.HandleFunc("/markets", handlers.MarketInfo).Methods("GET")
	// secure.HandleFunc("/orders", handlers.OrderInfo).Methods("GET")
	// secure.HandleFunc("/orders/{status}", handlers.GetOrdersByStatus).Methods("GET")
	// secure.HandleFunc("/golems", handlers.GetGolems).Methods("GET")
	// secure.HandleFunc("/golems/{archetype}", handlers.GetGolemsByArchetype).Methods("GET")
	// secure.HandleFunc("/golem/{symbol}", handlers.GolemInfo).Methods("GET")
	// secure.HandleFunc("/golem/{symbol}", handlers.ChangeGolemTask).Methods("PUT")
	// secure.HandleFunc("/rituals", handlers.ListRituals).Methods("GET")
	// secure.HandleFunc("/rituals/{ritual}", handlers.GetRitualInfo).Methods("GET")
	// secure.HandleFunc("/rituals/summon-invoker", handlers.NewInvoker).Methods("POST")
	// secure.HandleFunc("/rituals/summon-harvester", handlers.NewHarvester).Methods("POST")
	// secure.HandleFunc("/rituals/summon-courier", handlers.NewCourier).Methods("POST")
	// secure.HandleFunc("/rituals/summon-merchant", handlers.NewMerchant).Methods("POST")

	// Start listening
	log.Info.Printf("Listening on %s", ListenPort)
	log.Error.Fatal(http.ListenAndServe(ListenPort, mxr))
}