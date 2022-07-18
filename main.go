package main

import (
	"apricate/auth"
	"apricate/filemngr"
	"apricate/log"
	"apricate/metrics"
	"apricate/rdb"
	"apricate/schema"
	"context"
	"strings"
)

// Global vars
var (
	server_config_env_path = "data/server_config.env"
	auth_secret_path = "data/secrets.env"
	slur_filter_path = "data/slur_filter.txt"
	flush_DBs = false
	regen_auth_sec = false
	// Define relationship between string database name and redis db
	dbs = make(map[string]rdb.Database)
// 	world schema.World
// 	main_dictionary = schema.MainDictionary{}
)

func HandleServerConfigFlushDBs(lines []string) []string {
	fdb_lines_index, flush_dbs, lines := GetValueFromServerConfigByKey(lines, "flush_dbs", "false")
	if flush_dbs == "true" || flush_dbs == "dev" {
		// flush
		log.Important.Printf("flush_dbs = %s, flushing", flush_dbs)
		flush_DBs = true
	} else {
		log.Debug.Printf("flush_dbs neither true nor dev, skipping flush")
	}
	// Update
	if flush_dbs != "dev" {
		lines[fdb_lines_index] = "flush_dbs=false"
	}
	return lines
}

func HandleServerConfigRegenerateAuthSecret(lines []string, auth_secret_path string) []string {
	ras_lines_index, regenerate_auth_secret, lines := GetValueFromServerConfigByKey(lines, "regenerate_auth_secret", "false")
	if regenerate_auth_secret == "true" {
		// regen
		log.Important.Printf("regenerate_auth_secret = %s, flushing", regenerate_auth_secret)
		auth.CreateOrUpdateAuthSecretInFile(auth_secret_path)
		regen_auth_sec = true
		// Update lines
		lines[ras_lines_index] = "regenerate_auth_secret=false"
	} else {
		log.Debug.Printf("regenerate_auth_secret neither true nor dev, skipping flush")
	}
	return lines
}

func GetValueFromServerConfigByKey(lines []string, key string, default_value string) (int, string, []string) {
	key_lines_index, key_value := filemngr.GetKeyFromLines(key, lines)
	if key_value == "" {
		log.Important.Printf("key %s not found in server_config.env, creating with default value of %s", key, default_value)
		// Create key:value in env file since could not find one to update
		// If empty file then replace 1st line else append to end
		log.Debug.Printf("Creating new key:value in env file. server_config.env[0] == ''? %v", lines[0] == "")
		if lines[0] == "" {
			log.Debug.Printf("Line 0 empty in server_config.env, replacing line 0")
			lines[0] = key + "=" + default_value
			key_lines_index = 0
			key_value = default_value
		} else {
			log.Debug.Printf("Not blank server_config.env, appending to end")
			lines = append(lines, key + "=" + default_value)
			key_lines_index = len(lines)-1
			key_value = default_value
		}
	}
	return key_lines_index, key_value, lines
}

func ProcessServerConfig(lines []string, auth_secret_path string) ([]string, []string) {
	lines = HandleServerConfigFlushDBs(lines)
	lines = HandleServerConfigRegenerateAuthSecret(lines, auth_secret_path)
	_, listen_port, lines := GetValueFromServerConfigByKey(lines, "listen_port", ":8080")
	_, redis_addr, lines := GetValueFromServerConfigByKey(lines, "redis_addr", "rdb:6379")
	_, api_version, lines := GetValueFromServerConfigByKey(lines, "api_version", "0.5.0")
	misc_res := []string{listen_port, redis_addr, api_version}
	return lines, misc_res
}

