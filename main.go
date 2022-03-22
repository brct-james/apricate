package main

import (
	"net/http"
	"strings"

	"apricate/auth"
	"apricate/filemngr"
	"apricate/handlers"
	"apricate/log"
	"apricate/metrics"
	"apricate/rdb"
	"apricate/responses"
	"apricate/schema"

	"github.com/gorilla/mux"
)

// Global Vars
var (
	ListenPort = ":50250"
	RedisAddr = "localhost:6382"
	apiVersion = "0.5.0"
	// Define relationship between string database name and redis db
	dbs = make(map[string]rdb.Database)
	world schema.World
	main_dictionary = schema.MainDictionary{}
	flush_DBs = false
	regenerate_auth_secret = false
)

func initialize_dbs() {
	dbs["users"] = rdb.NewDatabase(RedisAddr, 0)
	dbs["assistants"] = rdb.NewDatabase(RedisAddr, 1)
	dbs["farms"] = rdb.NewDatabase(RedisAddr, 2)
	dbs["contracts"] = rdb.NewDatabase(RedisAddr, 3)
	dbs["warehouses"] = rdb.NewDatabase(RedisAddr, 4)
	dbs["clearinghouse"] = rdb.NewDatabase(RedisAddr, 5)

	if flush_DBs || regenerate_auth_secret {
		for _, db := range dbs {
			db.Flush()
		}
	}
}

func initialize_dictionaries() {
	// Load Seeds from YAML
	log.Debug.Println("Loading seeds list")
	main_dictionary.Seeds = schema.Seeds_load("./yaml/items/seeds.yaml")
	for k := range main_dictionary.Seeds {
		log.Debug.Println(k)
	}
	log.Info.Printf("Loaded seeds list")

	// Load Produce from YAML
	log.Debug.Println("Loading produce list")
	main_dictionary.Produce = schema.Produce_load("./yaml/items/produce.yaml")
	for k := range main_dictionary.Produce {
		log.Debug.Println(k)
	}
	log.Info.Printf("Loaded produce list")

	/*
	****** TODO: Validate 1:1 mapping for every seed and plant after loading both
	*/

	// Load Plants from YAML
	log.Debug.Println("Loading plant dictionary")
	main_dictionary.Plants = schema.Plants_load("./yaml/plants.yaml")
	for k, p := range main_dictionary.Plants {
		log.Debug.Printf("%s: %s", k, p.Name)
	}
	// log.Debug.Println(responses.JSON(dictionaries["plants"]))
	log.Info.Printf("Loaded plant dictionary")

	// Load Goods from YAML
	log.Debug.Println("Loading goods list")
	main_dictionary.Goods = schema.GoodListGenerator("./yaml/items/goods.yaml")
	log.Debug.Println(responses.JSON(main_dictionary.Goods))
	log.Info.Printf("Loaded goods list")

	// Load Markets from YAML
	log.Debug.Println("Loading Markets list")
	main_dictionary.Markets = schema.Markets_load("./yaml/world/markets.yaml")
	log.Debug.Println(responses.JSON(main_dictionary.Markets))
	log.Info.Printf("Loaded Markets list")
}

