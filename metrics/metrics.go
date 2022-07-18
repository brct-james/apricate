// Package metrics defines functions for tracking and displaying various game and server metrics
package metrics

import (
	"apricate/filemngr"
	"apricate/log"
	"apricate/schema"
	"apricate/timecalc"
	"fmt"
	"time"
)

// Helper Functions

// TODO: Test
// Write out metrics
func SaveMetrics() {
	log.Debug.Printf("Save metrics")
	mYaml := schema.SaveMetricsYaml{
		UniqueUsers: TrackingUniqueUsers.Usernames, // Handled by TrackUserCall
		UserActivity: TrackingActiveUsers.UserActivity, // Handled by TrackUserCall
		Coins: TrackingUserCoins.Coins, // Handled by TrackUserCall and TrackMarketBuySell
		MarketData: TrackingMarket.MarketData, // Handled by TrackMarketBuySell
		HarvestData: TrackingHarvests.HarvestData, // Handled by TrackHarvest
		RitualData: TrackingRituals.RitualData, // Handled by TrackRitual
		UserMagic: TrackingUserMagic.Magic, // Handled by TrackUserMagic
	}
	schema.Metrics_to_yaml("data/metrics.yaml", mYaml)
}

// TODO: Test
// Read out metrics
func LoadMetrics() {
	log.Debug.Printf("Load metrics")
	mYaml, found := schema.Metrics_from_yaml("data/metrics.yaml")
	if !found {
		// Failed to load
		log.Important.Printf("Failed to load metrics from YAML, saved metrics may not exist, creating and continuing.")
		filemngr.Touch("data/metrics.yaml")
		return
	}
	log.Debug.Printf("Found: %v", mYaml)
	TrackingUniqueUsers.Usernames = mYaml.UniqueUsers
	TrackingActiveUsers.UserActivity = mYaml.UserActivity
	TrackingUserCoins.Coins = mYaml.Coins
	TrackingMarket.MarketData = mYaml.MarketData
	TrackingHarvests.HarvestData = mYaml.HarvestData
	TrackingRituals.RitualData = mYaml.RitualData
	TrackingUserMagic.Magic = mYaml.UserMagic

	// If was blank, make sure it is still initialized
	if TrackingHarvests.HarvestData == nil {
		TrackingHarvests.HarvestData = make(map[string]uint64)
	}
	if TrackingMarket.MarketData == nil {
		TrackingMarket.MarketData = make(map[string]schema.GMBSMarketData)
	}
	if TrackingUserCoins.Coins == nil {
		TrackingUserCoins.Coins = make(map[string]uint64)
	}
	if TrackingActiveUsers.UserActivity == nil {
		TrackingActiveUsers.UserActivity = make(map[string]int64, 0)
	}
	if TrackingUniqueUsers.Usernames == nil {
		TrackingUniqueUsers.Usernames = make([]string, 0)
	}
	if TrackingRituals.RitualData == nil {
		TrackingRituals.RitualData = make(map[string]uint64, 0)
	}
	if TrackingUserMagic.Magic == nil {
		TrackingUserMagic.Magic = make(map[string]map[string]float64, 0)
	}
}

// TODO: Test
// Get metrics response
func GetMetricsResponse() (schema.MetricsResponse) {
	return schema.MetricsResponse {
		MarketBuySell: TrackingMarket,
		UserCoins: *TrackingUserCoins,
		Harvests: TrackingHarvests,
		Rituals: TrackingRituals,
		UserMagic: *TrackingUserMagic,
	}
}

// TODO: Test
// Assemble users metrics for json response
func AssembleUsersMetrics() (schema.UsersMetricEndpointResponse) {
	return schema.UsersMetricEndpointResponse{
		UniqueUsers:CalculateUniqueUsers(),
		ActiveUsers:CalculateActiveUsers(),
		// UsersByAchievement:CalculateUsersByAchievement(),
	}
}

// Metrics

// Unique Users
var TrackingUniqueUsers = schema.UniqueUsersMetric {
	Metric: schema.Metric{Name:"Unique Users", Description:"List of every user who has made an account since the last wipe."},
	Usernames: make([]string, 0),
}
// TODO: Test
func CalculateUniqueUsers() ([]string) {
	return TrackingUniqueUsers.Usernames
}
// TODO: Test
func TrackNewUser(username string) {
	log.Debug.Printf("Metrics:TrackNewUser")
	TrackingUniqueUsers.Usernames = append(TrackingUniqueUsers.Usernames, username)
	TrackUserCall(username)
}

