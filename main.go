package main

import (
	"net/http"

	"apricate/auth"
	"apricate/handlers"
	"apricate/log"
	"apricate/rdb"

	"github.com/gorilla/mux"
)

// Global Vars
var (
	ListenPort = ":50250"
	RedisAddr = "localhost:6382"
	apiVersion = "0.1.0"
	// Define relationship between string database name and redis db
	dbs = make(map[string]rdb.Database)
	flush_DBs = true
	regenerate_auth_secret = true
)

func initialize_dbs() {
	dbs["users"] = rdb.NewDatabase(RedisAddr, 0)
	dbs["assistants"] = rdb.NewDatabase(RedisAddr, 1)
	dbs["farms"] = rdb.NewDatabase(RedisAddr, 2)
	dbs["contracts"] = rdb.NewDatabase(RedisAddr, 3)
	dbs["inventories"] = rdb.NewDatabase(RedisAddr, 4)
	dbs["clearinghouse"] = rdb.NewDatabase(RedisAddr, 5)
	dbs["plots"] = rdb.NewDatabase(RedisAddr, 6)

	if flush_DBs || regenerate_auth_secret {
		for _, db := range dbs {
			db.Flush()
		}
	}
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

	// Begin Serving
	handle_requests()
}

func handle_requests() {
	// Define Routes
	//mux router
	mxr := mux.NewRouter().StrictSlash(true)
	// mxr.Use(handlers.GenerateHandlerMiddlewareFunc(userDatabase,worldDatabase))
	mxr.HandleFunc("/", handlers.Homepage).Methods("GET")
	// mxr.HandleFunc("/api", handlers.ApiSelection).Methods("GET")
	// mxr.HandleFunc("/api/v0", handlers.V0Status).Methods("GET")
	// mxr.HandleFunc("/api/v0/leaderboards", handlers.LeaderboardDescriptions).Methods("GET")
	// mxr.HandleFunc("/api/v0/leaderboards/{board}", handlers.GetLeaderboards).Methods("GET")
	// mxr.HandleFunc("/api/v0/users", handlers.UsersSummary).Methods("GET")
	// mxr.HandleFunc("/api/v0/users/{username}", handlers.UsernameInfo).Methods("GET")
	mxr.Handle("/api/users/{username}/claim", &handlers.UsernameClaim{Dbs: &dbs}).Methods("POST")
	// mxr.HandleFunc("/api/v0/locations", handlers.LocationsOverview).Methods("GET")

	// // secure subrouter for account-specific routes
	// secure := mxr.PathPrefix("/api/v0/my").Subrouter()
	// secure.Use(auth.GenerateTokenValidationMiddlewareFunc(userDatabase))
	// secure.HandleFunc("/account", handlers.AccountInfo).Methods("GET")
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