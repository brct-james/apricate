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
	UserActivity []UserCallTimestamp `json:"user_activity" binding:"required"` //usernames
}
type UserCallTimestamp struct {
	Username string `json:"username" binding:"required"`
	LastCallTimestamp int64 `json:"last_call_timestamp" binding:"required"`
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