// Active Users
var ActivityThresholdInMinutes int = 60
var TrackingActiveUsers = schema.ActiveUsersMetric {
	Metric: schema.Metric{Name:"Active Users", Description:fmt.Sprintf("List of every user who is considered active: have registered as a new user or hit a secure endpoint in the last %d minutes.", ActivityThresholdInMinutes)},
	UserActivity: make(map[string]int64, 0),
}
// TODO: Test
func CalculateActiveUsers() ([]string) {
	res := make([]string, 0)
	for username, timestamp := range TrackingActiveUsers.UserActivity {
		exclusion_time := timecalc.AddMinutesToTimestamp(time.Unix(timestamp, 0), ActivityThresholdInMinutes)
		if exclusion_time.After(time.Now()) {
			//include user from active users, as exclusion time in future
			res = append(res, username)
		}
	}
	return res
}
// TODO: Test
func TrackUserCall(username string) {
	log.Debug.Printf("Metrics:TrackUserCall")
	TrackingActiveUsers.UserActivity[username] = time.Now().Unix()
	SaveMetrics()
}

// User Coins
// See schema.User
var TrackingUserCoins = &schema.TrackingUserCoins
// See schema.User
var TrackingUserMagic = &schema.TrackingUserMagic

// Global Market Buy/Sell
var TrackingMarket = schema.GlobalMarketBuySellMetric {
	Metric: schema.Metric{Name:"Global Market Buy/Sell", Description:"Map of all items that have been bought or sold, and how many times each has been bought and sold."},
	MarketData: make(map[string]schema.GMBSMarketData),
}
// TODO: Test
func TrackMarketBuySell(itemName string, isBuy bool, quantity uint64) {
	log.Debug.Printf("Metrics:TrackMarketBuySell")
	existingData, edOK := TrackingMarket.MarketData[itemName]
	if !edOK {
		// New Data
		if isBuy {
			TrackingMarket.MarketData[itemName] = schema.GMBSMarketData{
				Bought: quantity,
				Sold: 0,
			}
		} else {
			TrackingMarket.MarketData[itemName] = schema.GMBSMarketData{
				Bought: 0,
				Sold: quantity,
			}
		}
	} else {
		// Existing Data
		if isBuy {
			existingData.Bought += quantity
			TrackingMarket.MarketData[itemName] = existingData
		} else {
			existingData.Sold += quantity
			TrackingMarket.MarketData[itemName] = existingData
		}
	}
	SaveMetrics()
}

// Plants Harvested
var TrackingHarvests = schema.TrackingHarvestsMetric {
	Metric: schema.Metric{Name:"Plants Harvested", Description:"Map of all plants that have been harvested and how many times that has occurred."},
	HarvestData: make(map[string]uint64),
}
// TODO: Test
func TrackHarvest(plantName string) {
	log.Debug.Printf("Metrics:TrackHarvest")
	TrackingHarvests.HarvestData[plantName] ++
	SaveMetrics()
}

// Rituals Cast
var TrackingRituals = schema.TrackingRitualsMetric {
	Metric: schema.Metric{Name:"Rituals Cast", Description:"Map of all rituals that have been cast and how many times that has occurred."},
	RitualData: make(map[string]uint64),
}
// TODO: Test
func TrackRitual(riteRunes string, riteName string) {
	log.Debug.Printf("Metrics:TrackRitual")
	TrackingRituals.RitualData[riteRunes + ": " + riteName] ++
	SaveMetrics()
}


// // Users by Achievement
// var TrackingUsersByAchievement = schema.UsersByAchievementMetric {
// 	Metric: schema.Metric{Name:"Users By Achievement", Description:"List of all achievements and the users who have achieved them."},
// 	UsersByAchievement: make([]schema.AchievementMetric, 0),
// }
// func CalculateUsersByAchievement() ([]schema.AchievementMetric) {
// 	return TrackingUsersByAchievement.UsersByAchievement
// }
