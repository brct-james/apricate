// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

import (
	"apricate/filemngr"
	"apricate/log"

	"gopkg.in/yaml.v3"
)

type Metric struct {
	Name string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type SaveMetricsYaml struct {
	UniqueUsers []string `yaml:"UniqueUsers"`
	UserActivity map[string]int64 `yaml:"UserActivity"`
	Coins map[string]uint64 `yaml:"UserCoins"`
	MarketData map[string]GMBSMarketData `yaml:"MarketData"`
}

type UsersMetricEndpointResponse struct {
	UniqueUsers []string `json:"unique_users" binding:"required"`
	ActiveUsers []string `json:"active_users" binding:"required"`
	// UsersByAchievement []AchievementMetric `json:"users-by-achievement" binding:"required"`
}

// Unique Users
type UniqueUsersMetric struct {
	Metric
	Usernames []string `yaml:"Usernames" json:"usernames" binding:"required"` //usernames
}

// Active Users
type ActiveUsersMetric struct {
	Metric
	UserActivity map[string]int64 `yaml:"UserActivity" json:"user_activity" binding:"required"` //usernames
}

// User Coins
type UserCoinsMetric struct {
	Metric
	Coins map[string]uint64 `yaml:"UserCoins" json:"coins" binding:"required"`
}

// Global Market Buy/Sell
type GlobalMarketBuySellMetric struct {
	Metric
	MarketData map[string]GMBSMarketData `yaml:"GlobalMarketBuySell" json:"market_item_data" binding:"required"`
}
type GMBSMarketData struct {
	Bought uint64 `yaml:"Bought" json:"bought" binding:"required"`
	Sold uint64 `yaml:"Sold" json:"sold" binding:"required"`
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

// Load metrics struct by unmarhsalling given yaml file
func Metrics_from_yaml(path_to_metrics_yaml string) (SaveMetricsYaml, bool) {
	metricsBytes, readErr := filemngr.ReadFileToBytes(path_to_metrics_yaml)
	if readErr != nil {
		log.Error.Printf("Read Error in metrics_from_yaml: %v", readErr)
		return SaveMetricsYaml{}, false
	}
	var metrics SaveMetricsYaml
	err := yaml.Unmarshal(metricsBytes, &metrics)
	if err != nil {
		log.Error.Printf("Error in metrics_from_yaml: %v", err)
		return SaveMetricsYaml{}, false
	}
	return metrics, true
}

// Save metrics struct by marhsalling given yaml file
func Metrics_to_yaml(path_to_metrics_yaml string, mYaml SaveMetricsYaml) {
	data, err := yaml.Marshal(&mYaml)
	if err != nil {
		log.Error.Printf("Error in metrics_to_yaml: %v", err)
	}
	writeErr := filemngr.WriteBytesToFile(path_to_metrics_yaml, data)
	if writeErr != nil {
		log.Error.Printf("Write Error in metrics_to_yaml: %v", err)
	}
}