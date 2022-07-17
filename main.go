package main

import (
	"apricate/filemngr"
	"apricate/log"
)

// TODO: Test
func LoadEnvFileToLines(env_path string) []string {
	// Ensure exists
	filemngr.Touch(env_path)
	// Load config file
	lines, readErr := filemngr.ReadFileToLineSlice(env_path)
	if readErr != nil {
		// is mission-critical, using Fatal
		log.Error.Fatalf("Could not read lines from %s. Err: %v", env_path, readErr)
	}
	return lines
}

// TODO: Test
func WriteEnvLinesToFile(env_path string, lines []string) {
	// Join and write out
	writeErr := filemngr.WriteLinesToFile(env_path, lines)
	if writeErr != nil {
		log.Error.Fatalf("Could not write %s: %v", env_path, writeErr)
	}
	log.Info.Printf("Wrote config to %s", env_path)
}

func HandleServerConfigFlushDBs(lines []string) []string {
	fdb_lines_index, flush_dbs := filemngr.GetKeyFromLines("flush_dbs", lines)
	if flush_dbs == "" {
		log.Important.Printf("flush_dbs not found in server_config.env, creating with default value of false")
		// Create secret in env file since could not find one to update
		// If empty file then replace 1st line else append to end
		log.Debug.Printf("Creating new secret in env file. server_config.env[0] == ''? %v", lines[0] == "")
		if lines[0] == "" {
			log.Debug.Printf("Line 0 empty in server_config.env, replacing line 0")
			lines[0] = "flush_dbs=false"
			fdb_lines_index = 0
			flush_dbs = "false"
		} else {
			log.Debug.Printf("Not blank server_config.env, appending to end")
			lines = append(lines, "flush_dbs=false")
			fdb_lines_index = len(lines)-1
			flush_dbs = "false"
		}
	}
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

func HandleServerConfigRegenerateAuthSecret(lines []string) []string {
	ras_lines_index, regenerate_auth_secret := filemngr.GetKeyFromLines("regenerate_auth_secret", lines)
	if regenerate_auth_secret == "" {
		log.Important.Printf("regenerate_auth_secret not found in server_config.env, creating with default value of false")
		// Create secret in env file since could not find one to update
		// If empty file then replace 1st line else append to end
		log.Debug.Printf("Creating new secret in env file. server_config.env[0] == ''? %v", lines[0] == "")
		if lines[0] == "" {
			log.Debug.Printf("Line 0 empty in server_config.env, replacing line 0")
			lines[0] = "regenerate_auth_secret=false"
			ras_lines_index = 0
			regenerate_auth_secret = "false"
		} else {
			log.Debug.Printf("Not blank server_config.env, appending to end")
			lines = append(lines, "regenerate_auth_secret=false")
			ras_lines_index = len(lines)-1
			regenerate_auth_secret = "false"
		}
	}
	if regenerate_auth_secret == "true" {
		// regen
		log.Important.Printf("regenerate_auth_secret = %s, flushing", regenerate_auth_secret)
		//TODO
		// Update
		lines[ras_lines_index] = "regenerate_auth_secret=false"
	} else {
		log.Debug.Printf("regenerate_auth_secret neither true nor dev, skipping flush")
	}
	return lines
}

func ProcessServerConfig(lines []string) []string {
	lines = HandleServerConfigFlushDBs(lines)
	lines = HandleServerConfigRegenerateAuthSecret(lines)

	return lines
}

func main() {
	log.Important.Printf("Hello, World")
	server_config_lines := LoadEnvFileToLines("data/server_config.env")
	server_config_lines = ProcessServerConfig(server_config_lines)
	WriteEnvLinesToFile("data/server_config.env", server_config_lines)
}