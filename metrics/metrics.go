// Package metrics defines functions for tracking and displaying various game and server metrics
package metrics

import (
	"apricate/schema"
	"apricate/timecalc"
	"fmt"
	"strings"
	"time"
)

// Helper Functions

// Return found, index, UserCallTimestamp
func findActiveUserByName(userActivityTimestamps []schema.UserCallTimestamp, username string) (bool, int, schema.UserCallTimestamp) {
	for index, user := range userActivityTimestamps {
		if strings.EqualFold(username, user.Username) {
			// found user
			return true, index, user
		}
	}
	return false, -1, schema.UserCallTimestamp{}
}

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

// // Active Users
var ActivityThresholdInMinutes int = 60
var TrackingActiveUsers = schema.ActiveUsersMetric {
	Metric: schema.Metric{Name:"Active Users", Description:fmt.Sprintf("List of every user who is considered active: have registered as a new user or hit a secure endpoint in the last %d minutes.", ActivityThresholdInMinutes)},
	UserActivity: make([]schema.UserCallTimestamp, 0),
}
func CalculateActiveUsers() ([]string) {
	res := make([]string, 0)
	for _, user := range TrackingActiveUsers.UserActivity {
		exclusion_time := timecalc.AddMinutesToTimestamp(time.Unix(user.LastCallTimestamp, 0), ActivityThresholdInMinutes)
		if exclusion_time.After(time.Now()) {
			//include user from active users, as exclusion time in future
			res = append(res, user.Username)
		}
	}
	return res
}
func TrackUserCall(username string) {
	foundUser, userIndex, _ := findActiveUserByName(TrackingActiveUsers.UserActivity, username)
	if !foundUser {
		// New User
		TrackingActiveUsers.UserActivity = append(TrackingActiveUsers.UserActivity, schema.UserCallTimestamp{Username: username, LastCallTimestamp: time.Now().Unix()})
		return
	}
	// Existing user
	TrackingActiveUsers.UserActivity[userIndex].LastCallTimestamp = time.Now().Unix()
}

// // Users by Achievement
// var TrackingUsersByAchievement = schema.UsersByAchievementMetric {
// 	Metric: schema.Metric{Name:"Users By Achievement", Description:"List of all achievements and the users who have achieved them."},
// 	UsersByAchievement: make([]schema.AchievementMetric, 0),
// }
// func CalculateUsersByAchievement() ([]schema.AchievementMetric) {
// 	return TrackingUsersByAchievement.UsersByAchievement
// }
