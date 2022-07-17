package main

import (
	"apricate/filemngr"
	"apricate/log"
)

// TODO: Test
func LoadEnvToLines(env_path string) []string {
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

func ProcessServerConfig(lines []string) {
	filemngr.GetKeyFromLines("", lines)
}

func main() {
	log.Important.Printf("Hello, World")
	server_config_lines := LoadEnvToLines("data/server_config.env")
	ProcessServerConfig(server_config_lines)
}