// TODO: Test
func initialize_dbs(redis_addr string) {
	log.Info.Printf("Connecting to Redis server at %s", redis_addr)

	dbs["users"] = rdb.NewDatabase(redis_addr, 0)
	dbs["assistants"] = rdb.NewDatabase(redis_addr, 1)
	dbs["farms"] = rdb.NewDatabase(redis_addr, 2)
	dbs["contracts"] = rdb.NewDatabase(redis_addr, 3)
	dbs["warehouses"] = rdb.NewDatabase(redis_addr, 4)
	dbs["caravans"] = rdb.NewDatabase(redis_addr, 5)
	dbs["clearinghouse"] = rdb.NewDatabase(redis_addr, 5)

	// Ping server
	_, err := dbs["users"].Goredis.Ping(context.Background()).Result()
	if err != nil {
		log.Error.Fatalf("Could not ping redis server at %s", redis_addr)
	}

	// Check to flush DBs
	regenerate_auth_secret := regen_auth_sec
	log.Info.Printf("Check Flush DBs: %v || %v : %v", flush_DBs, regenerate_auth_secret, flush_DBs || regenerate_auth_secret)
	if flush_DBs || regenerate_auth_secret {
		for _, db := range dbs {
			db.Flush()
		}
	}
}

// TODO: Test
func setup_my_character() {
	if flush_DBs || regen_auth_sec {
		schema.PregenerateUser("Greenitthe", dbs, true)
		metrics.TrackNewUser("Greenitthe")
		schema.PregenerateUser("Viridis", dbs, false)
		metrics.TrackNewUser("Viridis")
		schema.PregenerateUser("Green", dbs, true)
		metrics.TrackNewUser("Green")
	}
	log.Info.Println("Neither flushing DBs, nor regenerating auth secret. Token for user: Greenitthe should already exist in secrets.env. Skipping creation")
}

func main() {
	log.Info.Printf("--==Apricate.io REST API Server==--")
	
	// Load config or use defaults
	log.Info.Printf("Loading Server Config from %s", server_config_env_path)
	server_config_lines := filemngr.LoadEnvFileToLines(server_config_env_path)
	server_config_lines, misc_server_config := ProcessServerConfig(server_config_lines, auth_secret_path)
	listen_port := misc_server_config[0]
	redis_addr := misc_server_config[1]
	api_version := misc_server_config[2]
	log.Info.Println(listen_port, redis_addr, api_version)
	filemngr.WriteEnvLinesToFile(server_config_env_path, server_config_lines)

	log.Info.Printf("Listen Port %s", listen_port)
	log.Info.Printf("Redis Address %s", redis_addr)
	log.Info.Printf("API Version %s", api_version)
	log.Info.Printf("Finished Loading Server Config")

	log.Info.Println("Loading secrets from envfile")
	auth.LoadSecretsToEnv()

	// Setup redis databases for each namespace
	initialize_dbs(redis_addr)

	// Reset or Load Metrics
	log.Info.Printf("Loading metrics.yaml")
	if flush_DBs || regen_auth_sec {
		// Need to reset metrics
		log.Important.Printf("Cleared data/metrics.yaml")
		filemngr.DeleteIfExists("data/metrics.yaml")
	}
	metrics.LoadMetrics()

	setup_my_character()

	// Preload 
	// Ensure exists
	filemngr.Touch(slur_filter_path)
	// Read file to lines array splitting by newline
	read_slur_filter, readErr := filemngr.ReadFileToLineSlice(slur_filter_path)
	if readErr != nil {
		// Auth is mission-critical, using Fatal
		log.Error.Fatalf("Could not read lines from slur_filter.txt. Err: %v", readErr)
	}
	slur_filter := make([]string, len(read_slur_filter))
	for i, word := range read_slur_filter {
		slur_filter[i] = strings.ToUpper(word)
	}
	log.Info.Printf("Created/Loaded Username Slur Filter")

	// // Load World from YAML
	// world = schema.World_load("./yaml/world/regions.yaml", "./yaml/world/islands", "./yaml/world/locations")
	// log.Debug.Println(world)
	// log.Info.Printf("Loaded world")

	// // Initialize dictionaries
	// initialize_dictionaries()

	// // Begin Serving
	// handle_requests(slur_filter)
}