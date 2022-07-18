package main

import (
	"apricate/auth"
	"apricate/filemngr"
	"apricate/log"
)

// Global vars
var (
	server_config_env_path = "data/server_config.env"
	auth_secret_path = "data/secrets.env"
//  // Define relationship between string database name and redis db
// 	dbs = make(map[string]rdb.Database)
// 	world schema.World
// 	main_dictionary = schema.MainDictionary{}
)

func HandleServerConfigFlushDBs(lines []string) []string {
	fdb_lines_index, flush_dbs, lines := GetValueFromServerConfigByKey(lines, "flush_dbs", "false")
	if flush_dbs == "true" || flush_dbs == "dev" {
		// flush
		log.Important.Printf("flush_dbs = %s, flushing", flush_dbs)
		//TODO
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

	// Setup redis databases for each namespace
	// initialize_dbs()


}