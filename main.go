package main

import (
	"fmt"
	"net/http"
	
	"apricate/handlers"
	"apricate/log"

	"github.com/gorilla/mux"
)

func main() {
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
	// mxr.HandleFunc("/api/v0/users/{username}/claim", handlers.UsernameClaim).Methods("POST")
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
	ListenPort := ":50250"
	log.Info.Printf("Listening on %s", ListenPort)
	log.Error.Fatal(http.ListenAndServe(ListenPort, mxr))
}