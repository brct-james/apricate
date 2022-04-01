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

type MetricsResponse struct {
	MarketBuySell GlobalMarketBuySellMetric `json:"Global Market Buy/Sell" binding:"required"`
	UserCoins UserCoinsMetric `json:"User Coins" binding:"required"`
	Harvests TrackingHarvestsMetric `json:"Harvests" binding:"required"`
	Rituals TrackingRitualsMetric `json:"Rituals" binding:"required"`
	UserMagic UserMagicMetric `json:"User Magic" binding:"required"`
}

type SaveMetricsYaml struct {
	UniqueUsers []string `yaml:"UniqueUsers"`
	UserActivity map[string]int64 `yaml:"UserActivity"`
	Coins map[string]uint64 `yaml:"UserCoins"`
	MarketData map[string]GMBSMarketData `yaml:"MarketData"`
	HarvestData map[string]uint64 `yaml:"HarvestData"`
	RitualData map[string]uint64 `yaml:"RitualData"`
	UserMagic map[string]map[string]float64 `yaml:"UserMagic"`
}

type UsersMetricEndpointResponse struct {
	UniqueUsers []string `json:"unique_users" binding:"required"`
	ActiveUsers []string `json:"active_users" binding:"required"`
	// UsersByAchievement []AchievementMetric `json:"users-by-achievement" binding:"required"`
}

// Tracking User Coins for Metrics
var TrackingUserCoins = UserCoinsMetric {
	Metric: Metric{Name:"User Coins", Description:"Map of every registered user and their coins",},
	Coins: make(map[string]uint64),
}
func TrackUserCoins(username string, coins uint64) {
	log.Debug.Printf("Metrics:TrackUserCoins")
	TrackingUserCoins.Coins[username] = coins
}

// Tracking User Flux for Metrics
var TrackingUserMagic = UserMagicMetric {
	Metric: Metric{Name:"User Magic Stats", Description:"Map of every registered user, their arcane flux, and distortion tier",},
	Magic: make(map[string]map[string]float64),
}
func TrackUserMagic(username string, flux float64, distortionTier float64) {
	log.Debug.Printf("Metrics:TrackUserMagic")
	userMagic, mOk := TrackingUserMagic.Magic[username]
	if !mOk {
		userMagic = make(map[string]float64)
	}
	userMagic["Arcane Flux"] = flux
	userMagic["Distortion Tier"] = distortionTier

	TrackingUserMagic.Magic[username] = userMagic
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

// User Magic
type UserMagicMetric struct {
	Metric
	Magic map[string]map[string]float64 `yaml:"UserMagic" json:"magic" binding:"required"`
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

// Plants Harvested
type TrackingHarvestsMetric struct {
	Metric
	HarvestData map[string]uint64 `yaml:"TotalHarvests" json:"total_harvests" binding:"required"`
}

// Rituals Cast
type TrackingRitualsMetric struct {
	Metric
	RitualData map[string]uint64 `yaml:"TotalRituals" json:"total_rituals" binding:"required"`
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
	log.Debug.Printf("Load metrics from %s", path_to_metrics_yaml)
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
	log.Debug.Printf("Save metrics to %s", path_to_metrics_yaml)
	data, err := yaml.Marshal(&mYaml)
	if err != nil {
		log.Error.Printf("Error in metrics_to_yaml: %v", err)
	}
	writeErr := filemngr.WriteBytesToFile(path_to_metrics_yaml, data)
	if writeErr != nil {
		log.Error.Printf("Write Error in metrics_to_yaml: %v", err)
	}
}