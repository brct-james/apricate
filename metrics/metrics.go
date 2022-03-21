// Package metrics defines functions for tracking and displaying various game and server metrics
package metrics

import (
	"apricate/schema"
	"apricate/timecalc"
	"fmt"
	"time"
)

// Helper Functions

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
func CalculateUniqueUsers() ([]string) {
	return TrackingUniqueUsers.Usernames
}
func TrackNewUser(username string) {
	TrackingUniqueUsers.Usernames = append(TrackingUniqueUsers.Usernames, username)
	TrackUserCall(username)
}

// Active Users
var ActivityThresholdInMinutes int = 60
var TrackingActiveUsers = schema.ActiveUsersMetric {
	Metric: schema.Metric{Name:"Active Users", Description:fmt.Sprintf("List of every user who is considered active: have registered as a new user or hit a secure endpoint in the last %d minutes.", ActivityThresholdInMinutes)},
	UserActivity: make(map[string]int64, 0),
}
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
func TrackUserCall(username string) {
	TrackingActiveUsers.UserActivity[username] = time.Now().Unix()
}

// User Coins
// See schema.User
var TrackingUserCoins = &schema.TrackingUserCoins

// Global Market Buy/Sell
var TrackingMarket = schema.GlobalMarketBuySellMetric {
	Metric: schema.Metric{Name:"Global Market Buy/Sell", Description:"Map of all items that have been bought or sold, and how many times each has been bought and sold."},
	MarketData: make(map[string]schema.GMBSMarketData),
}
func TrackMarketBuySell(itemName string, isBuy bool, quantity uint64) {
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
}

// // Users by Achievement
// var TrackingUsersByAchievement = schema.UsersByAchievementMetric {
// 	Metric: schema.Metric{Name:"Users By Achievement", Description:"List of all achievements and the users who have achieved them."},
// 	UsersByAchievement: make([]schema.AchievementMetric, 0),
// }
// func CalculateUsersByAchievement() ([]schema.AchievementMetric) {
// 	return TrackingUsersByAchievement.UsersByAchievement
// }
