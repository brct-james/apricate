// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

type Metric struct {
	Name string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type UsersMetricEndpointResponse struct {
	UniqueUsers []string `json:"unique_users" binding:"required"`
	ActiveUsers []string `json:"active_users" binding:"required"`
	// UsersByAchievement []AchievementMetric `json:"users-by-achievement" binding:"required"`
}

// Unique Users
type UniqueUsersMetric struct {
	Metric
	Usernames []string `json:"usernames" binding:"required"` //usernames
}

// Active Users
type ActiveUsersMetric struct {
	Metric
	UserActivity map[string]int64 `json:"user_activity" binding:"required"` //usernames
}

// User Coins
type UserCoinsMetric struct {
	Metric
	Coins map[string]uint64 `json:"coins" binding:"required"`
}

// Global Market Buy/Sell
type GlobalMarketBuySellMetric struct {
	Metric
	MarketData map[string]GMBSMarketData `json:"market_item_data" binding:"required"`
}
type GMBSMarketData struct {
	Bought uint64 `json:"bought" binding:"required"`
	Sold uint64 `json:"sold" binding:"required"`
}

// // Users by Achievement
// type UsersByAchievementMetric struct {
// 	Metric
// 	UsersByAchievement []AchievementMetric `json:"users_by_achievement" binding:"required"`
// }
// type AchievementMetric struct {
// 	Thing // name,symbol,description of particular achievement - may want to substitute this once achievements are made
// 	Users []string `json:"users" binding:"required"` //usernames
// }