func setup_my_character() {
	if flush_DBs || regenerate_auth_secret {
		schema.PregenerateUser("Greenitthe", dbs)
		metrics.TrackNewUser("Greenitthe")
		schema.PregenerateUser("Viridis", dbs)
		metrics.TrackNewUser("Viridis")
		schema.PregenerateUser("Green", dbs)
		metrics.TrackNewUser("Green")
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

	// Reset or Load Metrics
	log.Info.Printf("Loading metrics.yaml")
	if flush_DBs || regenerate_auth_secret {
		// Need to reset metrics
		log.Important.Printf("Cleared metrics.yaml")
		filemngr.DeleteIfExists("metrics.yaml")
	}
	metrics.LoadMetrics()

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
	world = schema.World_load("./yaml/world/regions.yaml", "./yaml/world/islands", "./yaml/world/locations")
	log.Debug.Println(world)
	log.Info.Printf("Loaded world")

	// Initialize dictionaries
	initialize_dictionaries()

	// Begin Serving
	handle_requests(slur_filter)
}

func handle_requests(slur_filter []string) {
	// Define Routes
	//mux router
	mxr := mux.NewRouter().StrictSlash(true)
	// mxr.Use(handlers.GenerateHandlerMiddlewareFunc(userDatabase,worldDatabase))
	mxr.HandleFunc("/", handlers.Homepage).Methods("GET")
	mxr.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://apricate.stoplight.io/docs/apricate/YXBpOjQ1NTU3NTc2-apricate-api", http.StatusPermanentRedirect)
	})
	mxr.HandleFunc("/docs/alpha-guide", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://apricate.stoplight.io/docs/apricate/ZG9jOjQ3MDIzNTgw-alpha-guide", http.StatusPermanentRedirect)
	})
	// mxr.HandleFunc("/api", handlers.ApiSelection).Methods("GET")
	// mxr.HandleFunc("/api/leaderboards", handlers.LeaderboardDescriptions).Methods("GET")
	// mxr.HandleFunc("/api/leaderboards/{board}", handlers.GetLeaderboards).Methods("GET")
	mxr.HandleFunc("/api/users", handlers.UsersSummary).Methods("GET")
	mxr.Handle("/api/users/{username}", &handlers.UsernameInfo{Dbs: &dbs}).Methods("GET")
	mxr.Handle("/api/users/{username}/claim", &handlers.UsernameClaim{Dbs: &dbs, SlurFilter: &slur_filter}).Methods("POST")
	mxr.Handle("/api/islands", &handlers.IslandsOverview{World: &world}).Methods("GET")
	mxr.Handle("/api/islands/{island-symbol}", &handlers.IslandOverview{World: &world}).Methods("GET")
	mxr.Handle("/api/regions", &handlers.RegionsOverview{World: &world}).Methods("GET")
	mxr.Handle("/api/regions/{regionName}", &handlers.RegionOverview{World: &world}).Methods("GET")
	mxr.Handle("/api/plants", &handlers.PlantsOverview{MainDictionary: &main_dictionary}).Methods("GET")
	mxr.Handle("/api/plants/{plantName}", &handlers.PlantOverview{MainDictionary: &main_dictionary}).Methods("GET")
	mxr.Handle("/api/plants/{plantName}/stage/{stageNum}", &handlers.PlantStageOverview{MainDictionary: &main_dictionary}).Methods("GET")
	mxr.HandleFunc("/api/metrics", handlers.MetricsOverview).Methods("GET")

	// secure subrouter for account-specific routes
	secure := mxr.PathPrefix("/api/my").Subrouter()
	secure.Use(auth.GenerateTokenValidationMiddlewareFunc(dbs["users"]))
	secure.Handle("/user", &handlers.AccountInfo{Dbs: &dbs}).Methods("GET")
	secure.Handle("/assistants", &handlers.AssistantsInfo{Dbs: &dbs}).Methods("GET")
	secure.Handle("/assistants/{id}", &handlers.AssistantInfo{Dbs: &dbs}).Methods("GET")
	secure.Handle("/farms", &handlers.FarmsInfo{Dbs: &dbs}).Methods("GET")
	secure.Handle("/farms/{location-symbol}", &handlers.FarmInfo{Dbs: &dbs}).Methods("GET")
	secure.Handle("/contracts", &handlers.ContractsInfo{Dbs: &dbs}).Methods("GET")
	secure.Handle("/contracts/{id}", &handlers.ContractInfo{Dbs: &dbs}).Methods("GET")
	secure.Handle("/warehouses", &handlers.WarehousesInfo{Dbs: &dbs}).Methods("GET")
	secure.Handle("/warehouses/{location-symbol}", &handlers.WarehouseInfo{Dbs: &dbs}).Methods("GET")
	secure.Handle("/nearby-locations", &handlers.NearbyLocationsInfo{Dbs: &dbs, World: &world}).Methods("GET")
	secure.Handle("/locations", &handlers.LocationsInfo{Dbs: &dbs, World: &world}).Methods("GET")
	secure.Handle("/locations/{location-symbol}", &handlers.LocationInfo{Dbs: &dbs, World: &world}).Methods("GET")
	secure.Handle("/markets", &handlers.MarketsInfo{Dbs: &dbs, MainDictionary: &main_dictionary}).Methods("GET")
	secure.Handle("/markets/{location-symbol}", &handlers.MarketInfo{Dbs: &dbs, MainDictionary: &main_dictionary}).Methods("GET")
	secure.Handle("/markets/{location-symbol}/order", &handlers.MarketOrder{Dbs: &dbs, MainDictionary: &main_dictionary}).Methods("PATCH")
	secure.Handle("/plots", &handlers.PlotsInfo{Dbs: &dbs}).Methods("GET")
	secure.Handle("/plots/{plot-id}", &handlers.PlotInfo{Dbs: &dbs}).Methods("GET")
	secure.Handle("/plots/{plot-id}/plant", &handlers.PlantPlot{Dbs: &dbs, MainDictionary: &main_dictionary}).Methods("POST")
	secure.Handle("/plots/{plot-id}/clear", &handlers.ClearPlot{Dbs: &dbs}).Methods("PUT")
	secure.Handle("/plots/{plot-id}/interact", &handlers.InteractPlot{Dbs: &dbs, MainDictionary: &main_dictionary}).Methods("PATCH")

	// Start listening
	log.Info.Printf("Listening on %s", ListenPort)
	log.Error.Fatal(http.ListenAndServe(ListenPort, mxr))
}